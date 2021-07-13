package log

import (
	"net"
	"validation_service/pkg/config"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	if config.Settings.Env == "production" {
		addLogstashHook()
	}
}

func addLogstashHook() {
	// TODO: В логгер тэгаи вкидывать PS + AgreementID, вместо того, чтоб писать внутрь message

	conn, err := net.Dial("tcp", config.Settings.LogstageUrl)
	if err != nil {
		Logger.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(
		logrus.Fields{"type": "validation"},
		))

	Logger.Hooks.Add(hook)
}
