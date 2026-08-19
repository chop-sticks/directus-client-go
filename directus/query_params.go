package directus

import (
	"fmt"
	"net/url"
)

func buildQueryString(params interface{}) string {
	values := url.Values{}
	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

func buildURL(baseURL, path string, params interface{}) string {
	urlPath := fmt.Sprintf("%s%s", baseURL, path)
	queryString := buildQueryString(params)
	return urlPath + queryString
}
