package api

import (
	"os"
	"testing"

	"github.com/evertoncolling/poc-requests-go/pkg/dto"
)

// setupIntegrationTest loads environment variables and creates a real client
func setupIntegrationTest(t *testing.T) CogniteClient {
	t.Helper()

	// Check if required environment variables are set
	// Note: For local development, load .env before running tests with: source .env && go test
	// For CI, environment variables are set via GitHub secrets
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tenantID := os.Getenv("TENANT_ID")
	cluster := os.Getenv("CDF_CLUSTER")
	project := os.Getenv("CDF_PROJECT")

	if clientID == "" || clientSecret == "" || tenantID == "" || cluster == "" || project == "" {
		t.Fatal("Integration test failed: required environment variables not set (CLIENT_ID, CLIENT_SECRET, TENANT_ID, CDF_CLUSTER, CDF_PROJECT)")
	}

	credentials := AzureADClientCredentials(
		clientID,
		clientSecret,
		tenantID,
		cluster,
	)

	clientConfig := ClientConfig{
		ClientName:  "poc-requests-go-integration-test",
		Cluster:     cluster,
		Project:     project,
		Credentials: credentials,
	}

	return NewCogniteClient(clientConfig)
}

func TestIntegration_TimeSeries_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test listing time series
	tsList, err := client.TimeSeries.List(5, false, "", "", nil, nil, "")
	if err != nil {
		t.Fatalf("Failed to list time series: %v", err)
	}

	t.Logf("Successfully fetched %d time series", len(tsList.Items))

	// Verify we got some results
	if len(tsList.Items) == 0 {
		t.Log("Warning: No time series found - this might be expected for empty projects")
	}

	// Verify the structure of returned data
	for i, ts := range tsList.Items {
		if i >= 3 { // Only check first 3 to avoid too much logging
			break
		}
		t.Logf("Time series %d: ID=%d, ExternalID=%s", i+1, ts.Id, ts.ExternalId)

		if ts.Id == 0 && ts.ExternalId == "" {
			t.Errorf("Time series %d has both empty ID and ExternalID", i+1)
		}
	}
}

func TestIntegration_TimeSeries_Filter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test filtering time series with a simple filter
	filter := &dto.TimeSeriesFilter{
		UnitQuantity: "Temperature",
	}

	filteredTsList, err := client.TimeSeries.Filter(filter, nil, 3, "", "", nil)
	if err != nil {
		t.Fatalf("Failed to filter time series: %v", err)
	}

	t.Logf("Successfully filtered time series, found %d results", len(filteredTsList.Items))

	for i, ts := range filteredTsList.Items {
		if i >= 2 { // Only check first 2
			break
		}
		t.Logf("Filtered time series %d: ID=%d, ExternalID=%s", i+1, ts.Id, ts.ExternalId)
	}
}

func TestIntegration_Units_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test fetching the unit catalog
	unitList, err := client.Units.List()
	if err != nil {
		t.Fatalf("Failed to list units: %v", err)
	}

	t.Logf("Successfully fetched %d units", len(unitList.Items))

	// Units are available in all CDF projects
	if len(unitList.Items) == 0 {
		t.Error("Expected at least some units in the catalog")
	}

	// Verify the structure of returned data
	foundTemperatureUnit := false
	for i, unit := range unitList.Items {
		if i >= 5 { // Only check first 5 to avoid too much logging
			break
		}
		t.Logf("Unit %d: ExternalID=%s, Name=%s", i+1, unit.ExternalId, unit.Name)

		if unit.ExternalId == "" {
			t.Errorf("Unit %d has empty ExternalID", i+1)
		}

		// Look for a common unit like temperature
		if unit.ExternalId == "temperature:deg_c" {
			foundTemperatureUnit = true
		}
	}

	if !foundTemperatureUnit && len(unitList.Items) > 0 {
		t.Error("No temperature unit found")
	}
}

func TestIntegration_DataModeling_ListDataModels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test listing data models
	dataModelsList, err := client.DataModeling.ListDataModels(
		5,     // limit
		nil,   // cursor
		nil,   // space
		false, // includeGlobal
		true,  // allVersions
	)
	if err != nil {
		t.Fatalf("Failed to list data models: %v", err)
	}

	t.Logf("Successfully fetched %d data models", len(dataModelsList.Items))

	if len(dataModelsList.Items) == 0 {
		t.Log("No data models found - this might be expected for projects without data modeling")
		return
	}

	// Verify the structure of returned data
	for i, model := range dataModelsList.Items {
		if i >= 3 { // Only check first 3
			break
		}
		t.Logf("Data model %d: Space=%s, ExternalID=%s, Version=%s",
			i+1, model.Space, model.ExternalId, model.Version)

		if model.Space == "" || model.ExternalId == "" {
			t.Errorf("Data model %d has empty Space or ExternalID", i+1)
		}
	}
}

