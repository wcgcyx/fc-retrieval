/*
Package adminapi contains the API code for the admin client - provider communication.
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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// ListFilesHandler lists all the files.
func ListFilesHandler(data []byte) (byte, []byte, error) {
	logging.Debug("Handle list files from admin")
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	files, err := ioutil.ReadDir(c.Settings.RetrievalDir)
	if err != nil {
		err = fmt.Errorf("Error reading files from directory %v: %v", c.Settings.RetrievalDir, err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	tags := make([]string, 0)
	cids := make([]string, 0)
	sizes := make([]int64, 0)
	published := make([]bool, 0)
	frequency := make([]int, 0)
	for _, file := range files {
		tag := file.Name()
		tags = append(tags, tag)
		cidStr := c.OfferMgr.GetCIDByTag(tag)
		if cidStr == "" {
			reader, err := os.Open(filepath.Join(c.Settings.RetrievalDir, tag))
			if err != nil {
				err = fmt.Errorf("Fail to open file %v: %v", file, err.Error())
				ack := fcradminmsg.EncodeACK(false, err.Error())
				return fcradminmsg.ACKType, ack, err
			}
			cid, err := cid.NewContentIDFromFile(reader)
			if err != nil {
				err = fmt.Errorf("Invalid CID: %v", err.Error())
				ack := fcradminmsg.EncodeACK(false, err.Error())
				return fcradminmsg.ACKType, ack, err
			}
			cidStr = cid.ToString()
			c.OfferMgr.AddCIDTag(cid, tag)
		}
		cid, err := cid.NewContentID(cidStr)
		if err != nil {
			err = fmt.Errorf("Internal error in parsing cid string %v: %v", cidStr, err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		cids = append(cids, cidStr)
		sizes = append(sizes, file.Size())
		published = append(published, len(c.OfferMgr.GetOffers(cid)) > 0)
		frequency = append(frequency, c.OfferMgr.GetAccessCountByCID(cid))
	}

	response, err := fcradminmsg.EncodeListFilesResponse(tags, cids, sizes, published, frequency)
	if err != nil {
		err = fmt.Errorf("Error in encoding response: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	return fcradminmsg.ListFilesResponseType, response, nil
}
