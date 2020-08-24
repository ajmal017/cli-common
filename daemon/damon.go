package daemon

import (
	"github.com/vietnamz/cli-comoon/daemon/config"
)

// Daemon is entry to keep all the backend services to serve the API.
type Daemon struct {
	 config *config.Config
}

// Constructor.
func NewDaemon(cfg *config.Config) *Daemon {
	return &Daemon{
		config: cfg,
	}
}

// initialize all the backend services.
// Support:
//			+ Prime Generator service: to return a sample prime number.
func (d *Daemon) Init() error{
	return nil
}