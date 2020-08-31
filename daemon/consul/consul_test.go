package consul

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestConsulMonitor(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	consul := NewConsulMonitor(DefaultServiceConfig())
	consul.Start()
	consul.Register()
	time.Sleep(200 * time.Second)
	consul.Stop()
}

