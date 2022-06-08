package internal

import (
	"fmt"
	"io"
	"net/http"
)

var (
	GoodRespCodes = []int{http.StatusOK, http.StatusAccepted, http.StatusCreated}
)

func send_put(srv string, buf io.Reader) (*http.Response, error) {
	return send(http.MethodPut, srv, buf, GoodRespCodes)
}

func send_head(srv string) (*http.Response, error) {
	return send(http.MethodHead, srv, nil, GoodRespCodes)
}

func send_get(srv string) (*http.Response, error) {
	return send(http.MethodGet, srv, nil, GoodRespCodes)
}

func send_delete(srv string) (*http.Response, error) {
	return send(http.MethodDelete, srv, nil, GoodRespCodes)
}

func send_copy(srv string) (*http.Response, error) {
	return send("COPY", srv, nil, GoodRespCodes)
}

func send_lock(srv string) (*http.Response, error) {
	return send("LOCK", srv, nil, GoodRespCodes)
}

func send_unlock(srv string) (*http.Response, error) {
	return send("UNLOCK", srv, nil, GoodRespCodes)
}

func send_move(srv string) (*http.Response, error) {
	return send("MOVE", srv, nil, GoodRespCodes)
}

func send_prop_find(srv string) (*http.Response, error) {
	return send("PROPFIND", srv, nil, GoodRespCodes)
}

func send_prop_patch(srv string) (*http.Response, error) {
	return send("PROPPATCH", srv, nil, GoodRespCodes)
}

func send_mkcol(srv string) (*http.Response, error) {
	return send("MKCOL", srv, nil, GoodRespCodes)
}

func send(method, srv string, body io.Reader, goodResponseCodes []int) (*http.Response, error) {
	req, err := http.NewRequest(method, srv, body)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if !isGoodResponseCode(resp.StatusCode, goodResponseCodes) {
		return nil, fmt.Errorf("%s: bad status code %d", method, resp.StatusCode)
	}
	// if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
	// 	return nil, fmt.Errorf("%s: wrong status code %d", method, resp.StatusCode)
	// }
	return resp, nil
}

func isGoodResponseCode(responseCode int, goodResponseCodes []int) bool {
	for _, b := range goodResponseCodes {
		if b == responseCode {
			return true
		}
	}
	return false
}