func TestIntegration_DataModeling_InstancesSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test searching for CogniteTimeSeries instances
	properties := []string{"name", "description"}
	nodeList, err := client.DataModeling.InstancesSearch(
		dto.ViewReference{
			Type:       "view",
			Space:      "cdf_cdm",
			ExternalId: "CogniteTimeSeries",
			Version:    "v1",
		},
		"",          // query
		nil,         // filter
		&properties, // properties
		nil,         // sort
		nil,         // limit (will use default)
		nil,         // cursor
		3,           // limit
	)
	if err != nil {
		t.Logf("InstancesSearch failed: %v", err)
		return
	}

	t.Logf("Successfully searched instances, found %d results", len(nodeList.Items))

	// Verify the structure if we got results
	for i, node := range nodeList.Items {
		if i >= 2 { // Only check first 2
			break
		}
		t.Logf("Instance %d: Space=%s, ExternalID=%s", i+1, node.Space, node.ExternalId)

		if node.Space == "" || node.ExternalId == "" {
			t.Errorf("Instance %d has empty Space or ExternalID", i+1)
		}
	}
}

func TestIntegration_DataModeling_GraphQLQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test a simple GraphQL query for units
	graphQLQuery := `query MyQuery($nItems: Int) {
  listCogniteUnit(filter: {space: {eq: "cdf_cdm_units"}}, first: $nItems) {
    items {
      externalId
      name
      symbol
      quantity
    }
  }
}`

	variables := map[string]interface{}{
		"nItems": 3,
	}

	graphQLResponse, err := client.DataModeling.GraphQLQuery(
		"cdf_cdm",
		"CogniteCore",
		"v1",
		graphQLQuery,
		variables,
	)
	if err != nil {
		t.Logf("GraphQL query failed: %v", err)
		return
	}

	// Check for GraphQL errors
	if len(graphQLResponse.Errors) > 0 {
		t.Logf("GraphQL returned errors:")
		for _, gqlErr := range graphQLResponse.Errors {
			t.Logf("  - %s", gqlErr.Message)
		}
		return
	}

	t.Log("GraphQL query executed successfully")

	// Verify we got some data
	if graphQLResponse.Data == nil {
		t.Error("GraphQL response data is nil")
	} else {
		t.Logf("GraphQL response received (length: %d bytes)", len(graphQLResponse.Data))
	}
}

func TestIntegration_Authentication_TokenFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// Test that we can fetch a token
	if client.AccessToken == "" {
		t.Error("Access token is empty after client creation")
	} else {
		t.Logf("Successfully obtained access token (length: %d characters)", len(client.AccessToken))
	}

	// Verify client configuration
	if client.ClientConfig.Project == "" {
		t.Error("Client project is empty")
	} else {
		t.Logf("Client configured for project: %s", client.ClientConfig.Project)
	}

	if client.ClientConfig.Cluster == "" {
		t.Error("Client cluster is empty")
	} else {
		t.Logf("Client configured for cluster: %s", client.ClientConfig.Cluster)
	}
}

// Helper function to test actual data retrieval if we have time series
func TestIntegration_TimeSeries_RetrieveData_IfAvailable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := setupIntegrationTest(t)

	// First, try to get a time series
	tsList, err := client.TimeSeries.List(1, false, "", "", nil, nil, "")
	if err != nil {
		t.Fatalf("Failed to list time series: %v", err)
	}

	if len(tsList.Items) == 0 {
		t.Skip("No time series available for data retrieval test")
	}

	// Get the first time series
	ts := tsList.Items[0]
	var externalId string
	if ts.ExternalId != "" {
		externalId = ts.ExternalId
	} else {
		t.Skip("Time series has no external ID, cannot test data retrieval")
	}

	t.Logf("Testing data retrieval for time series: %s", externalId)

	// Try to get latest data point
	latestDataPointsQueryItems := []dto.LatestDataPointsQueryItem{
		{
			ExternalId: externalId,
		},
	}

	latestDataPoints, err := client.TimeSeries.RetrieveLatest(
		&latestDataPointsQueryItems,
		nil,
	)
	if err != nil {
		t.Logf("Failed to retrieve latest data points: %v", err)
		return
	}

	t.Logf("Successfully retrieved latest data points for %d time series", len(latestDataPoints.Items))

	for i, latestDataPoint := range latestDataPoints.Items {
		if i >= 1 { // Only check first one
			break
		}
		t.Logf("Latest data point for %s: type %T", latestDataPoint.ExternalId, latestDataPoint.DatapointType)
	}
}
