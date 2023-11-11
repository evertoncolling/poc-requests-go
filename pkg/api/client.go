package api

import (
	"context"
	"fmt"
	"log"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

const (
	VERSION = "0.0.1"
)

type OAuthClientCredentials struct {
	ClientId     string
	ClientSecret string
	AuthorityURI string
	Cluster      string
}

func AzureADClientCredentials(
	clientId string,
	clientSecret string,
	tenantId string,
	cluster string,
) OAuthClientCredentials {
	return OAuthClientCredentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AuthorityURI: fmt.Sprintf("https://login.microsoftonline.com/%s", tenantId),
		Cluster:      cluster,
	}
}

func (m *OAuthClientCredentials) FetchToken() string {
	scopes := []string{
		fmt.Sprintf("https://%s.cognitedata.com/.default", m.Cluster),
	}

	cred, err := confidential.NewCredFromSecret(m.ClientSecret)
	if err != nil {
		log.Fatalf("Error creating cred from secret: %v", err)
	}

	confidentialClient, err := confidential.New(m.AuthorityURI, m.ClientId, cred)
	if err != nil {
		log.Fatalf("Error creating confidential client: %v", err)
	}

	result, err := confidentialClient.AcquireTokenSilent(context.TODO(), scopes)
	if err != nil {
		// cache miss, authenticate with another AcquireToken... method
		result, err = confidentialClient.AcquireTokenByCredential(context.TODO(), scopes)
		if err != nil {
			log.Fatalf("Error acquiring token: %v", err)
		}
	}
	return result.AccessToken
}

type ClientConfig struct {
	ClientName  string
	Cluster     string
	Project     string
	Credentials OAuthClientCredentials
}

type CogniteClient struct {
	ClientConfig ClientConfig
	AccessToken  string
	BaseURL      string
	Headers      map[string]string
	TimeSeries   TimeSeries
	Units        Units
}

type TimeSeries struct {
	Client *CogniteClient // Add a reference to the CogniteClient
}

type Units struct {
	Client *CogniteClient // Add a reference to the CogniteClient
}

func NewCogniteClient(
	clientConfig ClientConfig,
) CogniteClient {
	baseURL := fmt.Sprintf("https://%s.cognitedata.com", clientConfig.Cluster)
	accessToken := clientConfig.Credentials.FetchToken()
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"x-cdp-app":     clientConfig.ClientName,
		"x-cdp-sdk":     fmt.Sprintf("poc-requests-go:%s", VERSION),
		"cdf-version":   "beta",
	}
	client := CogniteClient{
		ClientConfig: clientConfig,
		AccessToken:  accessToken,
		BaseURL:      baseURL,
		Headers:      headers,
	}
	client.TimeSeries = TimeSeries{Client: &client}
	client.Units = Units{Client: &client}
	return client
}
