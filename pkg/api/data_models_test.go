package api

import (
	"testing"
)

func TestDataModeling_struct(t *testing.T) {
	// Test basic struct initialization
	clientConfig := ClientConfig{
		ClientName:  "test-client",
		Cluster:     "test-cluster",
		Project:     "test-project",
		Credentials: &mockCredentialProvider{token: "test-token"},
	}

	client := NewCogniteClient(clientConfig)

	if client.DataModeling.Client == nil {
		t.Error("Expected DataModeling.Client to be non-nil")
	}

	if client.DataModeling.Client.ClientConfig.Project != client.ClientConfig.Project {
		t.Error("Expected DataModeling.Client to be properly initialized")
	}
}
