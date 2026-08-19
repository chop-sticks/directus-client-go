package directus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FieldsResponse struct {
	Data []Field `json:"data"`
}

type FieldResponse struct {
	Data Field `json:"data"`
}

func (c *Client) GetFields() ([]Field, error) {
	return c.GetFieldsByCollection("")
}

func (c *Client) GetFieldsByCollection(collection string) ([]Field, error) {
	urlPath := buildURL(c.HostURL, fmt.Sprintf("/fields/%s", collection), nil)
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var fields *FieldsResponse
	err = json.Unmarshal(body, &fields)
	if err != nil {
		return nil, err
	}
	return fields.Data, nil
}

func (c *Client) GetFieldByCollectionAndName(collection string, name string) (*Field, error) {
	if collection == "" || name == "" {
		return nil, fmt.Errorf("collection and name must be provided")
	}
	urlPath := buildURL(c.HostURL, fmt.Sprintf("/fields/%s/%s", collection, name), nil)
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var field *FieldResponse
	err = json.Unmarshal(body, &field)
	if err != nil {
		return nil, err
	}
	return &field.Data, nil
}
