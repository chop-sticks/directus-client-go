package directus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CollectionsResponse struct {
	Data []Collection `json:"data"`
}

func (c *Client) GetCollections() ([]Collection, error) {
	urlPath := buildURL(c.HostURL, "/collections", nil)
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var collections *CollectionsResponse
	err = json.Unmarshal(body, &collections)
	if err != nil {
		return nil, err
	}
	return collections.Data, nil
}

type CollectionResponse struct {
	Data Collection `json:"data"`
}

func (c *Client) GetCollectionByName(name string) (*Collection, error) {
	urlPath := buildURL(c.HostURL, fmt.Sprintf("/collections/%s", name), nil)
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var collection *CollectionResponse
	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}
	return &collection.Data, nil
}

func (c *Client) processCollection(name string, req *CollectionRequest) (*Collection, error) {
	rb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	urlPath := buildURL(c.HostURL, fmt.Sprintf("/collections/%s", name), nil)
	method := http.MethodPost
	if name != "" {
		method = http.MethodPatch
	}
	httpReq, err := http.NewRequest(method, urlPath, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	var collection *CollectionResponse
	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}

	return &collection.Data, nil
}

func (c *Client) CreateCollection(req *CollectionRequest) (*Collection, error) {
	return c.processCollection("", req)
}

func (c *Client) PatchCollection(name string, req *CollectionRequest) (*Collection, error) {
	return c.processCollection(name, req)
}
