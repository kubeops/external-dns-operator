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

// RecoveryPlan is a nested struct in alidns response
type RecoveryPlan struct {
	Status                string `json:"Status" xml:"Status"`
	LastRollbackTimestamp int64  `json:"LastRollbackTimestamp" xml:"LastRollbackTimestamp"`
	UpdateTime            string `json:"UpdateTime" xml:"UpdateTime"`
	Remark                string `json:"Remark" xml:"Remark"`
	CreateTime            string `json:"CreateTime" xml:"CreateTime"`
	RecoveryPlanId        int64  `json:"RecoveryPlanId" xml:"RecoveryPlanId"`
	UpdateTimestamp       int64  `json:"UpdateTimestamp" xml:"UpdateTimestamp"`
	LastExecuteTimestamp  int64  `json:"LastExecuteTimestamp" xml:"LastExecuteTimestamp"`
	LastExecuteTime       string `json:"LastExecuteTime" xml:"LastExecuteTime"`
	LastRollbackTime      string `json:"LastRollbackTime" xml:"LastRollbackTime"`
	Name                  string `json:"Name" xml:"Name"`
	FaultAddrPoolNum      int    `json:"FaultAddrPoolNum" xml:"FaultAddrPoolNum"`
	CreateTimestamp       int64  `json:"CreateTimestamp" xml:"CreateTimestamp"`
}
