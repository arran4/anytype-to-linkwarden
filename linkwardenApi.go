package anytype_to_linkwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Tag struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	OwnerId   int    `json:"ownerId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Count     struct {
		Links int `json:"links"`
	} `json:"_count"`
}

func GetTags(token string, baseUrl string) ([]*Tag, error) {
	var response *http.Response
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tags", baseUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var responseObject struct {
		Tags []*Tag `json:"response"`
	}
	if err := json.Unmarshal(b, &responseObject); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return responseObject.Tags, nil
}

type Collection struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
	Members     any    `json:"members"`
	Parent      any    `json:"parent"`
	OwnerId     int    `json:"ownerId"`
	ParentId    int    `json:"parentId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func GetCollections(token string, baseUrl string) ([]*Collection, error) {
	var response *http.Response
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/collections", baseUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var responseObject struct {
		Collections []*Collection `json:"response"`
	}
	if err := json.Unmarshal(b, &responseObject); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return responseObject.Collections, nil
}

type PartialCreateCollection struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
	ParentId    int    `json:"parentId"`
}

func CreateCollections(token string, baseUrl string, collection *PartialCreateCollection) (*Collection, error) {
	var response *http.Response
	b, err := json.Marshal(collection)
	if err != nil {
		return nil, fmt.Errorf("unmarshal collection: %w", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections", baseUrl), bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	b, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var responseObject struct {
		Collection *Collection `json:"response"`
	}
	if err := json.Unmarshal(b, &responseObject); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return responseObject.Collection, nil
}

type Link struct {
	Id            int         `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Url           string      `json:"url"`
	Collection    *Collection `json:"collection"`
	Tags          []*Tag      `json:"tags"`
	Parent        any         `json:"parent"`
	OwnerId       int         `json:"ownerId"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
	Type          string      `json:"type"`
	CollectionId  int         `json:"collectionId"`
	TextContent   any         `json:"textContent"`
	Preview       any         `json:"preview"`
	Image         any         `json:"image"`
	Pdf           any         `json:"pdf"`
	Readable      any         `json:"readable"`
	Monolith      any         `json:"monolith"`
	LastPreserved any         `json:"lastPreserved"`
	ImportDate    any         `json:"importDate"`
}

type CollectionReference struct {
	Id      *int   `json:"id,omitempty"`
	OwnerId *int   `json:"ownerId,omitempty"`
	Name    string `json:"name"`
}

type TagReference struct {
	Id   *int   `json:"id,omitempty"`
	Name string `json:"name"`
}

type PartialCreateLink struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Url         string               `json:"url"`
	Collection  *CollectionReference `json:"collection"`
	Tags        []*TagReference      `json:"tags"`
}

func PostLink(token string, baseUrl string, link *PartialCreateLink) (*Link, error) {
	var response *http.Response
	b, err := json.Marshal(link)
	if err != nil {
		return nil, fmt.Errorf("unmarshal link: %w", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/links", baseUrl), bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	b, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var responseObject struct {
		Link *Link `json:"response"`
	}
	if err := json.Unmarshal(b, &responseObject); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return responseObject.Link, nil
}
