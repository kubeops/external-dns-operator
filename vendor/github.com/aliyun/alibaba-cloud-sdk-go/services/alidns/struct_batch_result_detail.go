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

// BatchResultDetail is a nested struct in alidns response
type BatchResultDetail struct {
	Status         bool   `json:"Status" xml:"Status"`
	Type           string `json:"Type" xml:"Type"`
	Domain         string `json:"Domain" xml:"Domain"`
	Remark         string `json:"Remark" xml:"Remark"`
	RecordId       string `json:"RecordId" xml:"RecordId"`
	Rr             string `json:"Rr" xml:"Rr"`
	Priority       string `json:"Priority" xml:"Priority"`
	RrStatus       string `json:"RrStatus" xml:"RrStatus"`
	OperateDateStr string `json:"OperateDateStr" xml:"OperateDateStr"`
	NewValue       string `json:"NewValue" xml:"NewValue"`
	Value          string `json:"Value" xml:"Value"`
	Ttl            string `json:"Ttl" xml:"Ttl"`
	BatchType      string `json:"BatchType" xml:"BatchType"`
	Line           string `json:"Line" xml:"Line"`
	NewRr          string `json:"NewRr" xml:"NewRr"`
	Reason         string `json:"Reason" xml:"Reason"`
}
