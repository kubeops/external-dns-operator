package loadbalancer

import (
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetACLs retrieves a list of ACLs
// Currently, a target_group_id or listener_id filter must be provided for this to return data
func (s *Service) GetACLs(parameters connection.APIRequestParameters) ([]ACL, error) {
	return connection.InvokeRequestAll(s.GetACLsPaginated, parameters)
}

// GetACLsPaginated retrieves a paginated list of ACLs
// Currently, a target_group_id or listener_id filter must be provided for this to return data
func (s *Service) GetACLsPaginated(parameters connection.APIRequestParameters) (*connection.Paginated[ACL], error) {
	body, err := s.getACLsPaginatedResponseBody(parameters)
	return connection.NewPaginated(body, parameters, s.GetACLsPaginated), err
}

func (s *Service) getACLsPaginatedResponseBody(parameters connection.APIRequestParameters) (*connection.APIResponseBodyData[[]ACL], error) {
	body := &connection.APIResponseBodyData[[]ACL]{}

	response, err := s.connection.Get("/loadbalancers/v2/acls", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetACL retrieves a single ACL by id
func (s *Service) GetACL(aclID int) (ACL, error) {
	body, err := s.getACLResponseBody(aclID)

	return body.Data, err
}

func (s *Service) getACLResponseBody(aclID int) (*connection.APIResponseBodyData[ACL], error) {
	body := &connection.APIResponseBodyData[ACL]{}

	if aclID < 1 {
		return body, fmt.Errorf("invalid acl id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/loadbalancers/v2/acls/%d", aclID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ACLNotFoundError{ID: aclID}
		}

		return nil
	})
}

// CreateACL creates an ACL
func (s *Service) CreateACL(req CreateACLRequest) (int, error) {
	body, err := s.createACLResponseBody(req)

	return body.Data.ID, err
}

func (s *Service) createACLResponseBody(req CreateACLRequest) (*connection.APIResponseBodyData[ACL], error) {
	body := &connection.APIResponseBodyData[ACL]{}

	response, err := s.connection.Post("/loadbalancers/v2/acls", &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body)
}

// PatchACL patches an ACL
func (s *Service) PatchACL(aclID int, req PatchACLRequest) error {
	_, err := s.patchACLResponseBody(aclID, req)

	return err
}

func (s *Service) patchACLResponseBody(aclID int, req PatchACLRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if aclID < 1 {
		return body, fmt.Errorf("invalid acl id")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/loadbalancers/v2/acls/%d", aclID), &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ACLNotFoundError{ID: aclID}
		}

		return nil
	})
}

// DeleteACL deletes an ACL
func (s *Service) DeleteACL(aclID int) error {
	_, err := s.deleteACLResponseBody(aclID)

	return err
}

func (s *Service) deleteACLResponseBody(aclID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if aclID < 1 {
		return body, fmt.Errorf("invalid acl id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/loadbalancers/v2/acls/%d", aclID), nil)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ACLNotFoundError{ID: aclID}
		}

		return nil
	})
}
