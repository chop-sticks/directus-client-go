package directus

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestRelationsRequestBuildError(t *testing.T) {
	client := badHostClient(t)

	if _, err := client.GetRelations(); err == nil {
		t.Error("GetRelations: expected request build error, got nil")
	}
	if _, err := client.GetRelationsByCollection("articles"); err == nil {
		t.Error("GetRelationsByCollection: expected request build error, got nil")
	}
	if _, err := client.GetRelation("articles", "author"); err == nil {
		t.Error("GetRelation: expected request build error, got nil")
	}
	if _, err := client.CreateRelation(&Relation{}); err == nil {
		t.Error("CreateRelation: expected request build error, got nil")
	}
	if _, err := client.PatchRelation("articles", "author", &Relation{}, nil); err == nil {
		t.Error("PatchRelation: expected request build error, got nil")
	}
	if err := client.DeleteRelation("articles", "author"); err == nil {
		t.Error("DeleteRelation: expected request build error, got nil")
	}
}

func TestGetRelations(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"collection":"articles","field":"author","related_collection":"users","meta":{"id":7},"schema":{"table":"articles"}}]}`)
	})

	relations, err := client.GetRelations()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/relations" {
		t.Errorf("expected GET /relations, got %s %s", gotMethod, gotPath)
	}
	if len(relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(relations))
	}
	if relations[0].Collection != "articles" || relations[0].RelatedCollection != "users" {
		t.Errorf("unexpected relation: %+v", relations[0])
	}
	if relations[0].Meta == nil || relations[0].Meta.ID != 7 {
		t.Errorf("unexpected meta: %+v", relations[0].Meta)
	}
	if relations[0].Schema == nil || relations[0].Schema.Table != "articles" {
		t.Errorf("unexpected schema: %+v", relations[0].Schema)
	}
}

func TestGetRelationsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetRelations(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetRelationsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{`)
	})
	if _, err := client.GetRelations(); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetRelationsByCollection(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[]}`)
	})

	if _, err := client.GetRelationsByCollection("articles"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/relations/articles" {
		t.Errorf("expected GET /relations/articles, got %s %s", gotMethod, gotPath)
	}
}

func TestGetRelationsByCollectionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	_, err := client.GetRelationsByCollection("")
	if err == nil || !strings.Contains(err.Error(), "collection must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetRelation(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"articles","field":"author","related_collection":"users"}}`)
	})

	relation, err := client.GetRelation("articles", "author")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/relations/articles/author" {
		t.Errorf("expected GET /relations/articles/author, got %s %s", gotMethod, gotPath)
	}
	if relation.Field != "author" {
		t.Errorf("unexpected relation: %+v", relation)
	}
}

func TestGetRelationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	cases := []struct{ collection, field string }{
		{"", "author"},
		{"articles", ""},
		{"", ""},
	}
	for _, tc := range cases {
		_, err := client.GetRelation(tc.collection, tc.field)
		if err == nil || !strings.Contains(err.Error(), "collection and field must be provided") {
			t.Errorf("collection=%q field=%q: expected validation error, got %v", tc.collection, tc.field, err)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetRelationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetRelation("articles", "missing"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetRelationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetRelation("articles", "author"); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateRelation(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"articles","field":"author","related_collection":"users"}}`)
	})

	relation, err := client.CreateRelation(&Relation{Collection: "articles", Field: "author", RelatedCollection: "users"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/relations" {
		t.Errorf("expected POST /relations, got %s %s", gotMethod, gotPath)
	}
	if relation.Collection != "articles" {
		t.Errorf("unexpected relation: %+v", relation)
	}
}

func TestCreateRelationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateRelation(&Relation{}); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestCreateRelationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateRelation(&Relation{}); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchRelation(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"articles","field":"author","related_collection":"users"}}`)
	})

	relation, err := client.PatchRelation("articles", "author", &Relation{RelatedCollection: "users"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/relations/articles/author" {
		t.Errorf("expected PATCH /relations/articles/author, got %s %s", gotMethod, gotPath)
	}
	if relation.RelatedCollection != "users" {
		t.Errorf("unexpected relation: %+v", relation)
	}
}

func TestPatchRelationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchRelation("", "author", &Relation{}, nil); err == nil || !strings.Contains(err.Error(), "collection and field must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if _, err := client.PatchRelation("articles", "", &Relation{}, nil); err == nil {
		t.Error("expected validation error for empty field, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchRelationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchRelation("articles", "author", &Relation{}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteRelation(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteRelation("articles", "author"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/relations/articles/author" {
		t.Errorf("expected DELETE /relations/articles/author, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteRelationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if err := client.DeleteRelation("", "author"); err == nil || !strings.Contains(err.Error(), "collection and field must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if err := client.DeleteRelation("articles", ""); err == nil {
		t.Error("expected validation error for empty field, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteRelationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if err := client.DeleteRelation("articles", "author"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}
