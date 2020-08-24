package server

import (
	"github.com/sirupsen/logrus"
	"github.com/vietnamz/cli-common/api/server/httputils"
	"github.com/vietnamz/cli-common/api/server/middleware"
)

func (s *Server) handlerWithGlobalMiddlewares( handler httputils.APIFunc) httputils.APIFunc {
	next := handler
	for _, m := range s.middlewares {
		next = m.WrapHandler(next)
	}
	if s.cfg.Logging && logrus.GetLevel() == logrus.DebugLevel {
		next = middleware.DebugRequestMiddleware(next)
	}
	return next
}
