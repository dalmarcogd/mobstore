package reqgetter

import (
	"net/http"
)

func GetCid(req *http.Request) *string {
	if req != nil && req.Header != nil {
		xCid := req.Header.Get("x-cid")
		if xCid != "" {
			return &xCid
		}
	}
	return nil
}
