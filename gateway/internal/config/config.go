/*
Package config - combines operations used to setup parameters for Gateway node in FileCoin network
*/
package config

import (
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/settings"
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

// Map sets the config for the Gateway. NB: Gateways start without a private key. Private keys are provided by a gateway admin client.
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

	defaultSearchPrice := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("SEARCH_PRICE"), defaultSearchPrice)
	if err != nil {
		// defaultSearchPrice is the default search price "0.001".
		defaultSearchPrice = big.NewInt(1_000_000_000_000_000)
	}

	defaultOfferPrice := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("OFFER_PRICE"), defaultOfferPrice)
	if err != nil {
		// defaultOfferPrice is the default offer price "0.001".
		defaultOfferPrice = big.NewInt(1_000_000_000_000_000)
	}

	defaultTopUpAmount := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("TOPUP_AMOUNT"), defaultTopUpAmount)
	if err != nil {
		// defaultTopUpAmount is the default top up amount "0.1".
		defaultTopUpAmount = big.NewInt(100_000_000_000_000_000)
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

		BindAdminAPI:   conf.GetInt("BIND_ADMIN_API"),
		AdminKeyFile:   conf.GetString("ADMIN_KEY_FILE"),
		RetrievalDir:   conf.GetString("RETRIEVAL_DIR"),
		StoreFullOffer: conf.GetBool("STORE_FULL_OFFER"),

		SyncDuration:             syncDuration,
		MsgKeyUpdateDuration:     msgKeyUpdateDuration,
		TCPInactivityTimeout:     tcpInactivityTimeout,
		TCPLongInactivityTimeout: tcpLongInactivityTimeout,

		SearchPrice: defaultSearchPrice,
		OfferPrice:  defaultOfferPrice,
		TopupAmount: defaultTopUpAmount,
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
