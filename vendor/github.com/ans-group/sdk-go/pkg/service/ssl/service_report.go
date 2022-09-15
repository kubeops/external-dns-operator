package ssl

import (
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
)

// GetReport retrieves a single report by id
func (s *Service) GetReport(domainName string) (Report, error) {
	body, err := s.getReportResponseBody(domainName)

	return body.Data, err
}

func (s *Service) getReportResponseBody(domainName string) (*connection.APIResponseBodyData[Report], error) {
	body := &connection.APIResponseBodyData[Report]{}

	if domainName == "" {
		return body, fmt.Errorf("invalid domain name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ssl/v1/reports/%s", domainName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}
