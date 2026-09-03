package directus

import "fmt"

// GetAsset retrieves the raw bytes of a file's asset by id. Query parameters
// (such as transformation presets) are passed through via q.
func (c *Client) GetAsset(id string, q *Query) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return c.requestRaw("GET", fmt.Sprintf("/assets/%s", id), q, nil)
}

// DownloadFilesZip downloads the given files as a single ZIP archive and
// returns its raw bytes.
func (c *Client) DownloadFilesZip(ids []string) ([]byte, error) {
	body := map[string]any{"ids": ids}
	return c.requestRaw("POST", "/assets/files/", nil, body)
}

// DownloadFolderZip downloads the contents of a folder as a ZIP archive and
// returns its raw bytes.
func (c *Client) DownloadFolderZip(id string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return c.requestRaw("POST", fmt.Sprintf("/assets/folder/%s", id), nil, nil)
}
