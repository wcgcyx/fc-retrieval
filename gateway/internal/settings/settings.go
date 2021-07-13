/*
Package settings - holds configuration specific to a Retrieval Gateway node.
*/
package settings

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
	"math/big"
	"time"
)

// DefaultMsgKeyUpdateDuration is the default msg signing key update duration
const DefaultMsgKeyUpdateDuration = 12 * time.Hour

// DefaultSyncDuration is the default peer manager sync duration
const DefaultSyncDuration = 12 * time.Hour

// DefaultTCPInactivityTimeout is the default timeout for TCP inactivity
const DefaultTCPInactivityTimeout = 5000 * time.Millisecond

// DefaultLongTCPInactivityTimeout is the default timeout for long TCP inactivity. This timeout should never be ignored.
const DefaultLongTCPInactivityTimeout = 300000 * time.Millisecond

// AppSettings defines the server configuraiton
type AppSettings struct {
	// Logging related settings
	LogServiceName string `mapstructure:"LOG_SERVICE_NAME"` // Log service name
	LogLevel       string `mapstructure:"LOG_LEVEL"`        // Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget      string `mapstructure:"LOG_TARGET"`       // Log Level: STDOUT
	LogDir         string `mapstructure:"LOG_DIR"`          // Log Dir: /var/.fc-retrieval/log
	LogFile        string `mapstructure:"LOG_FILE"`         // Log File: gateway.log
	LogMaxBackups  int    `mapstructure:"LOG_MAX_BACKUPS"`  // Log max backups: 3
	LogMaxAge      int    `mapstructure:"LOG_MAX_AGE"`      // Log max age (days): 28
	LogMaxSize     int    `mapstructure:"LOG_MAX_SIZE"`     // Log max size (MB): 500
	LogCompress    bool   `mapstructure:"LOG_COMPRESS"`     // Log compress: false
	LogTimeFormat  string `mapstructure:"LOG_TIME_FORMAT"`  // Log time format: RFC3339

	// Admin related
	BindAdminAPI   int    `mapstructure:"BIND_ADMIN_API"`   // Port number to bind to for admin secured HTTP connection
	SystemDir      string `mapstructure:"SYSTEM_DIR"`       // // Dir storing all data of this gateway
	RetrievalDir   string `mapstructure:"RETRIEVAL_DIR"`    // Retrieval Dir: /var/.fc-retrieval/gateway/files
	AdminKeyFile   string `mapstructure:"ADMIN_KEY_FILE"`   // File storing the admin access key file
	ConfigFile     string `mapstructure:"CONFIG_FILE"`      // File storing the gateway config
	StoreFullOffer bool   `mapstructure:"STORE_FULL_OFFER"` // Boolean indicates whether this gateway stores full offer

	// Duration
	SyncDuration             time.Duration `mapstructure:"SYNC_DURATION"`               // Sync duration
	MsgKeyUpdateDuration     time.Duration `mapstructure:"MSG_KEY_UPDATE_DURATION"`     // Msg key update duration
	TCPInactivityTimeout     time.Duration `mapstructure:"TCP_INACTIVITY_TIMEOUT"`      // TCP inactivity timeout
	TCPLongInactivityTimeout time.Duration `mapstructure:"TCP_LONG_INACTIVITY_TIMEOUT"` // TCP long inactivity timeout

	// Price, this is not configurable at the moment.
	SearchPrice *big.Int `mapstructure:"SEARCH_PRICE"` // Search price
	OfferPrice  *big.Int `mapstructure:"OFFER_PRICE"`  // Offer price
	TopupAmount *big.Int `mapstructure:"TOPUP_AMOUNT"` // Topup amount
}
