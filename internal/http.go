package internal

import (
	"errors"
	"net/http"
	"strings"
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
			if err := kv.UploadFile(key, r.Body, r.ContentLength); err != nil {
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
