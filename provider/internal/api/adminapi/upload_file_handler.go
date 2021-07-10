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
	"os"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// UploadFileHandler handles upload file request.
func UploadFileHandler(data []byte) (byte, []byte, error) {
	logging.Debug("Handle upload file from admin")
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	tag, fileData, err := fcradminmsg.DecodeUploadFileStartRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	if _, err := os.Stat(filepath.Join(c.Settings.RetrievalDir, tag)); os.IsNotExist(err) {
		// Not exist, save
		f, err := os.Create(filepath.Join(c.Settings.RetrievalDir, tag))
		if err == nil {
			_, err = f.Write(fileData)
			f.Close()
		}
		if err != nil {
			err = fmt.Errorf("Error saving file: %v", err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		// Save cid, tag
		reader, err := os.Open(filepath.Join(c.Settings.RetrievalDir, tag))
		if err != nil {
			err = fmt.Errorf("Fail to open file for cid calculation %v: %v", tag, err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		cid, err := cid.NewContentIDFromFile(reader)
		if err != nil {
			err = fmt.Errorf("Invalid CID: %v", err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		c.OfferMgr.AddCIDTag(cid, tag)

	} else {
		// Exist
		err = fmt.Errorf("Filename already existed %v", tag)
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Succeed
	ack := fcradminmsg.EncodeACK(true, "Succeed")
	return fcradminmsg.ACKType, ack, nil
}
