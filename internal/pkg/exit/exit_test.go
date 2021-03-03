package exit

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
)

type loggerStub struct {
	CalledWith []interface{}
	t          *testing.T
	wg         *sync.WaitGroup
}

func (l *loggerStub) Info(s string, args ...zap.Field) {
	arguments := append([]interface{}{}, s, args)
	l.CalledWith = arguments
	l.t.Log(l.CalledWith...)
	l.wg.Done()
}

func TestExitListener(t *testing.T) {
	wg := sync.WaitGroup{}
	log := &loggerStub{[]interface{}{}, t, &wg}
	t.Run("closes the channel when receives a notification", func(t *testing.T) {
		wg.Add(1)

		sigs := make(chan os.Signal)
		exit := make(chan struct{})

		go exitListener(log, sigs, exit)

		sigs <- os.Interrupt
		got := <-exit

		wg.Wait()

		if !reflect.DeepEqual(got, struct{}{}) {
			t.Errorf("should receive a empty struct")
		}
	})
	t.Run("calls logger when receives a notification", func(t *testing.T) {
		wg.Add(1)

		sigs := make(chan os.Signal)
		exit := make(chan struct{})

		go exitListener(log, sigs, exit)

		sigs <- os.Interrupt
		got := <-exit

		wg.Wait()

		if !reflect.DeepEqual(got, struct{}{}) {
			t.Errorf("should receive a empty struct")
		}
		if len(log.CalledWith) == 0 {
			t.Errorf("Should have called logger, %v", log.CalledWith)
		}
	})
}
