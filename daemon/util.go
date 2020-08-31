package daemon

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetPublicIpAddresses() string  {
	resp, err := http.Get("http://ipecho.net/plain")
	if err != nil {
		// handle error
		logrus.Errorf("failed to get ip address %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to get ip address %s", err)
	}
	return string(body)
}

func GetLocalIpAddresses()  []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Errorf("failed to get ip address list %s", err)
	}
	var res []string
	for _, i := range  addrs {
		if !strings.Contains(i.String(), ":") && strings.Contains(i.String(), "127.0.0.1")  {
			res = append(res, strings.Split(i.String(), "/")[0])
		}
	}
	return res
}


