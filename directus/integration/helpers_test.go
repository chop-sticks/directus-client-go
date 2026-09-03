//go:build integration

package integration

import (
	"fmt"
	. "github.com/chop-sticks/directus-client-go/directus"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// e2eSeq guarantees unique fixture names within a single test run.
var e2eSeq atomic.Int64

// uniqueName returns a collision-free identifier with the given prefix.
func uniqueName(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().Unix(), e2eSeq.Add(1))
}

// newTestCollection creates a throwaway collection with an auto-increment
// integer primary key `id` and a string `title` field, registering cleanup.
func newTestCollection(t *testing.T, c *Client) string {
	t.Helper()
	name := uniqueName("e2e_col")
	_, err := c.CreateCollection(&CollectionRequest{
		Collection: name,
		Meta:       &CollectionMeta{Note: "e2e"},
		Schema:     &CollectionSchema{},
		Fields: []Field{
			{Field: "id", Type: "integer", Meta: &FieldMeta{Hidden: true}, Schema: &FieldSchema{IsPrimaryKey: true, HasAutoIncrement: true}},
			{Field: "title", Type: "string", Meta: &FieldMeta{Interface: "input"}, Schema: &FieldSchema{}},
		},
	}, nil)
	if err != nil {
		t.Fatalf("newTestCollection: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteCollection(name) })
	return name
}

// newTestItem inserts an item with the given title and returns its primary key.
func newTestItem(t *testing.T, c *Client, collection, title string) string {
	t.Helper()
	item, err := c.CreateItem(collection, map[string]any{"title": title}, nil)
	if err != nil {
		t.Fatalf("newTestItem: %v", err)
	}
	return fmt.Sprintf("%v", item["id"])
}

// newTestRole creates a throwaway role and returns its id.
func newTestRole(t *testing.T, c *Client) string {
	t.Helper()
	r, err := c.CreateRole(&Role{Name: uniqueName("e2e_role")}, nil)
	if err != nil {
		t.Fatalf("newTestRole: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteRole(r.ID) })
	return r.ID
}

// newTestPolicy creates a throwaway policy and returns its id.
func newTestPolicy(t *testing.T, c *Client) string {
	t.Helper()
	p, err := c.CreatePolicy(&Policy{Name: uniqueName("e2e_policy")}, nil)
	if err != nil {
		t.Fatalf("newTestPolicy: %v", err)
	}
	t.Cleanup(func() { _ = c.DeletePolicy(p.ID) })
	return p.ID
}

// newTestDashboard creates a throwaway dashboard and returns its id.
func newTestDashboard(t *testing.T, c *Client) string {
	t.Helper()
	d, err := c.CreateDashboard(&Dashboard{Name: uniqueName("e2e_dash"), Icon: "space_dashboard"}, nil)
	if err != nil {
		t.Fatalf("newTestDashboard: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteDashboard(d.ID) })
	return d.ID
}

// newTestFolder creates a throwaway folder and returns its id.
func newTestFolder(t *testing.T, c *Client) string {
	t.Helper()
	f, err := c.CreateFolder(&Folder{Name: uniqueName("e2e_folder")}, nil)
	if err != nil {
		t.Fatalf("newTestFolder: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteFolder(f.ID) })
	return f.ID
}

// newTestUser creates a throwaway user and returns its id.
func newTestUser(t *testing.T, c *Client) string {
	t.Helper()
	u, err := c.CreateUser(&User{
		Email:    uniqueName("e2e") + "@example.com",
		Password: "Testpassw0rd!",
	}, nil)
	if err != nil {
		t.Fatalf("newTestUser: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteUser(u.ID) })
	return u.ID
}

// uploadTestFile uploads a small text file and returns its id.
func uploadTestFile(t *testing.T, c *Client) string {
	t.Helper()
	f, err := c.UploadFile(strings.NewReader("e2e file contents"), "e2e.txt", map[string]string{"title": uniqueName("e2e_file")}, nil)
	if err != nil {
		t.Fatalf("uploadTestFile: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteFile(f.ID) })
	return f.ID
}

// TestE2EFixtures validates every shared fixture helper against live Directus.
func TestE2EFixtures(t *testing.T) {
	c := itestClient(t)

	coll := newTestCollection(t, c)
	if key := newTestItem(t, c, coll, "fixture"); key == "" {
		t.Error("newTestItem returned empty key")
	}
	if newTestRole(t, c) == "" {
		t.Error("newTestRole returned empty id")
	}
	if newTestPolicy(t, c) == "" {
		t.Error("newTestPolicy returned empty id")
	}
	if newTestDashboard(t, c) == "" {
		t.Error("newTestDashboard returned empty id")
	}
	if newTestFolder(t, c) == "" {
		t.Error("newTestFolder returned empty id")
	}
	if newTestUser(t, c) == "" {
		t.Error("newTestUser returned empty id")
	}
	if uploadTestFile(t, c) == "" {
		t.Error("uploadTestFile returned empty id")
	}
}
