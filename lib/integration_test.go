package lib_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/saeid-a/go-novu/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type IntegrationRequestDetails[T any] struct {
	Url    string
	Method string
	Body   T
}

type IntegrationResponseDetails struct {
	StatusCode int
	Body       interface{}
}

type IntegrationServerOptions[T any] struct {
	ExpectedRequest IntegrationRequestDetails[T]

	ExpectedResponse IntegrationResponseDetails
}

func ValidateIntegrationRequest[T any](t *testing.T, req *http.Request, expectedRequest IntegrationRequestDetails[T]) {
	t.Run("Request must be authorized", func(t *testing.T) {
		authKey := req.Header.Get("Authorization")
		assert.True(t, strings.Contains(authKey, novuApiKey))
		assert.True(t, strings.HasPrefix(authKey, "ApiKey"))
	})

	t.Run("URL and request method is as expected", func(t *testing.T) {
		assert.Equal(t, expectedRequest.Method, req.Method)
		assert.True(t, strings.HasPrefix(expectedRequest.Url, "/v1/integrations"))
		assert.Equal(t, expectedRequest.Url, req.RequestURI)
	})
}

func IntegrationTestServer[T any](t *testing.T, options IntegrationServerOptions[T]) *httptest.Server {
	var receivedBody T

	integrationService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ValidateIntegrationRequest(t, req, options.ExpectedRequest)

		if req.Body == http.NoBody {
			assert.Empty(t, options.ExpectedRequest.Body)
		} else {
			if err := json.NewDecoder(req.Body).Decode(&receivedBody); err != nil {
				log.Printf("error in unmarshalling %+v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			t.Run("Request is as expected", func(t *testing.T) {
				assert.Equal(t, options.ExpectedRequest.Body, receivedBody)
			})
		}

		w.WriteHeader(options.ExpectedResponse.StatusCode)

		bb, _ := json.Marshal(options.ExpectedResponse.Body)
		w.Write(bb)

	}))

	t.Cleanup(func() {
		integrationService.Close()
	})

	return integrationService
}

func TestCreateIntegration_Success(t *testing.T) {
	createIntegrationRequest := lib.CreateIntegrationRequest{
		ProviderID: "sendgrid",
		Channel:    "email",
		Credentials: lib.IntegrationCredentials{
			ApiKey:    "api_key",
			SecretKey: "secret_key",
		},
		Active: true,
		Check:  false,
	}

	var response *lib.IntegrationResponse
	fileToStruct(filepath.Join("../testdata", "integration_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[lib.CreateIntegrationRequest]{
		ExpectedRequest: IntegrationRequestDetails[lib.CreateIntegrationRequest]{
			Url:    "/v1/integrations",
			Method: http.MethodPost,
			Body:   createIntegrationRequest,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusCreated,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	res, err := novuClient.IntegrationsApi.Create(ctx, createIntegrationRequest)

	assert.Equal(t, response, res)

	require.NoError(t, err)
}

func TestGetAllIntegration_Success(t *testing.T) {

	var response *lib.GetIntegrationsResponse
	fileToStruct(filepath.Join("../testdata", "get_integrations_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    "/v1/integrations",
			Method: http.MethodGet,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.GetAll(ctx)

	assert.Equal(t, response, res)

	require.NoError(t, err)
}

func TestGetActiveIntegration_Success(t *testing.T) {

	var response *lib.GetIntegrationsResponse
	fileToStruct(filepath.Join("../testdata", "get_active_integrations_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    "/v1/integrations/active",
			Method: http.MethodGet,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.GetActive(ctx)

	assert.Equal(t, response, res)

	require.NoError(t, err)
}

func TestGetWebhookSupportStatusIntegration_Success(t *testing.T) {
	providerId := "sendgrid"
	response := true
	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    fmt.Sprintf("/v1/integrations/webhook/provider/%s/status", providerId),
			Method: http.MethodGet,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.GetWebhookSupportStatus(ctx, providerId)
	assert.Equal(t, response, res)
	require.NoError(t, err)
}

func TestUpdateIntegration_Success(t *testing.T) {
	const integrationId = "integrationId"

	updateIntegrationRequest := lib.UpdateIntegrationRequest{
		Credentials: lib.IntegrationCredentials{
			ApiKey:    "new_api_key",
			SecretKey: "new_secret_key",
		},
		Active: true,
		Check:  false,
	}

	var response *lib.IntegrationResponse
	fileToStruct(filepath.Join("../testdata", "integration_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[lib.UpdateIntegrationRequest]{
		ExpectedRequest: IntegrationRequestDetails[lib.UpdateIntegrationRequest]{
			Url:    fmt.Sprintf("/v1/integrations/%s", integrationId),
			Method: http.MethodPut,
			Body:   updateIntegrationRequest,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.Update(ctx, integrationId, updateIntegrationRequest)

	assert.Equal(t, response, res)

	require.NoError(t, err)
}

func TestDeleteActiveIntegration_Success(t *testing.T) {
	const integrationId = "integrationId"

	var response *lib.IntegrationResponse
	fileToStruct(filepath.Join("../testdata", "delete_integration_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    fmt.Sprintf("/v1/integrations/%s", integrationId),
			Method: http.MethodDelete,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.Delete(ctx, integrationId)

	assert.Equal(t, response, res)

	require.NoError(t, err)
}

func TestSetIntegrationAsPrimary_Success(t *testing.T) {
	const integrationId = "IntegrationId"

	var response *lib.SetIntegrationAsPrimaryResponse
	fileToStruct(filepath.Join("../testdata", "set_integration_as_primary_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    fmt.Sprintf("/v1/integrations/%s/set-primary", integrationId),
			Method: http.MethodPost,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.SetIntegrationAsPrimary(ctx, integrationId)

	assert.Equal(t, response, res)
	require.NoError(t, err)
}

func TestGetChannelLimit_Success(t *testing.T) {
	const channelType = "ChannelType"

	var response *lib.IntegrationChannelLimitResponse
	fileToStruct(filepath.Join("../testdata", "integration_channel_limit_response.json"), &response)

	httpServer := IntegrationTestServer(t, IntegrationServerOptions[interface{}]{
		ExpectedRequest: IntegrationRequestDetails[interface{}]{
			Url:    fmt.Sprintf("/v1/integrations/%s/limit", channelType),
			Method: http.MethodGet,
		},
		ExpectedResponse: IntegrationResponseDetails{
			StatusCode: http.StatusOK,
			Body:       response,
		},
	})

	ctx := context.Background()
	novuClient := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})

	res, err := novuClient.IntegrationsApi.GetChannelLimit(ctx, channelType)

	assert.Equal(t, response, res)
	require.NoError(t, err)
}
