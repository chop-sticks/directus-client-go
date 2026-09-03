//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

func TestE2ESchemaRelations(t *testing.T) {
	c := itestClient(t)

	collA := newTestCollection(t, c)
	collB := newTestCollection(t, c)

	if _, err := c.CreateField(collA, &Field{
		Field:  "b_id",
		Type:   "integer",
		Meta:   &FieldMeta{},
		Schema: &FieldSchema{},
	}, nil); err != nil {
		t.Fatalf("CreateField: %v", err)
	}

	created, err := c.CreateRelation(&Relation{
		Collection:        collA,
		Field:             "b_id",
		RelatedCollection: collB,
	})
	if err != nil {
		t.Fatalf("CreateRelation: %v", err)
	}
	if created == nil || created.Collection != collA || created.Field != "b_id" {
		t.Fatalf("CreateRelation returned unexpected relation: %+v", created)
	}
	t.Cleanup(func() { _ = c.DeleteRelation(collA, "b_id") })

	all, err := c.GetRelations()
	if err != nil {
		t.Fatalf("GetRelations: %v", err)
	}
	if all == nil {
		t.Fatalf("GetRelations returned nil slice")
	}

	byColl, err := c.GetRelationsByCollection(collA)
	if err != nil {
		t.Fatalf("GetRelationsByCollection: %v", err)
	}
	if len(byColl) == 0 {
		t.Fatalf("GetRelationsByCollection returned no relations for %s", collA)
	}

	rel, err := c.GetRelation(collA, "b_id")
	if err != nil {
		t.Fatalf("GetRelation: %v", err)
	}
	if rel == nil || rel.Field != "b_id" {
		t.Fatalf("GetRelation returned unexpected relation: %+v", rel)
	}

	patched, err := c.PatchRelation(collA, "b_id", &Relation{
		Meta: &RelationMeta{OneField: "a_items"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchRelation: %v", err)
	}
	if patched == nil || patched.Field != "b_id" {
		t.Fatalf("PatchRelation returned unexpected relation: %+v", patched)
	}

	if err := c.DeleteRelation(collA, "b_id"); err != nil {
		t.Fatalf("DeleteRelation: %v", err)
	}
}

func TestE2ESchemaSettings(t *testing.T) {
	c := itestClient(t)

	settings, err := c.GetSettings(nil)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if settings == nil {
		t.Fatalf("GetSettings returned nil")
	}

	original := settings.ProjectDescriptor
	t.Cleanup(func() {
		_, _ = c.PatchSettings(&Settings{ProjectDescriptor: original}, nil)
	})

	updated, err := c.PatchSettings(&Settings{ProjectDescriptor: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchSettings: %v", err)
	}
	if updated == nil {
		t.Fatalf("PatchSettings returned nil")
	}

	if _, err := c.PatchSettings(&Settings{ProjectDescriptor: original}, nil); err != nil {
		t.Fatalf("PatchSettings revert: %v", err)
	}
}

func TestE2ESchemaExtensions(t *testing.T) {
	c := itestClient(t)

	exts, err := c.GetExtensions()
	if err != nil {
		t.Fatalf("GetExtensions: %v", err)
	}
	if exts == nil {
		t.Fatalf("GetExtensions returned nil slice")
	}

	// Environment-dependent: registry access requires external network.
	if _, err := c.GetRegistryExtensions(nil); err != nil {
		t.Logf("GetRegistryExtensions: %v", err)
	}

	if _, err := c.PatchExtension("nonexistent", map[string]any{
		"meta": map[string]any{"enabled": true},
	}); err == nil {
		t.Fatalf("PatchExtension(nonexistent) expected error, got nil")
	}

	if _, err := c.PatchBundleExtension("nob", "non", map[string]any{
		"meta": map[string]any{"enabled": true},
	}); err == nil {
		t.Fatalf("PatchBundleExtension(nonexistent) expected error, got nil")
	}

	if err := c.DeleteExtension("nonexistent"); err == nil {
		t.Fatalf("DeleteExtension(nonexistent) expected error, got nil")
	}

	if err := c.InstallRegistryExtension("bad", "bad"); err == nil {
		t.Fatalf("InstallRegistryExtension(bad) expected error, got nil")
	}

	if err := c.UninstallRegistryExtension("bad"); err == nil {
		t.Fatalf("UninstallRegistryExtension(bad) expected error, got nil")
	}
}

func TestE2ESchemaSnapshot(t *testing.T) {
	c := itestClient(t)

	snapshot, err := c.GetSchemaSnapshot(nil, nil)
	if err != nil {
		t.Fatalf("GetSchemaSnapshot: %v", err)
	}
	if snapshot == nil {
		t.Fatalf("GetSchemaSnapshot returned nil")
	}

	// Diffing the current snapshot against itself yields no diff (nil, nil).
	diff, err := c.SchemaDiffSnapshot(snapshot, false, "")
	if err != nil {
		t.Fatalf("SchemaDiffSnapshot: %v", err)
	}
	if diff != nil {
		t.Fatalf("SchemaDiffSnapshot expected nil diff for identical snapshot, got %+v", diff)
	}

	// Applying an empty diff is not safely testable; call and log on error.
	if err := c.SchemaApply(&SchemaDiff{}, false); err != nil {
		t.Logf("SchemaApply(empty): %v", err)
	}
}
