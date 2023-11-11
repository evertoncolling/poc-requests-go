package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/evertoncolling/poc-requests-go/pkg/api"
	"github.com/evertoncolling/poc-requests-go/pkg/dto"

	grob "github.com/MetalBlueberry/go-plotly/graph_objects"
	"github.com/MetalBlueberry/go-plotly/offline"
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
		nil, nil, nil, nil, nil, nil, nil,
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

	// let's try a more reasonable number of data points
	fmt.Println("\n### Data Points (for plotting)")
	items = []dto.DataPointsQueryItem{
		{
			ExternalId:  "EVE-TI-FORNEBU-01-2",
			Start:       "300d-ago",
			End:         "now",
			Aggregates:  []string{"average"},
			Granularity: "1h",
			Limit:       10000,
		},
	}
	start = time.Now()
	dataPoints, err = client.TimeSeries.RetrieveData(
		&items,
		nil, nil, nil, nil, nil, nil, nil,
	)
	elapsed = time.Since(start)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Prepare the data for plotting
	dps = dataPoints.Items[0]
	aggregateDatapointsType := dps.DatapointType.(*dto.DataPointListItem_AggregateDatapoints)
	aggregateDatapoints := aggregateDatapointsType.AggregateDatapoints.Datapoints
	values := make([]float64, 0, len(aggregateDatapoints))
	timestamps := make([]time.Time, 0, len(aggregateDatapoints))
	for _, datapoint := range aggregateDatapoints {
		if datapoint != nil {
			values = append(values, datapoint.Average)
			timestamp := time.Unix(0, datapoint.Timestamp*int64(time.Millisecond))
			timestamps = append(timestamps, timestamp)
		}
	}
	fmt.Println("Data Points Count:", len(values))
	fmt.Printf("Time taken: %s\n", elapsed)
	unit, err := findUnitByExternalId(&unitList, dps.UnitExternalId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Plot the data
	fig := &grob.Fig{
		Data: grob.Traces{
			&grob.Bar{
				Type: grob.TraceTypeScatter,
				X:    timestamps,
				Y:    values,
				Name: dps.ExternalId,
			},
		},
		Layout: &grob.Layout{
			Title: &grob.LayoutTitle{
				Text: fmt.Sprintf("Fetched %d data points - hourly aggregates", len(values)),
				Font: &grob.LayoutTitleFont{
					Size: 24.0,
				},
			},
			Font: &grob.LayoutFont{
				Family: "Roboto",
				Size:   12,
				Color:  "black",
			},
			Xaxis: &grob.LayoutXaxis{
				Showspikes:     grob.True,
				Spikemode:      grob.LayoutXaxisSpikemode("across"),
				Spikethickness: 1.0,
				Spikedash:      "solid",
			},
			Yaxis: &grob.LayoutYaxis{
				Title: &grob.LayoutYaxisTitle{
					Text: fmt.Sprintf("%s [%s]", unit.Quantity, unit.Symbol),
				},
			},
			Spikedistance: -1,
			Legend: &grob.LayoutLegend{
				Orientation: grob.LayoutLegendOrientation("h"),
				Yanchor:     grob.LayoutLegendYanchor("bottom"),
				Y:           1.02,
				Xanchor:     grob.LayoutLegendXanchor("right"),
				X:           1.0,
			},
			Showlegend: grob.True,
			Height:     800,
		},
	}

	offline.Show(fig)
}
