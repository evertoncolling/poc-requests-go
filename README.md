# poc-requests-go

A Go client library for interacting with Cognite Data Fusion (CDF) APIs. This library provides a simple interface for accessing CDF services including time series, units, and data modeling capabilities.

[![CI](https://github.com/evertoncolling/poc-requests-go/actions/workflows/ci.yml/badge.svg)](https://github.com/evertoncolling/poc-requests-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/evertoncolling/poc-requests-go)](https://goreportcard.com/report/github.com/evertoncolling/poc-requests-go)
[![codecov](https://codecov.io/gh/evertoncolling/poc-requests-go/branch/main/graph/badge.svg)](https://codecov.io/gh/evertoncolling/poc-requests-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/evertoncolling/poc-requests-go.svg)](https://pkg.go.dev/github.com/evertoncolling/poc-requests-go)

## Installation

```bash
go get github.com/evertoncolling/poc-requests-go
```

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/evertoncolling/poc-requests-go/pkg/api"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Create credential provider
    credProvider := &api.AzureCredentialProvider{
        TenantID:     os.Getenv("TENANT_ID"),
        ClientID:     os.Getenv("CLIENT_ID"),
        ClientSecret: os.Getenv("CLIENT_SECRET"),
        Scopes:       []string{"https://bluefield.cognitedata.com/.default"},
    }

    // Create client
    client := api.NewCogniteClient(os.Getenv("CDF_PROJECT"), credProvider)

    // List time series
    ctx := context.Background()
    timeSeries, err := client.TimeSeries.List(ctx, 10, false)
    if err != nil {
        log.Fatal("Error fetching time series:", err)
    }

    fmt.Printf("Found %d time series\n", len(timeSeries.Items))
}
```

## Supported Endpoints

### Time Series API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `List()` | `GET /timeseries` | List time series with optional filtering |
| `Filter()` | `POST /timeseries/list` | Advanced filtering of time series |
| `RetrieveData()` | `POST /timeseries/data/list` | Retrieve time series data points |
| `RetrieveLatest()` | `POST /timeseries/data/latest` | Get latest data points for time series |

#### Examples

```go
// List time series
timeSeries, err := client.TimeSeries.List(ctx, 100, true)

// Filter time series
filter := dto.TimeSeriesFilter{
    Filter: dto.FilterDetails{
        Name: &dto.StringFilter{Value: "temperature"},
    },
    Limit: 10,
}
result, err := client.TimeSeries.Filter(ctx, filter)

// Retrieve data points
dataRequest := dto.DataRequest{
    Items: []dto.DataRequestItem{
        {ExternalID: "temperature_sensor_1"},
    },
    Start: "2d-ago",
    End:   "now",
}
data, err := client.TimeSeries.RetrieveData(ctx, dataRequest)
```

### Units API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `List()` | `GET /units` | Retrieve the complete units catalog |

#### Examples

```go
// List all units
units, err := client.Units.List(ctx)
fmt.Printf("Found %d units\n", len(units.Items))

// Find specific unit
for _, unit := range units.Items {
    if unit.ExternalID == "temperature:deg_c" {
        fmt.Printf("Unit: %s, Symbol: %s\n", unit.Name, unit.Symbol)
    }
}
```

### Data Modeling API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `ListDataModels()` | `GET /models/datamodels` | List available data models |
| `InstancesSearch()` | `POST /models/instances/search` | Search for instances in data models |
| `GraphQLQuery()` | `POST /models/graphql` | Execute GraphQL queries |

#### Examples

```go
// List data models
models, err := client.DataModeling.ListDataModels(ctx)

// Search instances
searchRequest := dto.InstancesSearchRequest{
    View: dto.ViewIdentifier{
        Type:       "view",
        Space:      "my-space",
        ExternalID: "my-view",
    },
    InstanceType: "node",
    Limit:        10,
}
instances, err := client.DataModeling.InstancesSearch(ctx, searchRequest)

// Execute GraphQL query
query := dto.GraphQLQueryRequest{
    Query: `
        query {
            listPumps {
                items {
                    externalId
                    name
                }
            }
        }
    `,
}
result, err := client.DataModeling.GraphQLQuery(ctx, query)
```

## Development

### Prerequisites

- Go 1.22 or higher
- Make (optional, for using Makefile commands)

### Building and Testing

```bash
# Install development tools
make install-tools

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Build the project
make build

# Run all quality checks
make check
```

### Available Make Commands

```bash
make help                  # Show available commands
make build                 # Build the application
make test                  # Run unit tests (excludes integration tests)
make test-unit             # Run unit tests only
make test-integration      # Run integration tests (loads .env locally)
make test-integration-ci   # Run integration tests (uses env vars directly)
make test-all              # Run all tests including integration tests
make test-coverage         # Run tests with coverage report
make lint                  # Run linter
make fmt                   # Format code
make vet                   # Run go vet
make check                 # Run all quality checks
make clean                 # Clean build artifacts
make deps                  # Download dependencies
```

### Testing

The project includes both unit tests and integration tests:

- **Unit tests**: Mock-based tests that run quickly without external dependencies
- **Integration tests**: Real API tests that make actual calls to CDF endpoints

#### Running Tests

```bash
# Run only unit tests (fast, no credentials needed)
make test-unit

# Run integration tests locally (loads .env file automatically)
make test-integration

# Run all tests (unit + integration)
make test-all
```

#### Setting up Integration Tests

##### Local Development

Create a `.env` file in the project root with your CDF credentials:

```bash
# .env file
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
TENANT_ID=your_tenant_id
CDF_CLUSTER=your_cluster_name
CDF_PROJECT=your_project_name
```

The `make test-integration` target will automatically load these variables.

Alternatively, you can use the provided script:

```bash
# Run integration tests with the script
./scripts/test-integration.sh
```

##### CI/CD Setup

For GitHub Actions, set the following repository secrets:

1. Go to your GitHub repository
2. Navigate to Settings > Secrets and variables > Actions
3. Add the following secrets:
   - `CLIENT_ID`
   - `CLIENT_SECRET`
   - `TENANT_ID`
   - `CDF_CLUSTER`
   - `CDF_PROJECT`

Integration tests will run automatically in CI for pushes and pull requests.

## Project Structure

```
.
├── .github/workflows/     # GitHub Actions CI/CD
├── pkg/
│   ├── api/              # API client implementations
│   ├── dto/              # Data transfer objects
│   └── proto/            # Protocol buffer definitions
├── main.go               # Example application
├── Makefile              # Build automation
└── .golangci.yml         # Linter configuration
```

## Running the Example

Create a `.env` file in your project root with the following variables:

```bash
CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TENANT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CDF_CLUSTER=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CDF_PROJECT=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

```bash
# Run the example
go run main.go
```
