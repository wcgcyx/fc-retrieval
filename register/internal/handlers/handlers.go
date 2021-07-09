/*
Package handlers - is remote API to call a Register in a FileCoin Network. Contains method exposed by Register to be
called by Retrieval Gateways and Retrieval Providers in order to get information about current network state.

E.g. - get an updated list of gateways, including just recently enrolled ones
*/
package handlers

import (
	"github.com/wcgcyx/fc-retrieval/register/config"
)

var apiconfig = config.Config()
