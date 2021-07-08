package adminapi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

func UploadFileHandler(data []byte) (byte, []byte, error) {
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
