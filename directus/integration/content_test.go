//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

func TestE2EContentPresets(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)

	// Single CRUD.
	created, err := c.CreatePreset(&Preset{Collection: coll, Bookmark: uniqueName("e2e_preset")}, nil)
	if err != nil {
		t.Fatalf("CreatePreset: %v", err)
	}
	if created == nil || created.ID == 0 {
		t.Fatalf("CreatePreset returned no id")
	}
	t.Cleanup(func() { _ = c.DeletePreset(created.ID) })

	got, err := c.GetPreset(created.ID, nil)
	if err != nil {
		t.Fatalf("GetPreset: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetPreset returned unexpected preset: %+v", got)
	}

	list, err := c.GetPresets(nil)
	if err != nil {
		t.Fatalf("GetPresets: %v", err)
	}
	if list == nil {
		t.Fatalf("GetPresets returned nil")
	}

	patched, err := c.PatchPreset(created.ID, &Preset{Icon: "star"}, nil)
	if err != nil {
		t.Fatalf("PatchPreset: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchPreset returned unexpected preset: %+v", patched)
	}

	// Batch create.
	many, err := c.CreatePresets([]Preset{
		{Collection: coll, Bookmark: uniqueName("e2e_preset")},
		{Collection: coll, Bookmark: uniqueName("e2e_preset")},
	}, nil)
	if err != nil {
		t.Fatalf("CreatePresets: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("CreatePresets returned %d presets, want 2", len(many))
	}
	keys := []int{many[0].ID, many[1].ID}

	// PatchPresets by keys.
	updated, err := c.PatchPresets(keys, &Preset{Icon: "bookmark"}, nil)
	if err != nil {
		t.Fatalf("PatchPresets: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchPresets returned %d presets, want 2", len(updated))
	}

	// PatchPresetsBatch.
	batched, err := c.PatchPresetsBatch([]Preset{
		{ID: many[0].ID, Icon: "flag"},
		{ID: many[1].ID, Icon: "flag"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchPresetsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchPresetsBatch returned %d presets, want 2", len(batched))
	}

	// DeletePresets (many).
	if err := c.DeletePresets(keys); err != nil {
		t.Fatalf("DeletePresets: %v", err)
	}

	// DeletePreset (single).
	if err := c.DeletePreset(created.ID); err != nil {
		t.Fatalf("DeletePreset: %v", err)
	}
}

func TestE2EContentTranslations(t *testing.T) {
	c := itestClient(t)

	created, err := c.CreateTranslation(&Translation{Language: "en-US", Key: uniqueName("e2e_tr"), Value: "hi"}, nil)
	if err != nil {
		t.Fatalf("CreateTranslation: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateTranslation returned no id")
	}
	t.Cleanup(func() { _ = c.DeleteTranslation(created.ID) })

	got, err := c.GetTranslation(created.ID, nil)
	if err != nil {
		t.Fatalf("GetTranslation: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetTranslation returned unexpected translation: %+v", got)
	}

	list, err := c.GetTranslations(nil)
	if err != nil {
		t.Fatalf("GetTranslations: %v", err)
	}
	if list == nil {
		t.Fatalf("GetTranslations returned nil")
	}

	patched, err := c.PatchTranslation(created.ID, &Translation{Value: "hello"}, nil)
	if err != nil {
		t.Fatalf("PatchTranslation: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchTranslation returned unexpected translation: %+v", patched)
	}

	many, err := c.CreateTranslations([]Translation{
		{Language: "en-US", Key: uniqueName("e2e_tr"), Value: "a"},
		{Language: "en-US", Key: uniqueName("e2e_tr"), Value: "b"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateTranslations: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("CreateTranslations returned %d translations, want 2", len(many))
	}
	keys := []string{many[0].ID, many[1].ID}

	updated, err := c.PatchTranslations(keys, &Translation{Value: "same"}, nil)
	if err != nil {
		t.Fatalf("PatchTranslations: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchTranslations returned %d translations, want 2", len(updated))
	}

	batched, err := c.PatchTranslationsBatch([]Translation{
		{ID: many[0].ID, Value: "x"},
		{ID: many[1].ID, Value: "y"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchTranslationsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchTranslationsBatch returned %d translations, want 2", len(batched))
	}

	if err := c.DeleteTranslations(keys); err != nil {
		t.Fatalf("DeleteTranslations: %v", err)
	}

	if err := c.DeleteTranslation(created.ID); err != nil {
		t.Fatalf("DeleteTranslation: %v", err)
	}
}

func TestE2EContentShares(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)
	itemKey := newTestItem(t, c, coll, "e2e-share")

	created, err := c.CreateShare(&Share{Name: uniqueName("e2e_share"), Collection: coll, Item: itemKey}, nil)
	if err != nil {
		t.Fatalf("CreateShare: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateShare returned no id")
	}
	t.Cleanup(func() { _ = c.DeleteShare(created.ID) })

	got, err := c.GetShare(created.ID, nil)
	if err != nil {
		t.Fatalf("GetShare: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetShare returned unexpected share: %+v", got)
	}

	list, err := c.GetShares(nil)
	if err != nil {
		t.Fatalf("GetShares: %v", err)
	}
	if list == nil {
		t.Fatalf("GetShares returned nil")
	}

	patched, err := c.PatchShare(created.ID, &Share{Name: uniqueName("e2e_share")}, nil)
	if err != nil {
		t.Fatalf("PatchShare: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchShare returned unexpected share: %+v", patched)
	}

	// Public share metadata (success).
	info, err := c.ReadShareInfo(created.ID)
	if err != nil {
		t.Fatalf("ReadShareInfo: %v", err)
	}
	if info == nil || info.ID != created.ID {
		t.Fatalf("ReadShareInfo returned unexpected info: %+v", info)
	}

	// AuthenticateShare and InviteShare are environment-dependent.
	if _, err := c.AuthenticateShare(created.ID, "", "json"); err != nil {
		t.Logf("AuthenticateShare: %v", err)
	}
	if err := c.InviteShare(created.ID, []string{"a@example.com"}); err != nil {
		t.Logf("InviteShare: %v", err)
	}

	// Batch create.
	many, err := c.CreateShares([]Share{
		{Name: uniqueName("e2e_share"), Collection: coll, Item: itemKey},
		{Name: uniqueName("e2e_share"), Collection: coll, Item: itemKey},
	}, nil)
	if err != nil {
		t.Fatalf("CreateShares: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("CreateShares returned %d shares, want 2", len(many))
	}
	keys := []string{many[0].ID, many[1].ID}

	updated, err := c.PatchShares(keys, &Share{MaxUses: 5}, nil)
	if err != nil {
		t.Fatalf("PatchShares: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchShares returned %d shares, want 2", len(updated))
	}

	batched, err := c.PatchSharesBatch([]Share{
		{ID: many[0].ID, MaxUses: 10},
		{ID: many[1].ID, MaxUses: 10},
	}, nil)
	if err != nil {
		t.Fatalf("PatchSharesBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchSharesBatch returned %d shares, want 2", len(batched))
	}

	if err := c.DeleteShares(keys); err != nil {
		t.Fatalf("DeleteShares: %v", err)
	}

	if err := c.DeleteShare(created.ID); err != nil {
		t.Fatalf("DeleteShare: %v", err)
	}
}

func TestE2EContentComments(t *testing.T) {
	c := itestClient(t)
	coll := newTestCollection(t, c)
	itemKey := newTestItem(t, c, coll, "e2e-comment")

	created, err := c.CreateComment(&Comment{Collection: coll, Item: itemKey, Comment: "hi"}, nil)
	if err != nil {
		t.Fatalf("CreateComment: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateComment returned no id")
	}
	t.Cleanup(func() { _ = c.DeleteComment(created.ID) })

	got, err := c.GetComment(created.ID, nil)
	if err != nil {
		t.Fatalf("GetComment: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetComment returned unexpected comment: %+v", got)
	}

	list, err := c.GetComments(nil)
	if err != nil {
		t.Fatalf("GetComments: %v", err)
	}
	if list == nil {
		t.Fatalf("GetComments returned nil")
	}

	patched, err := c.PatchComment(created.ID, &Comment{Comment: "updated"}, nil)
	if err != nil {
		t.Fatalf("PatchComment: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchComment returned unexpected comment: %+v", patched)
	}

	// Batch create.
	many, err := c.CreateComments([]Comment{
		{Collection: coll, Item: itemKey, Comment: "one"},
		{Collection: coll, Item: itemKey, Comment: "two"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateComments: %v", err)
	}
	if len(many) != 2 {
		t.Fatalf("CreateComments returned %d comments, want 2", len(many))
	}
	keys := []string{many[0].ID, many[1].ID}

	if err := c.DeleteComments(keys); err != nil {
		t.Fatalf("DeleteComments: %v", err)
	}

	if err := c.DeleteComment(created.ID); err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
}
