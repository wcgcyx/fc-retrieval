/*
Package adminapi contains the API code for the admin client - gateway communication.
*/
package adminapi

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
	"errors"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// ForceSyncHandler handles force sync handler
func ForceSyncHandler(data []byte) (byte, []byte, error) {
	logging.Debug("Handle force sync from admin")
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	min, max := c.PeerMgr.GetCurrentCIDHashRange()
	logging.Debug("Before syncing min: %v, max: %v", min, max)

	// Do a sync
	c.PeerMgr.Sync()

	min, max = c.PeerMgr.GetCurrentCIDHashRange()
	logging.Debug("After syncing min: %v, max: %v", min, max)

	// Succeed
	ack := fcradminmsg.EncodeACK(true, "Succeed.")
	return fcradminmsg.ACKType, ack, nil
}
