//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"strings"
	"testing"
)

func TestE2EFilesCRUD(t *testing.T) {
	c := itestClient(t)

	// UploadFile (directly, so we control the created id and cleanup).
	uploaded, err := c.UploadFile(strings.NewReader("e2e file contents"), "e2e.txt", map[string]string{"title": uniqueName("e2e_file")}, nil)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	if uploaded == nil || uploaded.ID == "" {
		t.Fatalf("UploadFile returned no id")
	}
	t.Cleanup(func() { _ = c.DeleteFile(uploaded.ID) })

	// GetFiles.
	files, err := c.GetFiles(nil)
	if err != nil {
		t.Fatalf("GetFiles: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("GetFiles returned no files")
	}

	// GetFile.
	got, err := c.GetFile(uploaded.ID, nil)
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	if got == nil || got.ID != uploaded.ID {
		t.Fatalf("GetFile returned %+v, want id %s", got, uploaded.ID)
	}

	// PatchFile (title).
	newTitle := uniqueName("e2e_title")
	patched, err := c.PatchFile(uploaded.ID, &File{Title: newTitle}, nil)
	if err != nil {
		t.Fatalf("PatchFile: %v", err)
	}
	if patched == nil || patched.Title != newTitle {
		t.Fatalf("PatchFile title = %q, want %q", patched.Title, newTitle)
	}
}

func TestE2EFilesBatch(t *testing.T) {
	c := itestClient(t)

	f1 := uploadTestFile(t, c)
	f2 := uploadTestFile(t, c)

	// PatchFiles (keys + shared data).
	sharedTitle := uniqueName("e2e_shared")
	updated, err := c.PatchFiles([]string{f1, f2}, &File{Title: sharedTitle}, nil)
	if err != nil {
		t.Fatalf("PatchFiles: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchFiles returned %d files, want 2", len(updated))
	}

	// PatchFilesBatch (per-item data).
	batch := []File{
		{ID: f1, Title: uniqueName("e2e_b1")},
		{ID: f2, Title: uniqueName("e2e_b2")},
	}
	batched, err := c.PatchFilesBatch(batch, nil)
	if err != nil {
		t.Fatalf("PatchFilesBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchFilesBatch returned %d files, want 2", len(batched))
	}
}

func TestE2EFilesDelete(t *testing.T) {
	c := itestClient(t)

	// DeleteFile (single).
	single, err := c.UploadFile(strings.NewReader("delete me"), "e2e_del.txt", map[string]string{"title": uniqueName("e2e_del")}, nil)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	if err := c.DeleteFile(single.ID); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}

	// DeleteFiles (many).
	a, err := c.UploadFile(strings.NewReader("a"), "e2e_a.txt", map[string]string{"title": uniqueName("e2e_a")}, nil)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	b, err := c.UploadFile(strings.NewReader("b"), "e2e_b.txt", map[string]string{"title": uniqueName("e2e_b")}, nil)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	if err := c.DeleteFiles([]string{a.ID, b.ID}); err != nil {
		t.Fatalf("DeleteFiles: %v", err)
	}
}

func TestE2EFilesImport(t *testing.T) {
	c := itestClient(t)

	// ImportFile requires external network access; log-only.
	f, err := c.ImportFile("https://raw.githubusercontent.com/directus/directus/main/readme.md", &File{Title: uniqueName("e2e_import")}, nil)
	if err != nil {
		t.Logf("ImportFile: %v", err)
		return
	}
	if f != nil && f.ID != "" {
		t.Cleanup(func() { _ = c.DeleteFile(f.ID) })
	}
}

func TestE2EFoldersCRUD(t *testing.T) {
	c := itestClient(t)

	// CreateFolder (single).
	single, err := c.CreateFolder(&Folder{Name: uniqueName("e2e_folder")}, nil)
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	if single == nil || single.ID == "" {
		t.Fatalf("CreateFolder returned no id")
	}
	t.Cleanup(func() { _ = c.DeleteFolder(single.ID) })

	// GetFolders.
	folders, err := c.GetFolders(nil)
	if err != nil {
		t.Fatalf("GetFolders: %v", err)
	}
	if len(folders) == 0 {
		t.Fatalf("GetFolders returned no folders")
	}

	// GetFolder.
	got, err := c.GetFolder(single.ID, nil)
	if err != nil {
		t.Fatalf("GetFolder: %v", err)
	}
	if got == nil || got.ID != single.ID {
		t.Fatalf("GetFolder returned %+v, want id %s", got, single.ID)
	}
}

