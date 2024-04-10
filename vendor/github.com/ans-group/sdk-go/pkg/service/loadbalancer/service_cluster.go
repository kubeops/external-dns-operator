package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetClusters retrieves a list of clusters
func (s *Service) GetClusters(parameters connection.APIRequestParameters) ([]Cluster, error) {
	return connection.InvokeRequestAll(s.GetClustersPaginated, parameters)
}

// GetClustersPaginated retrieves a paginated list of clusters
func (s *Service) GetClustersPaginated(parameters connection.APIRequestParameters) (*connection.Paginated[Cluster], error) {
	body, err := s.getClustersPaginatedResponseBody(parameters)
	return connection.NewPaginated(body, parameters, s.GetClustersPaginated), err
}

func (s *Service) getClustersPaginatedResponseBody(parameters connection.APIRequestParameters) (*connection.APIResponseBodyData[[]Cluster], error) {
	body := &connection.APIResponseBodyData[[]Cluster]{}

	response, err := s.connection.Get("/loadbalancers/v2/clusters", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetCluster retrieves a single cluster by id
func (s *Service) GetCluster(clusterID int) (Cluster, error) {
	body, err := s.getClusterResponseBody(clusterID)

	return body.Data, err
}

func (s *Service) getClusterResponseBody(clusterID int) (*connection.APIResponseBodyData[Cluster], error) {
	body := &connection.APIResponseBodyData[Cluster]{}

	if clusterID < 1 {
		return body, fmt.Errorf("invalid cluster id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/loadbalancers/v2/clusters/%d", clusterID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ClusterNotFoundError{ID: clusterID}
		}

		return nil
	})
}

// PatchCluster patches a Cluster
func (s *Service) PatchCluster(clusterID int, req PatchClusterRequest) error {
	_, err := s.patchClusterResponseBody(clusterID, req)

	return err
}

func (s *Service) patchClusterResponseBody(clusterID int, req PatchClusterRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if clusterID < 1 {
		return body, fmt.Errorf("invalid cluster id")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/loadbalancers/v2/clusters/%d", clusterID), &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ClusterNotFoundError{ID: clusterID}
		}

		return nil
	})
}

// DeployCluster deploys a Cluster
func (s *Service) DeployCluster(clusterID int) error {
	_, err := s.deployClusterResponseBody(clusterID)

	return err
}

func (s *Service) deployClusterResponseBody(clusterID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if clusterID < 1 {
		return body, fmt.Errorf("invalid cluster id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/loadbalancers/v2/clusters/%d/deploy", clusterID), nil)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ClusterNotFoundError{ID: clusterID}
		}

		return nil
	})
}

// ValidateCluster validates a cluster
func (s *Service) ValidateCluster(clusterID int) error {
	response := &connection.APIResponse{}

	if clusterID < 1 {
		return fmt.Errorf("invalid cluster id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/loadbalancers/v2/clusters/%d/validate", clusterID), connection.APIRequestParameters{})
	if err != nil {
		return err
	}

	if response.StatusCode == 422 {
		body := &validateClusterResponseBody{}

		return errors.New(body.Error())
	}

	return response.HandleResponse(&connection.APIResponseBody{}, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ClusterNotFoundError{ID: clusterID}
		}

		return nil
	})
}

type validateClusterResponseBody struct {
	connection.APIResponseBody
}

// GetCluster retrieves a single cluster by id
func (s *Service) GetClusterACLTemplates(clusterID int) (ACLTemplates, error) {
	body, err := s.getClusterACLTemplatesResponseBody(clusterID)

	return body.Data, err
}

func (s *Service) getClusterACLTemplatesResponseBody(clusterID int) (*connection.APIResponseBodyData[ACLTemplates], error) {
	body := &connection.APIResponseBodyData[ACLTemplates]{}

	if clusterID < 1 {
		return body, fmt.Errorf("invalid cluster id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/loadbalancers/v2/clusters/%d/acl-templates", clusterID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ClusterNotFoundError{ID: clusterID}
		}

		return nil
	})
}
