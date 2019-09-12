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
		logrus.WithFields(logrus.Fields{"method": method, "endpoint": endpoint, "ip": ip, "response": response}).Error("Request failed due to server side error")
		return
	}

	if response >= 400 && response < 500 {
		logrus.WithFields(logrus.Fields{"method": method, "endpoint": endpoint, "ip": ip, "response": response}).Warn("Request impossible to handle")
		return
	}

	logrus.WithFields(logrus.Fields{"method": method, "endpoint": endpoint, "ip": ip, "response": response}).Info("Request sucessfully handled")

}
