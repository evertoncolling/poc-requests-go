package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"poc-requests-go/pkg/api"
	"poc-requests-go/pkg/dto"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	grob "github.com/MetalBlueberry/go-plotly/graph_objects"
	"github.com/MetalBlueberry/go-plotly/offline"
	"github.com/joho/godotenv"
)

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

	baseURL := fmt.Sprintf("https://%s.cognitedata.com", cluster)
	scopes := []string{
		fmt.Sprintf("https://%s.cognitedata.com/.default", cluster),
		"offline_access",
		"openid",
		"profile",
	}

	// confidential clients have a credential, such as a secret or a certificate
	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		log.Fatalf("Error creating cred from secret: %v", err)
	}

	authorityURI := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
	confidentialClient, err := confidential.New(authorityURI, clientID, cred)
	if err != nil {
		log.Fatalf("Error creating confidential client: %v", err)
	}

	// Get a token for the app itself
	result, err := confidentialClient.AcquireTokenSilent(context.TODO(), scopes)
	if err != nil {
		// cache miss, authenticate with another AcquireToken... method
		result, err = confidentialClient.AcquireTokenByCredential(context.TODO(), scopes)
		if err != nil {
			log.Fatalf("Error acquiring token: %v", err)
		}
	}
	accessToken := result.AccessToken

	fmt.Println("### Testing fetching some time series")

	// List time series
	tsList, err := api.ListTimeSeries(project, accessToken, baseURL, 100, false, "", "", nil, nil, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Time Series Count:", len(tsList.Items))

	// Filter time series
	filter := &dto.TimeSeriesFilter{
		UnitQuantity: "Pressure",
	}
	filteredTsList, err := api.FilterTimeSeries(
		project, accessToken, baseURL, filter, nil, 100, "", "", nil,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Filtered Time Series Count:", len(filteredTsList.Items))

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
	dataPoints, err := api.RetrieveData(
		project,
		accessToken,
		baseURL,
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
	dataPoints, err = api.RetrieveData(
		project,
		accessToken,
		baseURL,
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
					Text: dps.Unit,
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
