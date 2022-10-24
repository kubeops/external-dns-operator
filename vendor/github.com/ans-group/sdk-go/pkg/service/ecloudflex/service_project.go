package ecloudflex

import (
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetProjects retrieves a list of projects
func (s *Service) GetProjects(parameters connection.APIRequestParameters) ([]Project, error) {
	return connection.InvokeRequestAll(s.GetProjectsPaginated, parameters)
}

// GetProjectsPaginated retrieves a paginated list of projects
func (s *Service) GetProjectsPaginated(parameters connection.APIRequestParameters) (*connection.Paginated[Project], error) {
	body, err := s.getProjectsPaginatedResponseBody(parameters)
	return connection.NewPaginated(body, parameters, s.GetProjectsPaginated), err
}

func (s *Service) getProjectsPaginatedResponseBody(parameters connection.APIRequestParameters) (*connection.APIResponseBodyData[[]Project], error) {
	body := &connection.APIResponseBodyData[[]Project]{}

	response, err := s.connection.Get("/ecloud-flex/v1/projects", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetProject retrieves a single project by id
func (s *Service) GetProject(projectID int) (Project, error) {
	body, err := s.getProjectResponseBody(projectID)

	return body.Data, err
}

func (s *Service) getProjectResponseBody(projectID int) (*connection.APIResponseBodyData[Project], error) {
	body := &connection.APIResponseBodyData[Project]{}

	if projectID < 1 {
		return body, fmt.Errorf("invalid project id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud-flex/v1/projects/%d", projectID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ProjectNotFoundError{ID: projectID}
		}

		return nil
	})
}
