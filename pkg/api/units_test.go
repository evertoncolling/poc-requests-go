package api

import (
	"testing"
)

func TestUnits_struct(t *testing.T) {
	// Test basic struct initialization
	clientConfig := ClientConfig{
		ClientName:  "test-client",
		Cluster:     "test-cluster",
		Project:     "test-project",
		Credentials: &mockCredentialProvider{token: "test-token"},
	}

	client := NewCogniteClient(clientConfig)

	if client.Units.Client == nil {
		t.Error("Expected Units.Client to be non-nil")
	}

	if client.Units.Client.ClientConfig.Project != client.ClientConfig.Project {
		t.Error("Expected Units.Client to be properly initialized")
	}
}
