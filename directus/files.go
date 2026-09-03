package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// File models a Directus file (directus_files). The id is a UUID string.
type File struct {
	ID               string         `json:"id,omitempty"`
	Storage          string         `json:"storage,omitempty"`
	FilenameDisk     string         `json:"filename_disk,omitempty"`
	FilenameDownload string         `json:"filename_download,omitempty"`
	Title            string         `json:"title,omitempty"`
	Type             string         `json:"type,omitempty"`
	Folder           any            `json:"folder,omitempty"`
	CreatedOn        string         `json:"created_on,omitempty"`
	UploadedBy       any            `json:"uploaded_by,omitempty"`
	UploadedOn       string         `json:"uploaded_on,omitempty"`
	ModifiedBy       any            `json:"modified_by,omitempty"`
	ModifiedOn       string         `json:"modified_on,omitempty"`
	Charset          string         `json:"charset,omitempty"`
	Filesize         string         `json:"filesize,omitempty"`
	Width            int            `json:"width,omitempty"`
	Height           int            `json:"height,omitempty"`
	Duration         int            `json:"duration,omitempty"`
	Embed            string         `json:"embed,omitempty"`
	Description      string         `json:"description,omitempty"`
	Location         string         `json:"location,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	FocalPointX      int            `json:"focal_point_x,omitempty"`
	FocalPointY      int            `json:"focal_point_y,omitempty"`
	TusID            string         `json:"tus_id,omitempty"`
	TusData          map[string]any `json:"tus_data,omitempty"`
}

// GetFiles lists files.
func (c *Client) GetFiles(q *Query) ([]File, error) {
	return request[[]File](c, "GET", "/files", q, nil)
}

// GetFile retrieves a single file by id.
func (c *Client) GetFile(id string, q *Query) (*File, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*File](c, "GET", fmt.Sprintf("/files/%s", id), q, nil)
}

// ImportFile imports a file from a remote URL. Optional metadata is supplied
// via data.
func (c *Client) ImportFile(url string, data *File, q *Query) (*File, error) {
	if url == "" {
		return nil, fmt.Errorf("url must be provided")
	}
	body := map[string]any{"url": url}
	if data != nil {
		body["data"] = data
	}
	return request[*File](c, "POST", "/files/import", q, body)
}

// PatchFile updates a single file's metadata.
func (c *Client) PatchFile(id string, item *File, q *Query) (*File, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*File](c, "PATCH", fmt.Sprintf("/files/%s", id), q, item)
}

// PatchFiles updates multiple files identified by keys with the same data.
func (c *Client) PatchFiles(keys []string, item *File, q *Query) ([]File, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]File](c, "PATCH", "/files", q, body)
}

// PatchFilesBatch updates multiple files, each with its own values.
func (c *Client) PatchFilesBatch(items []File, q *Query) ([]File, error) {
	return request[[]File](c, "PATCH", "/files", q, items)
}

// DeleteFile deletes a single file by id.
func (c *Client) DeleteFile(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/files/%s", id), nil, nil)
}

// DeleteFiles deletes multiple files identified by keys.
func (c *Client) DeleteFiles(keys []string) error {
	return c.execute("DELETE", "/files", nil, keys)
}

// UploadFile uploads a new file using a multipart/form-data request. Each entry
// of fields is added as a form field before the binary "file" part, whose bytes
// are copied from data.
//
// This is the one allowed exception to the "no hand-rolled http.NewRequest"
// rule: multipart uploads require the writer-generated boundary in the
// Content-Type header, which the shared helpers cannot express. The request is
// still sent through c.doRequest and the {"data": ...} envelope is decoded here.
func (c *Client) UploadFile(data io.Reader, filename string, fields map[string]string, q *Query) (*File, error) {
	if data == nil {
		return nil, fmt.Errorf("data must be provided")
	}
	if filename == "" {
		return nil, fmt.Errorf("filename must be provided")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", buildURL(c.HostURL, "/files", q), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	raw, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	var resp apiResponse[*File]
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
