package internal

import (
	"fmt"
	"io"
	"net/http"
)

func send_put(remote string) (*http.Response, error) {
	return send("PUT", remote, nil)
}

func send_get(remote string) (*http.Response, error) {
	return send("GET", remote, nil)
}

func send_delete(remote string) (*http.Response, error) {
	return send("DELETE", remote, nil)
}

func send_copy(remote string) (*http.Response, error) {
	return send("COPY", remote, nil)
}

func send_lock(remote string) (*http.Response, error) {
	return send("LOCK", remote, nil)
}

func send_unlock(remote string) (*http.Response, error) {
	return send("UNLOCK", remote, nil)
}

func send_move(remote string) (*http.Response, error) {
	return send("MOVE", remote, nil)
}

func send_prop_find(remote string) (*http.Response, error) {
	return send("PROPFIND", remote, nil)
}

func send_prop_patch(remote string) (*http.Response, error) {
	return send("PROPPATCH", remote, nil)
}

func send_mkcol(remote string) (*http.Response, error) {
	return send("MKCOL", remote, nil)
}

func send(method, remote string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, remote, body)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("%s: wrong status code %d", method, resp.StatusCode)
	}
	return resp, nil
}
