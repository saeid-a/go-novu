package lib_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/saeid-a/go-novu/lib"
	"github.com/stretchr/testify/require"
)

func TestCreateWorkflow_Success(t *testing.T) {
	workflowCreate := lib.CreateWorkflowRequest{}
	httpServer := createTestServer(t, TestServerOptions[lib.CreateWorkflowRequest, map[string]string]{
		expectedURLPath:    "/v1/workflows",
		expectedSentMethod: http.MethodPost,
		expectedSentBody:   lib.CreateWorkflowRequest{},
		responseStatusCode: http.StatusCreated,
		responseBody:       map[string]string{},
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	err := c.WorkflowApi.Create(ctx, workflowCreate)

	require.NoError(t, err)
}

func TestListWorkflows_Success(t *testing.T) {
	body := map[string]string{}
	var expectedResponse *lib.ListWorkflowsResponse = &lib.ListWorkflowsResponse{
		Page:       0,
		PageSize:   20,
		TotalCount: 1,
		Data: []lib.GetWorkflowResponse{{
			ID:             "id",
			OrganizationID: "orgId",
			EnvironmentID:  "envId",
			Name:           "topicName",
		}},
	}

	httpServer := createTestServer(t, TestServerOptions[map[string]string, *lib.ListWorkflowsResponse]{
		expectedURLPath:    "/v1/workflows",
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   body,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.WorkflowApi.List(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}
func TestGetWorkflow_Success(t *testing.T) {
	body := map[string]string{}
	var expectedResponse *lib.ListWorkflowsResponse = &lib.ListWorkflowsResponse{
		Page:       0,
		PageSize:   20,
		TotalCount: 1,
		Data: []lib.GetWorkflowResponse{{
			ID:             "id",
			OrganizationID: "orgId",
			EnvironmentID:  "envId",
			Name:           "topicName",
		}},
	}

	httpServer := createTestServer(t, TestServerOptions[map[string]string, *lib.ListWorkflowsResponse]{
		expectedURLPath:    "/v1/workflows",
		expectedSentMethod: http.MethodGet,
		expectedSentBody:   body,
		responseStatusCode: http.StatusOK,
		responseBody:       expectedResponse,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	resp, err := c.WorkflowApi.List(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)
}

func TestDeleteWorkflow_Success(t *testing.T) {
	key := "topicKey"
	body := map[string]string{}

	httpServer := createTestServer(t, TestServerOptions[map[string]string, map[string]string]{
		expectedURLPath:    fmt.Sprintf("/v1/workflows/%s", key),
		expectedSentMethod: http.MethodDelete,
		expectedSentBody:   body,
		responseStatusCode: http.StatusNoContent,
		responseBody:       body,
	})

	ctx := context.Background()
	c := lib.NewAPIClient(novuApiKey, &lib.Config{BackendURL: lib.MustParseURL(httpServer.URL)})
	err := c.WorkflowApi.Delete(ctx, key)

	require.NoError(t, err)
}
