package shutdown

import (
	"io"
	"os"
	"os/signal"
	"validation_service/pkg/log"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	log.Logger.Infof("Caught signal %s. Shutting down...", sig)

	// Here we can do graceful shutdown (close connections and etc)
	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			log.Logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
