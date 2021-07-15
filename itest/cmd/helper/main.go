/*
Package main - containing code for helper method for running test in container
*/
package main

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"fmt"

	"github.com/wcgcyx/fc-retrieval/itest/pkg/util"
)

func main() {
	lotusAPI := util.GetLotusAPI()
	registerAPI := util.GetRegisterAPI()
	lotusToken, superAcct := util.GetLotusToken()
	res := fmt.Sprintf("%v;%v;%v;%v", lotusAPI, registerAPI, lotusToken, superAcct)
	ips := util.GetContainerInfo(false)
	for _, ip := range ips {
		res = fmt.Sprintf("%v;%v", res, ip)
	}
	ips = util.GetContainerInfo(true)
	for _, ip := range ips {
		res = fmt.Sprintf("%v;%v", res, ip)
	}
	fmt.Println(res)
}
