package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LogRequest(request *http.Request, response int) {
	method := request.Method
	endpoint := request.URL
	ip := request.RemoteAddr

	if response >= 500 && response < 600 {
		logrus.Errorf("%s %s resulted in %d from %s", method, endpoint, response, ip)
		return
	}

	if response >= 400 && response < 500 {
		logrus.Warnf("%s %s resulted in %d from %s", method, endpoint, response, ip)
		return
	}

	logrus.Infof("%s %s resulted in %d from %s", method, endpoint, response, ip)

}
