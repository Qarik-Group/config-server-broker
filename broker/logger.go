package broker

import (
	"time"

	"code.cloudfoundry.org/lager"
)

type BrokerLogger struct {
	Logger lager.Logger
}

func (logger BrokerLogger) DisplayHeader(name string, value string) error {
	logger.Logger.Info("Header", lager.Data{
		"name":  name,
		"value": value})
	return nil
}

func (logger BrokerLogger) DisplayHost(name string) error {
	logger.Logger.Info("DisplayHost", lager.Data{
		"name": name,
	})
	return nil
}

func (logger BrokerLogger) DisplayJSONBody(body []byte) error {
	logger.Logger.Info("DisplayJSONBody", lager.Data{
		"body": body,
	})
	return nil
}

func (logger BrokerLogger) DisplayBody(body []byte) error {
	logger.Logger.Info("DisplayBody", lager.Data{
		"body": body,
	})
	return nil
}

func (logger BrokerLogger) DisplayMessage(msg string) error {
	logger.Logger.Info("DisplayMessage", lager.Data{
		"msg": msg,
	})
	return nil
}

func (logger BrokerLogger) DisplayRequestHeader(method string, uri string, httpProtocol string) error {
	logger.Logger.Info("DisplayRequestHeader", lager.Data{
		"method":       method,
		"uri":          uri,
		"httpProtocol": httpProtocol,
	})
	return nil
}

func (logger BrokerLogger) DisplayResponseHeader(httpProtocol string, status string) error {
	logger.Logger.Info("DisplayResponseHeader", lager.Data{
		"httpProtocol": httpProtocol,
		"status":       status,
	})
	return nil
}

func (logger BrokerLogger) DisplayType(name string, requestDate time.Time) error {
	logger.Logger.Info("DisplayType", lager.Data{
		"name":        name,
		"requestDate": requestDate,
	})
	return nil
}

func (logger BrokerLogger) HandleInternalError(err error) {
	logger.Logger.Error("HandleInternalError", err)
}

func (logger BrokerLogger) Start() error {
	logger.Logger.Info("Start")
	return nil
}

func (logger BrokerLogger) Stop() error {
	logger.Logger.Info("Stop")
	return nil
}
