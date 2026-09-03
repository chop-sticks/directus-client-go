//go:build integration

// Package directus integration tests exercise a live Directus instance.
//
// They are excluded from the normal unit-test build via the `integration`
// build tag. Run them against a running stack (see docker-compose.yml) with:
//
//	go test -tags=integration -count=1 ./directus/...
//
// or, more conveniently:
//
//	task test:integration
//
// The target instance and credentials are read from the environment, falling
// back to the values baked into docker-compose.yml.
package integration

import (
	"fmt"
	. "github.com/chop-sticks/directus-client-go/directus"
	"os"
	"strings"
	"testing"
)

const (
	defaultURL      = "http://localhost:8055"
	defaultToken    = "eiriezashohnai1xohjuC2aem7duuDie"
	defaultEmail    = "test@example.com"
	defaultPassword = "testAtExampleDotCom+1"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func itestClient(t *testing.T) *Client {
	t.Helper()
	host := env("DIRECTUS_URL", defaultURL)
	token := env("DIRECTUS_TOKEN", defaultToken)
	c, err := NewClient(&host, &token)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestIntegrationServer(t *testing.T) {
	c := itestClient(t)

	pong, err := c.ServerPing()
	if err != nil {
		t.Fatalf("ServerPing: %v", err)
	}
	if !strings.Contains(pong, "pong") {
		t.Errorf("ServerPing = %q, want it to contain \"pong\"", pong)
	}

	info, err := c.ServerInfo()
	if err != nil {
		t.Fatalf("ServerInfo: %v", err)
	}
	if len(info) == 0 {
		t.Error("ServerInfo returned empty data")
	}

	health, err := c.ServerHealth()
	if err != nil {
		t.Fatalf("ServerHealth: %v", err)
	}
	if health == nil || health.Status == "" {
		t.Errorf("ServerHealth returned no status: %+v", health)
	}
}

func TestIntegrationAuthAndMe(t *testing.T) {
	c := itestClient(t)

	auth, err := c.Login(defaultEmail, defaultPassword, "json", "")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if auth == nil || auth.AccessToken == "" {
		t.Fatalf("Login returned no access token: %+v", auth)
	}

	me, err := c.GetUsersMe(&Query{Fields: []string{"id", "email"}})
	if err != nil {
		t.Fatalf("GetUsersMe: %v", err)
	}
	if me == nil || me.Email != defaultEmail {
		t.Errorf("GetUsersMe email = %q, want %q", me.Email, defaultEmail)
	}
}

func TestIntegrationCollectionAndItems(t *testing.T) {
	c := itestClient(t)
	const coll = "sdk_itest"

	// Remove any leftover from a previous run, then ensure cleanup.
	_ = c.DeleteCollection(coll)
	t.Cleanup(func() { _ = c.DeleteCollection(coll) })

	created, err := c.CreateCollection(&CollectionRequest{
		Collection: coll,
		Meta:       &CollectionMeta{Icon: "science", Note: "SDK integration test", Collapse: "open"},
		Schema:     &CollectionSchema{},
		Fields: []Field{
			{
				Field:  "id",
				Type:   "integer",
				Meta:   &FieldMeta{Hidden: true, Readonly: true, Interface: "input"},
				Schema: &FieldSchema{IsPrimaryKey: true, HasAutoIncrement: true},
			},
			{
				Field:  "title",
				Type:   "string",
				Meta:   &FieldMeta{Interface: "input"},
				Schema: &FieldSchema{},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("CreateCollection: %v", err)
	}
	if created.Collection != coll {
		t.Fatalf("CreateCollection returned %q, want %q", created.Collection, coll)
	}

	t.Run("field lifecycle", func(t *testing.T) {
		field, err := c.CreateField(coll, &Field{
			Field:  "count",
			Type:   "integer",
			Meta:   &FieldMeta{Interface: "input"},
			Schema: &FieldSchema{},
		}, nil)
		if err != nil {
			t.Fatalf("CreateField: %v", err)
		}
		if field.Field != "count" {
			t.Errorf("CreateField returned %q, want \"count\"", field.Field)
		}

		got, err := c.GetFieldByCollectionAndName(coll, "count")
		if err != nil {
			t.Fatalf("GetFieldByCollectionAndName: %v", err)
		}
		if got.Type != "integer" {
			t.Errorf("field type = %q, want \"integer\"", got.Type)
		}

		if err := c.DeleteField(coll, "count"); err != nil {
			t.Fatalf("DeleteField: %v", err)
		}
	})

	t.Run("item CRUD", func(t *testing.T) {
		item, err := c.CreateItem(coll, map[string]any{"title": "hello"}, nil)
		if err != nil {
			t.Fatalf("CreateItem: %v", err)
		}
		key := fmt.Sprintf("%v", item["id"])
		if key == "" || key == "<nil>" {
			t.Fatalf("CreateItem returned no id: %+v", item)
		}

		items, err := c.GetItems(coll, &Query{Fields: []string{"*"}})
		if err != nil {
			t.Fatalf("GetItems: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("GetItems returned %d items, want 1", len(items))
		}

		patched, err := c.PatchItem(coll, key, map[string]any{"title": "updated"}, nil)
		if err != nil {
			t.Fatalf("PatchItem: %v", err)
		}
		if patched["title"] != "updated" {
			t.Errorf("PatchItem title = %v, want \"updated\"", patched["title"])
		}

		if err := c.DeleteItem(coll, key); err != nil {
			t.Fatalf("DeleteItem: %v", err)
		}

		remaining, err := c.GetItems(coll, nil)
		if err != nil {
			t.Fatalf("GetItems after delete: %v", err)
		}
		if len(remaining) != 0 {
			t.Errorf("expected 0 items after delete, got %d", len(remaining))
		}
	})
}
