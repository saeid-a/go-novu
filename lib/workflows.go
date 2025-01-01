package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type IWorkflows interface {
	List(ctx context.Context, options *ListTopicsOptions) (*ListTopicsResponse, error)
	Create(ctx context.Context, workflow CreateWorkflowRequest) error
	Get(ctx context.Context, key string) (*GetWorkflowResponse, error)
	Update(ctx context.Context, key string, workflow UpdateWorkflowRequest) error
	Delete(ctx context.Context, key string) error
	UpdateStatus(ctx context.Context, key string, subscribers []string) error
}

type WorkflowService service

func (t *WorkflowService) Create(ctx context.Context, workflow CreateWorkflowRequest) error {
	var resp interface{}
	URL := t.client.config.BackendURL.JoinPath("topics")

	reqBody := workflow

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	httpResponse, err := t.client.sendRequest(req, &resp)
	if err != nil {
		return err
	}

	if httpResponse.StatusCode != HTTPStatusCreated {
		return errors.Wrap(err, "unable to create topic")
	}

	return nil
}

func (t *WorkflowService) List(ctx context.Context, options *ListTopicsOptions) (*ListTopicsResponse, error) {
	var resp ListTopicsResponse
	URL := t.client.config.BackendURL.JoinPath("topics")

	if options == nil {
		options = &ListTopicsOptions{}
	}
	queryParams, _ := json.Marshal(options)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), bytes.NewBuffer(queryParams))
	if err != nil {
		return nil, err
	}

	_, err = t.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (t *WorkflowService) Get(ctx context.Context, key string) (*GetWorkflowResponse, error) {
	var resp GetWorkflowResponse
	URL := t.client.config.BackendURL.JoinPath("workflows", key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	_, err = t.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (t *WorkflowService) Update(ctx context.Context, key string, workflow UpdateWorkflowRequest) error {
	var resp GetWorkflowResponse
	URL := t.client.config.BackendURL.JoinPath("workflows", key)

	queryParams, _ := json.Marshal(workflow)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL.String(), bytes.NewBuffer(queryParams))
	if err != nil {
		return err
	}

	_, err = t.client.sendRequest(req, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowService) UpdateStatus(ctx context.Context, key string, subscribers []string) error {
	URL := t.client.config.BackendURL.JoinPath("topics", key, "subscribers/removal")

	queryParams, _ := json.Marshal(SubscribersTopicRequest{
		Subscribers: subscribers,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(queryParams))
	if err != nil {
		return err
	}

	_, err = t.client.sendRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowService) Delete(ctx context.Context, key string) error {
	var resp interface{}
	URL := t.client.config.BackendURL.JoinPath("workflows", key)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, URL.String(), http.NoBody)
	if err != nil {
		return err
	}

	_, err = t.client.sendRequest(req, &resp)
	if err != nil {
		return err
	}

	return nil
}
