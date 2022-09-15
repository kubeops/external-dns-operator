package alidns

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CopyGtmConfig invokes the alidns.CopyGtmConfig API synchronously
func (client *Client) CopyGtmConfig(request *CopyGtmConfigRequest) (response *CopyGtmConfigResponse, err error) {
	response = CreateCopyGtmConfigResponse()
	err = client.DoAction(request, response)
	return
}

// CopyGtmConfigWithChan invokes the alidns.CopyGtmConfig API asynchronously
func (client *Client) CopyGtmConfigWithChan(request *CopyGtmConfigRequest) (<-chan *CopyGtmConfigResponse, <-chan error) {
	responseChan := make(chan *CopyGtmConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CopyGtmConfig(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CopyGtmConfigWithCallback invokes the alidns.CopyGtmConfig API asynchronously
func (client *Client) CopyGtmConfigWithCallback(request *CopyGtmConfigRequest, callback func(response *CopyGtmConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CopyGtmConfigResponse
		var err error
		defer close(result)
		response, err = client.CopyGtmConfig(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CopyGtmConfigRequest is the request struct for api CopyGtmConfig
type CopyGtmConfigRequest struct {
	*requests.RpcRequest
	SourceId     string `position:"Query" name:"SourceId"`
	TargetId     string `position:"Query" name:"TargetId"`
	CopyType     string `position:"Query" name:"CopyType"`
	UserClientIp string `position:"Query" name:"UserClientIp"`
	Lang         string `position:"Query" name:"Lang"`
}

// CopyGtmConfigResponse is the response struct for api CopyGtmConfig
type CopyGtmConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCopyGtmConfigRequest creates a request to invoke CopyGtmConfig API
func CreateCopyGtmConfigRequest() (request *CopyGtmConfigRequest) {
	request = &CopyGtmConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "CopyGtmConfig", "alidns", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCopyGtmConfigResponse creates a response to parse from CopyGtmConfig response
func CreateCopyGtmConfigResponse() (response *CopyGtmConfigResponse) {
	response = &CopyGtmConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
