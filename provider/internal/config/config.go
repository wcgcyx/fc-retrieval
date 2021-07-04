/*
Package config - combines operations used to setup parameters for Provider node in FileCoin network
*/
package config

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
	"flag"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/settings"
)

// NewConfig creates a new configuration
func NewConfig() *viper.Viper {
	conf := viper.New()
	conf.AutomaticEnv()
	defineFlags(conf)
	bindFlags(conf)
	setValues(conf)
	return conf
}

// Map sets the config for the Provider. NB: Providers start without a private key. Private keys are provided by a provider admin client.
func Map(conf *viper.Viper) settings.AppSettings {
	syncDuration, err := time.ParseDuration(conf.GetString("SYNC_DURATION"))
	if err != nil {
		syncDuration = settings.DefaultSyncDuration
	}
	msgKeyUpdateDuration, err := time.ParseDuration(conf.GetString("MSG_KEY_UPDATE_DURATION"))
	if err != nil {
		msgKeyUpdateDuration = settings.DefaultMsgKeyUpdateDuration
	}
	tcpInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpInactivityTimeout = settings.DefaultTCPInactivityTimeout
	}
	tcpLongInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_LONG_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpLongInactivityTimeout = settings.DefaultLongTCPInactivityTimeout
	}

	return settings.AppSettings{
		LogServiceName: conf.GetString("LOG_SERVICE_NAME"),
		LogLevel:       conf.GetString("LOG_LEVEL"),
		LogTarget:      conf.GetString("LOG_TARGET"),
		LogDir:         conf.GetString("LOG_DIR"),
		LogFile:        conf.GetString("LOG_FILE"),
		LogMaxBackups:  conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:      conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:     conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:    conf.GetBool("LOG_COMPRESS"),
		LogTimeFormat:  conf.GetString("LOG_TIME_FORMAT"),

		BindAdminAPI: conf.GetInt("BIND_ADMIN_API"),
		AdminKeyFile: conf.GetString("ADMIN_KEY_FILE"),
		RetrievalDir: conf.GetString("RETRIEVAL_DIR"),

		SyncDuration:             syncDuration,
		MsgKeyUpdateDuration:     msgKeyUpdateDuration,
		TCPInactivityTimeout:     tcpInactivityTimeout,
		TCPLongInactivityTimeout: tcpLongInactivityTimeout,
	}
}

func defineFlags(conf *viper.Viper) {
	flag.String("host", "0.0.0.0", "help message for host")
	flag.String("ip", "127.0.0.1", "help message for ip")
}

func bindFlags(conf *viper.Viper) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := conf.BindPFlags(pflag.CommandLine); err != nil {
		logging.Error("can't bind a command line flag")
	}
}

func setValues(conf *viper.Viper) {
	conf.Set("IP", conf.GetString("ip"))
}
