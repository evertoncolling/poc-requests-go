package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/evertoncolling/poc-requests-go/pkg/api"
	"github.com/evertoncolling/poc-requests-go/pkg/dto"

	"github.com/joho/godotenv"
)

// Find a unit by external ID
func findUnitByExternalId(unitList *dto.UnitList, externalId string) (*dto.Unit, error) {
	for _, unit := range unitList.Items {
		if unit.ExternalId == externalId {
			return &unit, nil
		}
	}
	return nil, fmt.Errorf("unit not found with ExternalId: %s", externalId)
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read client credentials from environment variables
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tenantID := os.Getenv("TENANT_ID")
	cluster := os.Getenv("CDF_CLUSTER")
	project := os.Getenv("CDF_PROJECT")

	credentials := api.AzureADClientCredentials(
		clientID,
		clientSecret,
		tenantID,
		cluster,
	)
	clientConfig := api.ClientConfig{
		ClientName:  "poc-requests-go",
		Cluster:     cluster,
		Project:     project,
		Credentials: credentials,
	}
	client := api.NewCogniteClient(clientConfig)

	fmt.Println("### Testing fetching some time series")

	// List time series
	tsList, err := client.TimeSeries.List(100, false, "", "", nil, nil, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Time Series Count:", len(tsList.Items))

	// Filter time series
	filter := &dto.TimeSeriesFilter{
		UnitQuantity: "Pressure",
	}
	filteredTsList, err := client.TimeSeries.Filter(filter, nil, 100, "", "", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Filtered Time Series Count:", len(filteredTsList.Items))

	// Fetch the unit catalog
	fmt.Println("\n### Testing fetching the Unit catalog")
	unitList, err := client.Units.List()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Unit Count:", len(unitList.Items))

	// Fetch the latest data points
	fmt.Println("\n### Latest Data Points")
	externalIds := []string{"EVE-TI-FORNEBU-01-2", "EVE-TI-FORNEBU-01-3"}
	var latestDataPointsQueryItems []dto.LatestDataPointsQueryItem
	for _, externalId := range externalIds {
		latestDataPointsQueryItems = append(latestDataPointsQueryItems, dto.LatestDataPointsQueryItem{
			ExternalId: externalId,
		})
	}
	latestDataPoints, err := client.TimeSeries.RetrieveLatest(
		&latestDataPointsQueryItems,
		nil,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, latestDataPoint := range latestDataPoints.Items {
		fmt.Println(latestDataPoint.ExternalId+":", latestDataPoint.DatapointType)
	}

	// Get many data points (performance test)
	fmt.Println("\n### Data Points (performance test)")
	items := []dto.DataPointsQueryItem{
		{
			ExternalId: "EVE-TI-FORNEBU-01-2",
			Start:      "300d-ago",
			End:        "now",
			Limit:      100000,
		},
	}
	start := time.Now()
	dataPoints, err := client.TimeSeries.RetrieveData(
		&items,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dps := dataPoints.Items[0]
	fmt.Println("Data Points External ID:", dps.ExternalId)
	fmt.Println("Data Points Unit:", dps.UnitExternalId)
	// not sure if there is a simpler way to get to the data points :)
	switch v := (*dps).DatapointType.(type) {
	case *dto.DataPointListItem_NumericDatapoints:
		fmt.Println("Data Points Count:", len(v.NumericDatapoints.Datapoints))
	case *dto.DataPointListItem_StringDatapoints:
		fmt.Println("Data Points Count:", len(v.StringDatapoints.Datapoints))
	case *dto.DataPointListItem_AggregateDatapoints:
		fmt.Println("Data Points Count:", len(v.AggregateDatapoints.Datapoints))
	default:
		fmt.Println("Unknown data point type:", v)
	}
	fmt.Printf("Time taken: %s\n", elapsed)

	// Fetch data models
	fmt.Println("\n### Testing fetching some data models")
	dataModelsList, err := client.DataModeling.ListDataModels(
		1000, nil, nil, false, true,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, dataModel := range dataModelsList.Items {
		fmt.Printf("Space: %s, External ID: %s, Version: %s\n", dataModel.Space, dataModel.ExternalId, dataModel.Version)
	}
	fmt.Println("Data Models Count:", len(dataModelsList.Items))

	// Search for CogniteTimeSeries instances
	fmt.Println("\n### Testing searching for CogniteTimeSeries instances")
	properties := []string{"name", "description"}
	nodeList, err := client.DataModeling.InstancesSearch(
		dto.ViewReference{
			Type:       "view",
			Space:      "cdf_cdm",
			ExternalId: "CogniteTimeSeries",
			Version:    "v1",
		},
		"",
		nil,
		&properties,
		nil,
		nil,
		nil,
		100,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Node List Count:", len(nodeList.Items))
	fmt.Println()

	for _, node := range nodeList.Items {
		fmt.Printf("Space: %s, External ID: %s\n", node.Space, node.ExternalId)

		// Navigate through the nested properties structure
		for _, spaceValue := range node.Properties {
			// Cast the space value to a map
			if spaceMap, ok := spaceValue.(map[string]interface{}); ok {
				// Iterate through the view keys
				for viewKey, viewProperties := range spaceMap {
					fmt.Printf("View: %s\n", viewKey)

					// Cast the view properties to a map
					if propertiesMap, ok := viewProperties.(map[string]interface{}); ok {
						// Now print each property on a new line with the desired format
						for propKey, propValue := range propertiesMap {
							fmt.Printf("- %s: %v\n", propKey, propValue)
						}
					}
				}
			}
		}
		// Fetch the latest data points for the CogniteTimeSeries instance before now
		instanceId := dto.InstanceId{
			Space:      node.Space,
			ExternalId: node.ExternalId,
		}
		var latestDataPointsQueryItems []dto.LatestDataPointsQueryItem
		latestDataPointsQueryItems = append(latestDataPointsQueryItems, dto.LatestDataPointsQueryItem{
			InstanceId: &instanceId,
			Before:     "now",
		})
		latestDataPoints, err := client.TimeSeries.RetrieveLatest(
			&latestDataPointsQueryItems,
			nil,
		)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, latestDataPoint := range latestDataPoints.Items {
			// need to update protobuf to include the instance id
			fmt.Println("- latest data point:", latestDataPoint.DatapointType)
		}

		fmt.Println()
	}

}
