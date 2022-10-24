package ecloud

import (
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetVPNEndpoints retrieves a list of VPN endpoints
func (s *Service) GetVPNEndpoints(parameters connection.APIRequestParameters) ([]VPNEndpoint, error) {
	return connection.InvokeRequestAll(s.GetVPNEndpointsPaginated, parameters)
}

// GetVPNEndpointsPaginated retrieves a paginated list of VPN endpoints
func (s *Service) GetVPNEndpointsPaginated(parameters connection.APIRequestParameters) (*connection.Paginated[VPNEndpoint], error) {
	body, err := s.getVPNEndpointsPaginatedResponseBody(parameters)
	return connection.NewPaginated(body, parameters, s.GetVPNEndpointsPaginated), err
}

func (s *Service) getVPNEndpointsPaginatedResponseBody(parameters connection.APIRequestParameters) (*connection.APIResponseBodyData[[]VPNEndpoint], error) {
	body := &connection.APIResponseBodyData[[]VPNEndpoint]{}

	response, err := s.connection.Get("/ecloud/v2/vpn-endpoints", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetVPNEndpoint retrieves a single VPN endpoint by id
func (s *Service) GetVPNEndpoint(endpointID string) (VPNEndpoint, error) {
	body, err := s.getVPNEndpointResponseBody(endpointID)

	return body.Data, err
}

func (s *Service) getVPNEndpointResponseBody(endpointID string) (*connection.APIResponseBodyData[VPNEndpoint], error) {
	body := &connection.APIResponseBodyData[VPNEndpoint]{}

	if endpointID == "" {
		return body, fmt.Errorf("invalid vpn endpoint id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v2/vpn-endpoints/%s", endpointID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &VPNEndpointNotFoundError{ID: endpointID}
		}

		return nil
	})
}

// CreateVPNEndpoint creates a new VPN endpoint
func (s *Service) CreateVPNEndpoint(req CreateVPNEndpointRequest) (TaskReference, error) {
	body, err := s.createVPNEndpointResponseBody(req)

	return body.Data, err
}

func (s *Service) createVPNEndpointResponseBody(req CreateVPNEndpointRequest) (*connection.APIResponseBodyData[TaskReference], error) {
	body := &connection.APIResponseBodyData[TaskReference]{}

	response, err := s.connection.Post("/ecloud/v2/vpn-endpoints", &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// PatchVPNEndpoint patches a VPN endpoint
func (s *Service) PatchVPNEndpoint(endpointID string, req PatchVPNEndpointRequest) (TaskReference, error) {
	body, err := s.patchVPNEndpointResponseBody(endpointID, req)

	return body.Data, err
}

func (s *Service) patchVPNEndpointResponseBody(endpointID string, req PatchVPNEndpointRequest) (*connection.APIResponseBodyData[TaskReference], error) {
	body := &connection.APIResponseBodyData[TaskReference]{}

	if endpointID == "" {
		return body, fmt.Errorf("invalid endpoint id")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/ecloud/v2/vpn-endpoints/%s", endpointID), &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &VPNEndpointNotFoundError{ID: endpointID}
		}

		return nil
	})
}

// DeleteVPNEndpoint deletes a VPN endpoint
func (s *Service) DeleteVPNEndpoint(endpointID string) (string, error) {
	body, err := s.deleteVPNEndpointResponseBody(endpointID)

	return body.Data.TaskID, err
}

func (s *Service) deleteVPNEndpointResponseBody(endpointID string) (*connection.APIResponseBodyData[TaskReference], error) {
	body := &connection.APIResponseBodyData[TaskReference]{}

	if endpointID == "" {
		return body, fmt.Errorf("invalid endpoint id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v2/vpn-endpoints/%s", endpointID), nil)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &VPNEndpointNotFoundError{ID: endpointID}
		}

		return nil
	})
}

// GetVPNEndpointTasks retrieves a list of VPN endpoint tasks
func (s *Service) GetVPNEndpointTasks(endpointID string, parameters connection.APIRequestParameters) ([]Task, error) {
	return connection.InvokeRequestAll(func(p connection.APIRequestParameters) (*connection.Paginated[Task], error) {
		return s.GetVPNEndpointTasksPaginated(endpointID, p)
	}, parameters)
}

// GetVPNEndpointTasksPaginated retrieves a paginated list of VPN endpoint tasks
func (s *Service) GetVPNEndpointTasksPaginated(endpointID string, parameters connection.APIRequestParameters) (*connection.Paginated[Task], error) {
	body, err := s.getVPNEndpointTasksPaginatedResponseBody(endpointID, parameters)

	return connection.NewPaginated(body, parameters, func(p connection.APIRequestParameters) (*connection.Paginated[Task], error) {
		return s.GetVPNEndpointTasksPaginated(endpointID, p)
	}), err
}

func (s *Service) getVPNEndpointTasksPaginatedResponseBody(endpointID string, parameters connection.APIRequestParameters) (*connection.APIResponseBodyData[[]Task], error) {
	body := &connection.APIResponseBodyData[[]Task]{}

	if endpointID == "" {
		return body, fmt.Errorf("invalid vpn endpoint id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v2/vpn-endpoints/%s/tasks", endpointID), parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &VPNEndpointNotFoundError{ID: endpointID}
		}

		return nil
	})
}
