package logging

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// LoggerFromConfig initializes a service-specific logger instance.
func Configure(name string) log.Logger {
	return log.NewLogger(
		log.Name(name),
		log.Level("debug"),
		log.Pretty(true),
		log.Color(true),
		log.File(""),
	)
}
