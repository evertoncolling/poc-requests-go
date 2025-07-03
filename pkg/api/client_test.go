package api

import (
	"testing"
)

type mockCredentialProvider struct {
	token string
}

func (m *mockCredentialProvider) FetchToken() string {
	return m.token
}

func TestNewCogniteClient(t *testing.T) {
	tests := []struct {
		name            string
		clientConfig    ClientConfig
		expectedProject string
		expectedCluster string
	}{
		{
			name: "Valid client creation",
			clientConfig: ClientConfig{
				ClientName:  "test-client",
				Cluster:     "test-cluster",
				Project:     "test-project",
				Credentials: &mockCredentialProvider{token: "test-token"},
			},
			expectedProject: "test-project",
			expectedCluster: "test-cluster",
		},
		{
			name: "Empty project",
			clientConfig: ClientConfig{
				ClientName:  "test-client",
				Cluster:     "test-cluster",
				Project:     "",
				Credentials: &mockCredentialProvider{token: "test-token"},
			},
			expectedProject: "",
			expectedCluster: "test-cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewCogniteClient(tt.clientConfig)

			if client.ClientConfig.Project != tt.expectedProject {
				t.Errorf("Expected project %s, got %s", tt.expectedProject, client.ClientConfig.Project)
			}

			if client.ClientConfig.Cluster != tt.expectedCluster {
				t.Errorf("Expected cluster %s, got %s", tt.expectedCluster, client.ClientConfig.Cluster)
			}

			expectedBaseURL := "https://" + tt.expectedCluster + ".cognitedata.com"
			if client.BaseURL != expectedBaseURL {
				t.Errorf("Expected base URL %s, got %s", expectedBaseURL, client.BaseURL)
			}

			if client.TimeSeries.Client == nil {
				t.Error("Expected TimeSeries service to be initialized")
			}

			if client.Units.Client == nil {
				t.Error("Expected Units service to be initialized")
			}

			if client.DataModeling.Client == nil {
				t.Error("Expected DataModeling service to be initialized")
			}

			if client.AccessToken != "test-token" {
				t.Errorf("Expected access token test-token, got %s", client.AccessToken)
			}
		})
	}
}

func TestOAuthClientCredentials_FetchToken(t *testing.T) {
	// This test mainly verifies the struct fields are set correctly
	// In production, you'd want to mock the Azure AD client
	provider := OAuthClientCredentials{
		ClientId:     "test-client",
		ClientSecret: "test-secret",
		AuthorityURI: "https://login.microsoftonline.com/test-tenant",
		Cluster:      "test-cluster",
	}

	if provider.ClientId != "test-client" {
		t.Errorf("Expected client ID test-client, got %s", provider.ClientId)
	}

	if provider.ClientSecret != "test-secret" {
		t.Errorf("Expected client secret test-secret, got %s", provider.ClientSecret)
	}

	if provider.AuthorityURI != "https://login.microsoftonline.com/test-tenant" {
		t.Errorf("Expected authority URI https://login.microsoftonline.com/test-tenant, got %s", provider.AuthorityURI)
	}

	if provider.Cluster != "test-cluster" {
		t.Errorf("Expected cluster test-cluster, got %s", provider.Cluster)
	}
}

func TestAzureADClientCredentials(t *testing.T) {
	result := AzureADClientCredentials("client-id", "client-secret", "tenant-id", "cluster")

	if result.ClientId != "client-id" {
		t.Errorf("Expected client ID client-id, got %s", result.ClientId)
	}

	if result.ClientSecret != "client-secret" {
		t.Errorf("Expected client secret client-secret, got %s", result.ClientSecret)
	}

	expectedAuthority := "https://login.microsoftonline.com/tenant-id"
	if result.AuthorityURI != expectedAuthority {
		t.Errorf("Expected authority URI %s, got %s", expectedAuthority, result.AuthorityURI)
	}

	if result.Cluster != "cluster" {
		t.Errorf("Expected cluster cluster, got %s", result.Cluster)
	}
}

func TestToken_FetchToken(t *testing.T) {
	token := Token{AccessToken: "test-access-token"}
	result := token.FetchToken()

	if result != "test-access-token" {
		t.Errorf("Expected test-access-token, got %s", result)
	}
}
