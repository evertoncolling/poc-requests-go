package api

import (
	"testing"
)

func TestTimeSeries_struct(t *testing.T) {
	// Test basic struct initialization
	clientConfig := ClientConfig{
		ClientName:  "test-client",
		Cluster:     "test-cluster",
		Project:     "test-project",
		Credentials: &mockCredentialProvider{token: "test-token"},
	}

	client := NewCogniteClient(clientConfig)

	if client.TimeSeries.Client == nil {
		t.Error("Expected TimeSeries.Client to be non-nil")
	}

	if client.TimeSeries.Client.ClientConfig.Project != client.ClientConfig.Project {
		t.Error("Expected TimeSeries.Client to be properly initialized")
	}
}
