package safedns

import (
	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetSettings retrieves account settings
func (s *Service) GetSettings() (Settings, error) {
	body, err := s.getSettingsResponseBody()

	return body.Data, err
}

func (s *Service) getSettingsResponseBody() (*connection.APIResponseBodyData[Settings], error) {
	body := &connection.APIResponseBodyData[Settings]{}

	response, err := s.connection.Get("/safedns/v1/settings", connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}
