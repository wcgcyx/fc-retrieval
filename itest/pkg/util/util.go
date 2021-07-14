/*
Package util - common functions used in end-to-end and integration testing. Allowing to start different types of
Retrieval network nodes for testing.
*/
package util

func GetLotusAPI() string {
	return ""
}

func GetRegisterAPI() string {
	return ""
}

func GetLotusToken() (string, string) {
	return "", ""
}

func Topup(lotusAPI string, token string, superAcct string, privKeys []string) {
}

func GetContainerInfo(pvd bool) ([]string, []string) {
	return nil, nil
}
