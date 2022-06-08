package internal

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/h2non/bimg"
)

func (kv *KVImg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := []byte(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		{
			rec, found, err := kv.GetFile(key)
			w.Header().Set("Content-Length", "0")
			if len(rec.Hash) > 0 {
				w.Header().Set("Content-Md5", rec.Hash)
			}
			if len(rec.Locations) > 0 {
				w.Header().Set("X-Srv-Locations", strings.Join(rec.Locations, ","))
			}
			if rec.Deleted == DELETED {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				if errors.Is(err, ErrNotMatching) {
					w.Header().Set("X-Srv-Status", "needs-rebuild")
				}
				if errors.Is(err, ErrImgDoesntExist) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			w.Header().Set("Location", getId(rec.Locations[found], keyToPath(key)))
			w.WriteHeader(http.StatusFound)
			break
		}
	case http.MethodPut:
		{
			if r.ContentLength == 0 {
				w.WriteHeader(http.StatusLengthRequired)
				return
			}
			// Check if header sent to optimize image
			var buf io.Reader = r.Body
			if r.Header.Get("X-Optimize") == "true" {
				// TODO: Optimize the buffer image
				data, err := io.ReadAll(buf)
				if err != nil {
					w.Header().Set("Error", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				res, err := bimg.NewImage(data).Convert(bimg.WEBP)
				if err != nil {
					w.Header().Set("Error", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				buf = bytes.NewReader(res)
			}
			if err := kv.UploadFile(key, buf, r.ContentLength); err != nil {
				w.Header().Set("Error", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			break
		}
	case http.MethodDelete:
		{
			if err := kv.DeleteFile(key); err != nil {
				w.Header().Set("Error", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			break
		}
	default:
		{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}
