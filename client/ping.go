package client
import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/vietnamz/cli-common/api/types"
	"github.com/vietnamz/cli-common/errdefs"
)

// Ping pings the server and returns the value of the "Docker-Experimental",
// "Builder-Version", "OS-Type" & "API-Version" headers. It attempts to use
// a HEAD request on the endpoint, but falls back to GET if HEAD is not supported
// by the daemon.
func (cli *Client) Ping(ctx context.Context) (types.Ping, error) {
	var ping types.Ping

	// Using cli.buildRequest() + cli.doRequest() instead of cli.sendRequest()
	// because ping requests are used during API version negotiation, so we want
	// to hit the non-versioned /_ping endpoint, not /v1.xx/_ping
	req, err := cli.buildRequest(http.MethodHead, path.Join(cli.basePath, "/ping"), nil, nil)
	if err != nil {
		return ping, err
	}
	serverResp, err := cli.doRequest(ctx, req)
	if err == nil {
		defer ensureReaderClosed(serverResp)
		switch serverResp.statusCode {
		case http.StatusOK, http.StatusInternalServerError:
			// Server handled the request, so parse the response
			return parsePingResponse(cli, serverResp)
		}
	} else if IsErrConnectionFailed(err) {
		return ping, err
	}

	req, err = cli.buildRequest(http.MethodGet, path.Join(cli.basePath, "/ping"), nil, nil)
	if err != nil {
		return ping, err
	}
	serverResp, err = cli.doRequest(ctx, req)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return ping, err
	}
	return parsePingResponse(cli, serverResp)
}

func parsePingResponse(cli *Client, resp serverResponse) (types.Ping, error) {
	var ping types.Ping
	if resp.header == nil {
		err := cli.checkResponseErr(resp)
		return ping, errdefs.FromStatusCode(err, resp.statusCode)
	}
	err :=json.NewDecoder(resp.body).Decode(&ping)
	return ping, err
}

