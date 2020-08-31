package system

import (
	"context"
	"github.com/vietnamz/cli-common/api/server/httputils"
	"github.com/vietnamz/cli-common/api/types"
	"net/http"
	"runtime"
)

func (s *systemRouter) pingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")

	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Length", "0")
	}
	ping := types.Ping{
		APIVersion: "v1",
		OSType: runtime.GOOS,
		Experimental: false,
		BuilderVersion: "HelloWorld",
	}
	err := httputils.WriteJSON(w, 200, ping)
	return err
}
