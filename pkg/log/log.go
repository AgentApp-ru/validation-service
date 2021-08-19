package log

import (
	"net"
	"time"
	"validation_service/pkg/config"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	go func() {
		for {
			Logger.ReplaceHooks(getNewHooks())
			time.Sleep(60 * time.Second)
		}
	}()
}

func getNewHooks() logrus.LevelHooks {
	hooks := make(logrus.LevelHooks)
	if hook, err := getLogstashHook(); err == nil {
		hooks.Add(hook)
	}
	return hooks
}

func getLogstashHook() (logrus.Hook, error) {
	// TODO: В логгер тэгаи вкидывать PS + AgreementID, вместо того, чтоб писать внутрь message

	conn, err := net.Dial("tcp", config.Settings.LogstageUrl)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(
		logrus.Fields{"type": "validation"},
	))

	return hook, nil
}
