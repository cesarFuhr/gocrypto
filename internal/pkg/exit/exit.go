package exit

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cesarFuhr/gocrypto/internal/pkg/logger"
	"go.uber.org/zap"
)

//ListenToExit notify when a sign to exit was made
func ListenToExit(e chan struct{}) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	l := logger.NewLogger()

	go exitListener(l, s, e)
}

func exitListener(l logger.Logger, sigs chan os.Signal, exit chan struct{}) {
	sig := <-sigs
	l.Info("Shuting down...", zap.Stringer("signal", sig))

	exit <- struct{}{}
}
