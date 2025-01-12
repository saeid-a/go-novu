package lib

import (
	"context"
	"fmt"
)

type IEnvironment interface {
	Current(ctx context.Context) (EnvironmentResponse, error)
	GetAll(ctx context.Context) (EnvironmentsResponse, error)
	Update(ctx context.Context, data BroadcastEventToAll) (EventResponse, error)
}

type EnvironmentService service

func (e EnvironmentService) Trigger(ctx context.Context, eventId string, data ITriggerPayloadOptions) (EventResponse, error) {
	//TODO implement me
	fmt.Println("implement me")
	return EventResponse{}, nil
}

func (e EnvironmentService) TriggerBulk(ctx context.Context, data []BulkTriggerOptions) ([]EventResponse, error) {
	//TODO implement me
	fmt.Println("implement me")
	return nil, nil
}

func (e EnvironmentService) BroadcastToAll(ctx context.Context, data BroadcastEventToAll) (EventResponse, error) {
	//TODO implement me
	fmt.Println("implement me")
	return EventResponse{}, nil
}

func (e EnvironmentService) CancelTrigger(ctx context.Context, transactionId string) (bool, error) {
	//TODO implement me
	fmt.Println("implement me")
	return true, nil

}

//func (e *EnvironmentService) Current(ctx context.Context) (EnvironmentResponse, error) {
//	var resp EventResponse
//	URL := e.client.config.BackendURL.JoinPath("events/trigger")
//
//	reqBody := EventRequest{
//		Name:          eventId,
//		To:            data.To,
//		Payload:       data.Payload,
//		Overrides:     data.Overrides,
//		TransactionId: data.TransactionId,
//		Actor:         data.Actor,
//	}
//
//	jsonBody, _ := json.Marshal(reqBody)
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
//	if err != nil {
//		return resp, err
//	}
//
//	_, err = e.client.sendRequest(req, &resp)
//	if err != nil {
//		return resp, err
//	}
//
//	return resp, nil
//}

var _ IEvent = &EnvironmentService{}