func TestE2EFoldersDeleteMany(t *testing.T) {
	c := itestClient(t)

	// CreateFolders (many).
	created, err := c.CreateFolders([]Folder{
		{Name: uniqueName("e2e_folders_a")},
		{Name: uniqueName("e2e_folders_b")},
	}, nil)
	if err != nil {
		t.Fatalf("CreateFolders: %v", err)
	}
	if len(created) != 2 {
		t.Fatalf("CreateFolders returned %d folders, want 2", len(created))
	}

	// DeleteFolder (single).
	if err := c.DeleteFolder(created[0].ID); err != nil {
		t.Fatalf("DeleteFolder: %v", err)
	}

	// DeleteFolders (many).
	if err := c.DeleteFolders([]string{created[1].ID}); err != nil {
		t.Fatalf("DeleteFolders: %v", err)
	}
}

func TestE2EAssets(t *testing.T) {
	c := itestClient(t)

	fileID := uploadTestFile(t, c)

	// GetAsset.
	raw, err := c.GetAsset(fileID, nil)
	if err != nil {
		t.Fatalf("GetAsset: %v", err)
	}
	if len(raw) == 0 {
		t.Fatalf("GetAsset returned no bytes")
	}

	// DownloadFilesZip.
	zip, err := c.DownloadFilesZip([]string{fileID})
	if err != nil {
		t.Fatalf("DownloadFilesZip: %v", err)
	}
	if len(zip) == 0 {
		t.Fatalf("DownloadFilesZip returned no bytes")
	}
}

func TestE2EAssetsFolderZip(t *testing.T) {
	c := itestClient(t)

	folderID := newTestFolder(t, c)

	// Upload a file directly into the folder.
	f, err := c.UploadFile(strings.NewReader("zip me"), "e2e_zip.txt", map[string]string{"title": uniqueName("e2e_zip")}, nil)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteFile(f.ID) })

	if _, err := c.PatchFile(f.ID, &File{Folder: folderID}, nil); err != nil {
		t.Fatalf("PatchFile folder: %v", err)
	}

	// DownloadFolderZip (log-only; may error on some backends).
	zip, err := c.DownloadFolderZip(folderID)
	if err != nil {
		t.Logf("DownloadFolderZip: %v", err)
		return
	}
	if len(zip) == 0 {
		t.Logf("DownloadFolderZip returned no bytes")
	}
}

func TestE2EUtilsCache(t *testing.T) {
	c := itestClient(t)

	if err := c.ClearCache(false); err != nil {
		t.Fatalf("ClearCache(false): %v", err)
	}
	if err := c.ClearCache(true); err != nil {
		t.Fatalf("ClearCache(true): %v", err)
	}
}

func TestE2EUtilsRandomString(t *testing.T) {
	c := itestClient(t)

	s, err := c.RandomString(16)
	if err != nil {
		t.Fatalf("RandomString(16): %v", err)
	}
	if len(s) != 16 {
		t.Fatalf("RandomString(16) len = %d, want 16", len(s))
	}

	// RandomString(0) uses the server default length.
	zero, err := c.RandomString(0)
	if err != nil {
		t.Fatalf("RandomString(0): %v", err)
	}
	if len(zero) == 0 {
		t.Fatalf("RandomString(0) returned empty string")
	}
}

func TestE2EUtilsExportImport(t *testing.T) {
	c := itestClient(t)

	coll := newTestCollection(t, c)
	_ = newTestItem(t, c, coll, uniqueName("e2e_export"))

	// UtilsExport.
	if err := c.UtilsExport(coll, "json", map[string]any{}, map[string]any{"filename_download": "e2e.json"}); err != nil {
		t.Fatalf("UtilsExport: %v", err)
	}

	// UtilsImport.
	if err := c.UtilsImport(coll, strings.NewReader("title\nfromcsv\n"), "import.csv"); err != nil {
		t.Fatalf("UtilsImport: %v", err)
	}
}

func TestE2EUtilsSort(t *testing.T) {
	c := itestClient(t)

	coll := newTestCollection(t, c)

	// Add an integer "sort" field and mark it as the collection sort field.
	if _, err := c.CreateField(coll, &Field{Field: "sort", Type: "integer", Meta: &FieldMeta{}, Schema: &FieldSchema{IsNullable: true}}, nil); err != nil {
		t.Fatalf("CreateField sort: %v", err)
	}
	if _, err := c.PatchCollection(coll, &CollectionRequest{Meta: &CollectionMeta{SortField: "sort"}}, nil); err != nil {
		t.Fatalf("PatchCollection sort field: %v", err)
	}

	item1 := newTestItem(t, c, coll, uniqueName("e2e_sort1"))
	item2 := newTestItem(t, c, coll, uniqueName("e2e_sort2"))

	// UtilitySort (log-only; behaviour depends on backend state).
	if err := c.UtilitySort(coll, item2, item1); err != nil {
		t.Logf("UtilitySort: %v", err)
	}
}
