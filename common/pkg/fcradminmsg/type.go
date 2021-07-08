/*
Package fcradminmsg - stores all the admin messages.
*/
package fcradminmsg

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
	InitialisationRequestType     = 0
	ListPeersRequestType          = 1
	ListPeersResponseType         = 2
	InspectPeerRequestType        = 3
	InspectPeerResponseType       = 5
	ChangePeerStatusRequestType   = 6
	ListCIDFrequencyRequestType   = 12
	ListCIDFrequencyResponseType  = 13
	GetOfferByCIDRequestType      = 14
	GetOfferByCIDResponseType     = 15
	CacheOfferByDigestRequestType = 16
	ListFilesRequestType          = 17
	ListFilesResponseType         = 18
	PublishOfferRequestType       = 19
	UploadFileRequestType         = 20
	ACKType                       = 21
)
