package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/evertoncolling/poc-requests-go/pkg/dto"

	"google.golang.org/protobuf/proto"
)

func ListTimeSeries(
	project string,
	token string,
	baseURL string,
	limit int,
	includeMetadata bool,
	cursor string,
	partition string,
	assetIDs []int64,
	rootAssetIDs []int64,
	externalIDPrefix string,
) (dto.TimeSeriesList, error) {
	// Create query parameters
	queryParams := make(map[string]interface{})
	queryParams["limit"] = limit
	queryParams["includeMetadata"] = includeMetadata
	queryParams["cursor"] = cursor
	queryParams["partition"] = partition
	queryParams["assetIds"] = assetIDs
	queryParams["rootAssetIds"] = rootAssetIDs
	queryParams["externalIdPrefix"] = externalIDPrefix

	endpoint := fmt.Sprintf("/api/v1/projects/%s/timeseries", project)

	// Build the URL with query parameters
	url := baseURL + endpoint + "?" + buildQueryParams(queryParams)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("cdf-version", "beta")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.TimeSeriesList{}, fmt.Errorf("failed to fetch timeseries: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}

	var tsList dto.TimeSeriesList
	if err := json.Unmarshal(body, &tsList); err != nil {
		return dto.TimeSeriesList{}, err
	}

	return tsList, nil
}

func FilterTimeSeries(
	project string,
	token string,
	baseURL string,
	filter *dto.TimeSeriesFilter,
	advancedFilter map[string]interface{},
	limit int,
	cursor string,
	partition string,
	sort []map[string]interface{},
) (dto.TimeSeriesList, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/timeseries/list", project)
	url := baseURL + endpoint

	body := map[string]interface{}{
		"filter":         filter,
		"advancedFilter": advancedFilter,
		"limit":          limit,
		"cursor":         cursor,
		"partition":      partition,
		"sort":           sort,
	}

	// Convert the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return dto.TimeSeriesList{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("cdf-version", "beta")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.TimeSeriesList{}, fmt.Errorf("failed to fetch timeseries: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.TimeSeriesList{}, err
	}

	var tsList dto.TimeSeriesList
	if err := json.Unmarshal(responseBody, &tsList); err != nil {
		return dto.TimeSeriesList{}, err
	}

	return tsList, nil
}

func RetrieveData(
	project string,
	token string,
	baseURL string,
	items *[]dto.DataPointsQueryItem,
	startTime *string,
	endTime *string,
	limit *int64,
	aggregates *[]string,
	granularity *string,
	includeOutsidePoints *bool,
	ignoreUnknownIds *bool,
) (*dto.DataPointListResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/timeseries/data/list", project)
	url := baseURL + endpoint

	body := make(map[string]interface{})
	body["items"] = items

	// Only add to body if parameters are not nil
	if startTime != nil {
		body["start"] = startTime
	}
	if endTime != nil {
		body["end"] = endTime
	}
	if limit != nil {
		body["limit"] = limit
	}
	if aggregates != nil {
		body["aggregates"] = aggregates
	}
	if granularity != nil {
		body["granularity"] = granularity
	}
	if includeOutsidePoints != nil {
		body["includeOutsidePoints"] = includeOutsidePoints
	}
	if ignoreUnknownIds != nil {
		body["ignoreUnknownIds"] = ignoreUnknownIds
	}

	// Convert the body to JSON
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Compress the JSON body (increase performance slightly)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(bodyJSON); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/protobuf")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("cdf-version", "beta")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch datapoints: %s", resp.Status)
	}

	// Read and decode the protobuf response
	var responseBody bytes.Buffer
	if _, err := io.Copy(&responseBody, resp.Body); err != nil {
		return nil, err
	}

	var dpList dto.DataPointListResponse
	if err := proto.Unmarshal(responseBody.Bytes(), &dpList); err != nil {
		return nil, err
	}

	return &dpList, nil
}
