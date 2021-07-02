/*
Package fcrmessages - stores all the p2p messages.
*/
package fcrmessages

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

const (
	StandardOfferDiscoveryRequestType  = 0
	StandardOfferDiscoveryResponseType = 1
	DHTOfferDiscoveryRequestType       = 2
	DHTOfferDiscoveryResponseType      = 3
	OfferPublishRequestType            = 4
	DataRetrievalRequestType           = 5 // Placeholder, TBD
	DataRetrievalResponseType          = 6 // Placeholder, TBD
	ACKType                            = 7
)
