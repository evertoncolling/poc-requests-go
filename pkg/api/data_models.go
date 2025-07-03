package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/evertoncolling/poc-requests-go/pkg/dto"
)

func (d *DataModeling) ListDataModels(
	limit int,
	cursor *string,
	// inlineViews bool,
	space *string,
	allVersions bool,
	includeGlobal bool,
) (dto.DataModelList, error) {
	// Create query parameters
	queryParams := make(map[string]interface{})
	if cursor != nil {
		queryParams["cursor"] = *cursor
	}
	if space != nil {
		queryParams["space"] = *space
	}
	queryParams["limit"] = limit
	// queryParams["inlineViews"] = inlineViews
	queryParams["allVersions"] = allVersions
	queryParams["includeGlobal"] = includeGlobal

	endpoint := fmt.Sprintf("/api/v1/projects/%s/models/datamodels", d.Client.ClientConfig.Project)

	// Build the URL with query parameters
	url := d.Client.BaseURL + endpoint + "?" + buildQueryParams(queryParams)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.DataModelList{}, err
	}

	for key, value := range d.Client.Headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return dto.DataModelList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.DataModelList{}, fmt.Errorf("failed to fetch units: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.DataModelList{}, err
	}

	var dataModelsList dto.DataModelList
	if err := json.Unmarshal(body, &dataModelsList); err != nil {
		return dto.DataModelList{}, err
	}

	return dataModelsList, nil
}

func (d *DataModeling) InstancesSearch(
	view dto.ViewReference,
	query string,
	instanceType *string,
	properties *[]string,
	targetUnits *[]dto.TargetUnitsDM,
	filter *map[string]interface{},
	// includeTyping bool,
	sort *[]dto.SearchSort,
	limit int,
) (dto.NodeList, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/models/instances/search", d.Client.ClientConfig.Project)
	url := d.Client.BaseURL + endpoint

	body := make(map[string]interface{})
	body["view"] = view
	body["query"] = query
	if instanceType != nil {
		body["instanceType"] = *instanceType
	}
	if properties != nil {
		body["properties"] = *properties
	}
	if targetUnits != nil {
		body["targetUnits"] = *targetUnits
	}
	if filter != nil {
		body["filter"] = *filter
	}
	// body["includeTyping"] = includeTyping
	if sort != nil {
		body["sort"] = *sort
	}
	body["limit"] = limit

	// Convert the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return dto.NodeList{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return dto.NodeList{}, err
	}

	for key, value := range d.Client.Headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return dto.NodeList{}, err
	}
	defer resp.Body.Close()

	// Check if status is not OK
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error message
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return dto.NodeList{}, fmt.Errorf("failed to search instances: %s - error reading response body: %v", resp.Status, err)
		}
		return dto.NodeList{}, fmt.Errorf("failed to search instances: %s - %s", resp.Status, string(respBody))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.NodeList{}, err
	}

	var nodeList dto.NodeList
	if err := json.Unmarshal(responseBody, &nodeList); err != nil {
		return dto.NodeList{}, err
	}

	return nodeList, nil
}

func (d *DataModeling) GraphQLQuery(
	space string,
	externalId string,
	version string,
	query string,
	variables map[string]interface{},
) (dto.GraphQLResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/userapis/spaces/%s/datamodels/%s/versions/%s/graphql",
		d.Client.ClientConfig.Project, space, externalId, version)
	url := d.Client.BaseURL + endpoint

	request := dto.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return dto.GraphQLResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return dto.GraphQLResponse{}, err
	}

	for key, value := range d.Client.Headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return dto.GraphQLResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.GraphQLResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return dto.GraphQLResponse{}, fmt.Errorf("GraphQL query failed: %s - %s", resp.Status, string(responseBody))
	}

	var graphQLResponse dto.GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphQLResponse); err != nil {
		return dto.GraphQLResponse{}, err
	}

	return graphQLResponse, nil
}
