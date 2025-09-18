package psc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
	WaitForResource(ctx context.Context, id string, resource interface{}) error
}

type Log interface {
	Printf(format string, args ...interface{})
}

type API struct {
	client     HttpClient
	taskWaiter TaskWaiter
	logger     Log
}

func NewAPI(client HttpClient, taskWaiter TaskWaiter, logger Log) *API {
	return &API{client: client, taskWaiter: taskWaiter, logger: logger}
}

func (a *API) GetService(ctx context.Context, subscription int) (*PrivateServiceConnectService, error) {
	message := fmt.Sprintf("get private service connect for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect", subscription)
	task, err := a.getService(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

func (a *API) GetActiveActiveService(ctx context.Context, subscription int, regionId int) (*PrivateServiceConnectService, error) {
	message := fmt.Sprintf("get private service connect for subscription %d in region %d", subscription, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect", subscription, regionId)
	task, err := a.getService(ctx, message, path)
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return task, nil
}

func (a *API) CreateService(ctx context.Context, subscription int) (int, error) {
	message := fmt.Sprintf("create private service connect for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect", subscription)
	resourceId, err := a.create(ctx, message, path)
	if err != nil {
		return 0, wrap404Error(subscription, err)
	}
	return resourceId, nil
}

func (a *API) CreateActiveActiveService(ctx context.Context, subscription int, regionId int) (int, error) {
	message := fmt.Sprintf("create private service connect for subscription %d in region %d", subscription, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect", subscription, regionId)
	resourceId, err := a.create(ctx, message, path)
	if err != nil {
		return 0, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return resourceId, nil
}

func (a *API) GetEndpointCreationScripts(ctx context.Context, subscription int, pscServiceId int, endpointId int, includeTerraformGcpScript bool) (*CreationScript, error) {
	message := fmt.Sprintf("get private service connect creation script for subscription %d, service %d and endpoint %d",
		subscription, pscServiceId, endpointId)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d/endpoints/%d/creationScripts",
		subscription, pscServiceId, endpointId)
	creationScript, err := a.getCreationScript(ctx, message, path, url.Values{
		"includeTerraformGcpScript": []string{strconv.FormatBool(includeTerraformGcpScript)},
	})
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return creationScript, nil
}

func (a *API) GetActiveActiveEndpointCreationScripts(ctx context.Context, subscription int, regionId int, pscServiceId int, endpointId int, includeTerraformGcpScript bool) (*CreationScript, error) {
	message := fmt.Sprintf("get private service connect creation script for subscription %d, service %d and endpoint %d in region %d",
		subscription, pscServiceId, endpointId, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d/endpoints/%d/creationScripts",
		subscription, regionId, pscServiceId, endpointId)
	creationScript, err := a.getCreationScript(ctx, message, path, url.Values{
		"includeTerraformGcpScript": []string{strconv.FormatBool(includeTerraformGcpScript)},
	})
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return creationScript, nil
}

func (a *API) GetEndpointDeletionScripts(ctx context.Context, subscription int, pscServiceId int, endpointId int) (*DeletionScript, error) {
	message := fmt.Sprintf("get private service connect deletion script for subscription %d, service %d and endpoint %d",
		subscription, pscServiceId, endpointId)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d/endpoints/%d/deletionScripts",
		subscription, pscServiceId, endpointId)
	deletionScript, err := a.getDeletionScript(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return deletionScript, nil
}

func (a *API) GetActiveActiveEndpointDeletionScripts(ctx context.Context, subscription int, regionId int, pscServiceId int, endpointId int) (*DeletionScript, error) {
	message := fmt.Sprintf("get private service connect deletion script for subscription %d, service %d and endpoint %d in region %d",
		subscription, pscServiceId, endpointId, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d/endpoints/%d/deletionScripts",
		subscription, regionId, pscServiceId, endpointId)
	deletionScript, err := a.getDeletionScript(ctx, message, path)
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return deletionScript, nil
}

func (a *API) GetEndpoints(ctx context.Context, subscription int, pscServiceId int) (*PrivateServiceConnectEndpoints, error) {
	message := fmt.Sprintf("get private service connect for subscription %d and service %d", subscription, pscServiceId)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d", subscription, pscServiceId)
	endpoints, err := a.getEndpoints(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return endpoints, nil
}

func (a *API) GetActiveActiveEndpoints(ctx context.Context, subscription int, regionId int, pscServiceId int) (*PrivateServiceConnectEndpoints, error) {
	message := fmt.Sprintf("get private service connect for subscription %d and service %d in region %d", subscription, pscServiceId, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d", subscription, regionId, pscServiceId)
	endpoints, err := a.getEndpoints(ctx, message, path)
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return endpoints, nil
}

func (a *API) CreateEndpoint(ctx context.Context, subscription int, pscServiceId int, endpoint CreatePrivateServiceConnectEndpoint) (int, error) {
	message := fmt.Sprintf("create private service connect endpoint for subscription %d and service %d", subscription, pscServiceId)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d", subscription, pscServiceId)

	var task internal.TaskResponse
	err := a.client.Post(ctx, message, path, endpoint, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for private service connect endpoint for subscription %d and service %d to finish being created", subscription, pscServiceId)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *API) CreateActiveActiveEndpoint(ctx context.Context, subscription int, regionId int, pscServiceId int, endpoint CreatePrivateServiceConnectEndpoint) (int, error) {
	message := fmt.Sprintf("create private service connect endpoint for subscription %d and service %d in region %d", subscription, pscServiceId, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d", subscription, regionId, pscServiceId)

	var task internal.TaskResponse
	err := a.client.Post(ctx, message, path, endpoint, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for private service connect endpoint for subscription %d and service %d in region %d to finish being created", subscription, pscServiceId, regionId)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *API) UpdateEndpoint(ctx context.Context, subscription int, pscServiceId int, endpointId int,
	endpoint *UpdatePrivateServiceConnectEndpoint) error {
	message := fmt.Sprintf("update private service connect endpoint %d/%d for subscription %d", pscServiceId, endpointId, subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d/endpoints/%d", subscription, pscServiceId, endpointId)
	err := a.update(ctx, message, path, endpoint)
	if err != nil {
		return wrap404Error(subscription, err)
	}
	return nil
}

func (a *API) UpdateActiveActiveEndpoint(ctx context.Context, subscription int, regionId int, pscServiceId int,
	endpointId int, endpoint *UpdatePrivateServiceConnectEndpoint) error {
	message := fmt.Sprintf("update private service connect endpoint  %d/%d for subscription %d in region %d", pscServiceId, endpointId, subscription, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d/endpoints/%d", subscription, regionId, pscServiceId, endpointId)
	err := a.update(ctx, message, path, endpoint)
	if err != nil {
		return wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return nil
}

func (a *API) DeleteService(ctx context.Context, subscription int) error {
	message := fmt.Sprintf("delete private service connect for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect", subscription)
	err := a.delete(ctx, message, path)
	if err != nil {
		return wrap404Error(subscription, err)
	}
	return nil
}

func (a *API) DeleteActiveActiveService(ctx context.Context, subscription int, regionId int) error {
	message := fmt.Sprintf("delete private service connect for subscription %d in region %d", subscription, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect", subscription, regionId)
	err := a.delete(ctx, message, path)
	if err != nil {
		return wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return nil
}

func (a *API) DeleteEndpoint(ctx context.Context, subscription int, pscServiceId int,
	endpointId int) error {
	message := fmt.Sprintf("delete private service connect endpoint %d/%d for subscription %d", pscServiceId, endpointId, subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-service-connect/%d/endpoints/%d", subscription, pscServiceId, endpointId)
	err := a.delete(ctx, message, path)
	if err != nil {
		return wrap404Error(subscription, err)
	}
	return nil
}

func (a *API) DeleteActiveActiveEndpoint(ctx context.Context, subscription int, regionId int, pscServiceId int,
	endpointId int) error {
	message := fmt.Sprintf("delete private service connect endpoint %d/%d for subscription %d in region %d", pscServiceId, endpointId, subscription, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-service-connect/%d/endpoints/%d", subscription, regionId, pscServiceId, endpointId)
	err := a.delete(ctx, message, path)
	if err != nil {
		return wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return nil
}

func (a *API) getService(ctx context.Context, message string, path string) (*PrivateServiceConnectService, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for private service connect request %d to complete", task.ID)

	var response PrivateServiceConnectService
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) getEndpoints(ctx context.Context, message string, path string) (*PrivateServiceConnectEndpoints, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for private service connect request %d to complete", task.ID)

	var response PrivateServiceConnectEndpoints
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) getCreationScript(ctx context.Context, message string, path string, values url.Values) (*CreationScript, error) {
	var task internal.TaskResponse
	err := a.client.GetWithQuery(ctx, message, path, values, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for private service connect creation script request %d to complete", task.ID)

	var response CreationScript
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) getDeletionScript(ctx context.Context, message string, path string) (*DeletionScript, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for private service connect deletion script request %d to complete", task.ID)

	var response DeletionScript
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) create(ctx context.Context, message string, path string) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, message, path, nil, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the Private Service Connect", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, fmt.Errorf("failed when creating Private Service Connect %d: %w", id, err)
	}

	return id, nil
}

func (a *API) update(ctx context.Context, message string, path string, body any) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, message, path, body, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish updating the Private Service Connect", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating Private Service Connect %w", err)
	}

	return nil
}

func (a *API) delete(ctx context.Context, message string, path string) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, message, path, nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish deleting the Private Service Connect", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting Private Service Connect %w", err)
	}

	return nil
}

func wrap404Error(subId int, err error) error {
	var e *internal.HTTPError
	if errors.As(err, &e) && e.StatusCode == http.StatusNotFound {
		return &NotFound{subscriptionID: subId}
	}
	var v *internal.Error
	if errors.As(err, &v) && v.StatusCode() == strconv.Itoa(http.StatusNotFound) {
		return &NotFound{subscriptionID: subId}
	}
	return err
}

func wrap404ErrorActiveActive(subId int, regionId int, err error) error {
	var e *internal.HTTPError
	if errors.As(err, &e) && e.StatusCode == http.StatusNotFound {
		return &NotFoundActiveActive{subscriptionID: subId, regionID: regionId}
	}
	var v *internal.Error
	if errors.As(err, &v) && v.StatusCode() == strconv.Itoa(http.StatusNotFound) {
		return &NotFoundActiveActive{subscriptionID: subId, regionID: regionId}
	}
	return err
}
