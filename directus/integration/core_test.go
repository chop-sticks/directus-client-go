//go:build integration

package integration

import (
	"fmt"
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

// TestE2ECoreCollections exercises the full collection lifecycle end to end,
// creating its own collection (rather than the shared helper) so CreateCollection,
// PatchCollection, PatchCollectionsBatch, and DeleteCollection are all covered.
func TestE2ECoreCollections(t *testing.T) {
	c := itestClient(t)
	name := uniqueName("e2e_core_col")

	created, err := c.CreateCollection(&CollectionRequest{
		Collection: name,
		Meta:       &CollectionMeta{Note: "e2e-core", Icon: "science"},
		Schema:     &CollectionSchema{},
		Fields: []Field{
			{Field: "id", Type: "integer", Meta: &FieldMeta{Hidden: true}, Schema: &FieldSchema{IsPrimaryKey: true, HasAutoIncrement: true}},
			{Field: "title", Type: "string", Meta: &FieldMeta{Interface: "input"}, Schema: &FieldSchema{}},
		},
	}, nil)
	if err != nil {
		t.Fatalf("CreateCollection: %v", err)
	}
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			_ = c.DeleteCollection(name)
		}
	})
	if created == nil || created.Collection != name {
		t.Fatalf("CreateCollection returned %+v, want collection %q", created, name)
	}

	cols, err := c.GetCollections()
	if err != nil {
		t.Fatalf("GetCollections: %v", err)
	}
	if len(cols) == 0 {
		t.Fatalf("GetCollections returned no collections")
	}

	got, err := c.GetCollectionByName(name)
	if err != nil {
		t.Fatalf("GetCollectionByName: %v", err)
	}
	if got == nil || got.Collection != name {
		t.Fatalf("GetCollectionByName returned %+v, want collection %q", got, name)
	}

	patched, err := c.PatchCollection(name, &CollectionRequest{
		Meta: &CollectionMeta{Note: "e2e-core-patched"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchCollection: %v", err)
	}
	if patched == nil || patched.Meta == nil || patched.Meta.Note != "e2e-core-patched" {
		t.Fatalf("PatchCollection note = %+v, want \"e2e-core-patched\"", patched)
	}

	batched, err := c.PatchCollectionsBatch([]CollectionRequest{
		{Collection: name, Meta: &CollectionMeta{Note: "e2e-core-batch"}},
	}, nil)
	if err != nil {
		t.Fatalf("PatchCollectionsBatch: %v", err)
	}
	if len(batched) == 0 {
		t.Fatalf("PatchCollectionsBatch returned no collections")
	}

	if err := c.DeleteCollection(name); err != nil {
		t.Fatalf("DeleteCollection: %v", err)
	}
	deleted = true
}

// TestE2ECoreFields exercises every field method against a throwaway collection.
func TestE2ECoreFields(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)

	created, err := c.CreateField(coll, &Field{
		Field:  "title2",
		Type:   "string",
		Meta:   &FieldMeta{Interface: "input"},
		Schema: &FieldSchema{},
	}, nil)
	if err != nil {
		t.Fatalf("CreateField: %v", err)
	}
	if created == nil || created.Field != "title2" {
		t.Fatalf("CreateField returned %+v, want field \"title2\"", created)
	}

	all, err := c.GetFields()
	if err != nil {
		t.Fatalf("GetFields: %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("GetFields returned no fields")
	}

	byColl, err := c.GetFieldsByCollection(coll)
	if err != nil {
		t.Fatalf("GetFieldsByCollection: %v", err)
	}
	if len(byColl) == 0 {
		t.Fatalf("GetFieldsByCollection returned no fields")
	}

	got, err := c.GetFieldByCollectionAndName(coll, "title2")
	if err != nil {
		t.Fatalf("GetFieldByCollectionAndName: %v", err)
	}
	if got == nil || got.Field != "title2" {
		t.Fatalf("GetFieldByCollectionAndName returned %+v, want field \"title2\"", got)
	}

	pf, err := c.PatchField(coll, "title2", &Field{
		Meta: &FieldMeta{Interface: "input", Note: "e2e-core"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchField: %v", err)
	}
	if pf == nil || pf.Field != "title2" {
		t.Fatalf("PatchField returned %+v, want field \"title2\"", pf)
	}

	pfs, err := c.PatchFields(coll, []Field{
		{Field: "title2", Type: "string", Meta: &FieldMeta{Interface: "input", Note: "e2e-core-batch"}, Schema: &FieldSchema{}},
	}, nil)
	if err != nil {
		t.Fatalf("PatchFields: %v", err)
	}
	if len(pfs) == 0 {
		t.Fatalf("PatchFields returned no fields")
	}

	if err := c.DeleteField(coll, "title2"); err != nil {
		t.Fatalf("DeleteField: %v", err)
	}
}

// TestE2ECoreItems exercises every item method against a throwaway collection.
func TestE2ECoreItems(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)

	item, err := c.CreateItem(coll, map[string]any{"title": "one"}, nil)
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	key1 := e2eCoreKey(item["id"])
	if key1 == "" {
		t.Fatalf("CreateItem returned no id: %+v", item)
	}

	created, err := c.CreateItems(coll, []map[string]any{
		{"title": "two"},
		{"title": "three"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateItems: %v", err)
	}
	if len(created) != 2 {
		t.Fatalf("CreateItems returned %d items, want 2", len(created))
	}
	key2 := e2eCoreKey(created[0]["id"])
	key3 := e2eCoreKey(created[1]["id"])

	items, err := c.GetItems(coll, &Query{Fields: []string{"*"}})
	if err != nil {
		t.Fatalf("GetItems: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("GetItems returned %d items, want 3", len(items))
	}

	one, err := c.GetItem(coll, key1, nil)
	if err != nil {
		t.Fatalf("GetItem: %v", err)
	}
	if one["title"] != "one" {
		t.Fatalf("GetItem title = %v, want \"one\"", one["title"])
	}

	patched, err := c.PatchItem(coll, key1, map[string]any{"title": "one-updated"}, nil)
	if err != nil {
		t.Fatalf("PatchItem: %v", err)
	}
	if patched["title"] != "one-updated" {
		t.Fatalf("PatchItem title = %v, want \"one-updated\"", patched["title"])
	}

	many, err := c.PatchItems(coll, []any{key2, key3}, map[string]any{"title": "bulk"}, nil)
	if err != nil {
		t.Fatalf("PatchItems: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("PatchItems returned %d items, want 2", len(many))
	}

	batch, err := c.PatchItemsBatch(coll, []map[string]any{
		{"id": key2, "title": "batch2"},
		{"id": key3, "title": "batch3"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchItemsBatch: %v", err)
	}
	if len(batch) == 0 {
		t.Fatalf("PatchItemsBatch returned no items")
	}

	if err := c.DeleteItem(coll, key1); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}

	if err := c.DeleteItems(coll, []any{key2, key3}); err != nil {
		t.Fatalf("DeleteItems: %v", err)
	}

	remaining, err := c.GetItems(coll, nil)
	if err != nil {
		t.Fatalf("GetItems after delete: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected 0 items after delete, got %d", len(remaining))
	}

	if _, err := c.CreateItem(coll, map[string]any{"title": "tofilter"}, nil); err != nil {
		t.Fatalf("CreateItem (tofilter): %v", err)
	}
	if err := c.DeleteItemsByQuery(coll, map[string]any{"title": map[string]any{"_eq": "tofilter"}}); err != nil {
		t.Fatalf("DeleteItemsByQuery: %v", err)
	}
	after, err := c.GetItems(coll, nil)
	if err != nil {
		t.Fatalf("GetItems after query delete: %v", err)
	}
	if len(after) != 0 {
		t.Fatalf("expected 0 items after query delete, got %d", len(after))
	}
}

// TestE2ECoreSingleton exercises the singleton read/write methods.
func TestE2ECoreSingleton(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)

	if _, err := c.PatchCollection(coll, &CollectionRequest{
		Meta: &CollectionMeta{Singleton: true},
	}, nil); err != nil {
		t.Fatalf("PatchCollection (singleton): %v", err)
	}

	if _, err := c.GetSingleton(coll, nil); err != nil {
		t.Fatalf("GetSingleton: %v", err)
	}

	if _, err := c.PatchSingleton(coll, map[string]any{"title": "x"}, nil); err != nil {
		t.Fatalf("PatchSingleton: %v", err)
	}

	got, err := c.GetSingleton(coll, &Query{Fields: []string{"*"}})
	if err != nil {
		t.Fatalf("GetSingleton after patch: %v", err)
	}
	if got["title"] != "x" {
		t.Fatalf("GetSingleton title = %v, want \"x\"", got["title"])
	}
}

// TestE2ECoreAggregate exercises the Aggregate method.
func TestE2ECoreAggregate(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)

	newTestItem(t, c, coll, "agg-one")
	newTestItem(t, c, coll, "agg-two")

	res, err := c.Aggregate(coll, &Query{Aggregate: map[string]any{"count": "*"}})
	if err != nil {
		t.Fatalf("Aggregate: %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("Aggregate returned no rows")
	}
}

// e2eCoreKey renders an item primary key returned as any into a URL-safe string.
func e2eCoreKey(v any) string {
	if v == nil {
		return ""
	}
	s := fmt.Sprintf("%v", v)
	if s == "<nil>" {
		return ""
	}
	return s
}
