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

// Record is a nested struct in alidns response
type Record struct {
	Value      string `json:"Value" xml:"Value"`
	TTL        int64  `json:"TTL" xml:"TTL"`
	Remark     string `json:"Remark" xml:"Remark"`
	RR         string `json:"RR" xml:"RR"`
	DomainName string `json:"DomainName" xml:"DomainName"`
	Priority   int64  `json:"Priority" xml:"Priority"`
	RecordId   string `json:"RecordId" xml:"RecordId"`
	Status     string `json:"Status" xml:"Status"`
	Weight     int    `json:"Weight" xml:"Weight"`
	Locked     bool   `json:"Locked" xml:"Locked"`
	Line       string `json:"Line" xml:"Line"`
	Type       string `json:"Type" xml:"Type"`
}
