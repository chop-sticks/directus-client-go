//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

func TestE2ETrackingNotifications(t *testing.T) {
	c := itestClient(t)
	userID := newTestUser(t, c)

	// CreateNotification
	created, err := c.CreateNotification(&Notification{Recipient: userID, Subject: "hi"}, nil)
	if err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}
	if created == nil || created.ID == 0 {
		t.Fatalf("CreateNotification: expected non-nil notification with id")
	}
	t.Cleanup(func() { _ = c.DeleteNotification(created.ID) })

	// GetNotification
	got, err := c.GetNotification(created.ID, nil)
	if err != nil {
		t.Fatalf("GetNotification: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetNotification: expected id %d, got %+v", created.ID, got)
	}

	// GetNotifications
	list, err := c.GetNotifications(nil)
	if err != nil {
		t.Fatalf("GetNotifications: %v", err)
	}
	if list == nil {
		t.Fatalf("GetNotifications: expected non-nil slice")
	}

	// PatchNotification
	patched, err := c.PatchNotification(created.ID, &Notification{Subject: "hi2"}, nil)
	if err != nil {
		t.Fatalf("PatchNotification: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchNotification: expected id %d, got %+v", created.ID, patched)
	}

	// CreateNotifications (batch)
	batch, err := c.CreateNotifications([]Notification{
		{Recipient: userID, Subject: "b1"},
		{Recipient: userID, Subject: "b2"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateNotifications: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreateNotifications: expected 2, got %d", len(batch))
	}
	keys := []int{batch[0].ID, batch[1].ID}

	// PatchNotifications (by keys)
	byKeys, err := c.PatchNotifications(keys, &Notification{Subject: "bk"}, nil)
	if err != nil {
		t.Fatalf("PatchNotifications: %v", err)
	}
	if len(byKeys) != 2 {
		t.Fatalf("PatchNotifications: expected 2, got %d", len(byKeys))
	}

	// PatchNotificationsBatch
	batched, err := c.PatchNotificationsBatch([]Notification{
		{ID: batch[0].ID, Subject: "bb1"},
		{ID: batch[1].ID, Subject: "bb2"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchNotificationsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchNotificationsBatch: expected 2, got %d", len(batched))
	}

	// DeleteNotifications (by keys)
	if err := c.DeleteNotifications(keys); err != nil {
		t.Fatalf("DeleteNotifications: %v", err)
	}

	// DeleteNotification
	if err := c.DeleteNotification(created.ID); err != nil {
		t.Fatalf("DeleteNotification: %v", err)
	}
}

func TestE2ETrackingActivity(t *testing.T) {
	c := itestClient(t)

	// generate some activity
	coll := newTestCollection(t, c)
	_ = newTestItem(t, c, coll, uniqueName("title"))

	// GetActivities
	activities, err := c.GetActivities(nil)
	if err != nil {
		t.Fatalf("GetActivities: %v", err)
	}
	if activities == nil {
		t.Fatalf("GetActivities: expected non-nil slice")
	}

	// GetActivity
	if len(activities) > 0 {
		act, err := c.GetActivity(activities[0].ID, nil)
		if err != nil {
			t.Fatalf("GetActivity: %v", err)
		}
		if act == nil || act.ID != activities[0].ID {
			t.Fatalf("GetActivity: expected id %d, got %+v", activities[0].ID, act)
		}
	}
}

func TestE2ETrackingRevisions(t *testing.T) {
	c := itestClient(t)

	// generate some revisions
	coll := newTestCollection(t, c)
	_ = newTestItem(t, c, coll, uniqueName("title"))

	// GetRevisions
	revisions, err := c.GetRevisions(nil)
	if err != nil {
		t.Fatalf("GetRevisions: %v", err)
	}
	if revisions == nil {
		t.Fatalf("GetRevisions: expected non-nil slice")
	}

	// GetRevision
	if len(revisions) > 0 {
		rev, err := c.GetRevision(revisions[0].ID, nil)
		if err != nil {
			t.Fatalf("GetRevision: %v", err)
		}
		if rev == nil || rev.ID != revisions[0].ID {
			t.Fatalf("GetRevision: expected id %d, got %+v", revisions[0].ID, rev)
		}
	}
}

func TestE2ETrackingVersions(t *testing.T) {
	c := itestClient(t)

	// collection with versioning enabled
	coll := newTestCollection(t, c)
	if _, err := c.PatchCollection(coll, &CollectionRequest{Meta: &CollectionMeta{Versioning: true}}, nil); err != nil {
		t.Fatalf("PatchCollection versioning: %v", err)
	}
	itemKey := newTestItem(t, c, coll, uniqueName("title"))

	// GetContentVersions
	versions, err := c.GetContentVersions(nil)
	if err != nil {
		t.Fatalf("GetContentVersions: %v", err)
	}
	if versions == nil {
		t.Fatalf("GetContentVersions: expected non-nil slice")
	}

	// CreateContentVersion
	created, err := c.CreateContentVersion(&Version{Key: uniqueName("e2e_ver"), Collection: coll, Item: itemKey}, nil)
	if err != nil {
		t.Fatalf("CreateContentVersion: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateContentVersion: expected non-nil version with id")
	}
	t.Cleanup(func() { _ = c.DeleteContentVersion(created.ID) })

	// GetContentVersion
	got, err := c.GetContentVersion(created.ID, nil)
	if err != nil {
		t.Fatalf("GetContentVersion: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetContentVersion: expected id %s, got %+v", created.ID, got)
	}

	// CreateContentVersions (batch)
	batch, err := c.CreateContentVersions([]Version{
		{Key: uniqueName("e2e_ver"), Collection: coll, Item: itemKey},
		{Key: uniqueName("e2e_ver"), Collection: coll, Item: itemKey},
	}, nil)
	if err != nil {
		t.Fatalf("CreateContentVersions: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreateContentVersions: expected 2, got %d", len(batch))
	}
	keys := []string{batch[0].ID, batch[1].ID}
	t.Cleanup(func() { _ = c.DeleteContentVersions(keys) })

	// PatchContentVersion
	patched, err := c.PatchContentVersion(created.ID, &Version{Name: "renamed"}, nil)
	if err != nil {
		t.Fatalf("PatchContentVersion: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchContentVersion: expected id %s, got %+v", created.ID, patched)
	}

	// PatchContentVersions (by keys)
	byKeys, err := c.PatchContentVersions(keys, &Version{Name: "renamed2"}, nil)
	if err != nil {
		t.Fatalf("PatchContentVersions: %v", err)
	}
	if len(byKeys) != 2 {
		t.Fatalf("PatchContentVersions: expected 2, got %d", len(byKeys))
	}

	// PatchContentVersionsBatch
	batched, err := c.PatchContentVersionsBatch([]Version{
		{ID: batch[0].ID, Name: "bb1"},
		{ID: batch[1].ID, Name: "bb2"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchContentVersionsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchContentVersionsBatch: expected 2, got %d", len(batched))
	}

	// SaveToContentVersion
	saved, err := c.SaveToContentVersion(created.ID, map[string]any{"title": "v2"})
	if err != nil {
		t.Fatalf("SaveToContentVersion: %v", err)
	}
	if saved == nil {
		t.Fatalf("SaveToContentVersion: expected non-nil result")
	}

	// CompareContentVersion
	compare, err := c.CompareContentVersion(created.ID)
	if err != nil {
		t.Fatalf("CompareContentVersion: %v", err)
	}
	if compare == nil {
		t.Fatalf("CompareContentVersion: expected non-nil result")
	}

	// PromoteContentVersion (environment/state dependent -> log only)
	if _, err := c.PromoteContentVersion(created.ID, compare.MainHash, nil); err != nil {
		t.Logf("PromoteContentVersion: %v", err)
	}

	// DeleteContentVersions (by keys)
	if err := c.DeleteContentVersions(keys); err != nil {
		t.Fatalf("DeleteContentVersions: %v", err)
	}

	// DeleteContentVersion
	if err := c.DeleteContentVersion(created.ID); err != nil {
		t.Fatalf("DeleteContentVersion: %v", err)
	}
}
