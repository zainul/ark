package xsend

import (
	"context"
	"encoding/json"
	"github.com/zainul/poskit/pkg/xlog"
	"net/http"
)

// Write is to make response of http
func Write(ctx context.Context, w http.ResponseWriter, request interface{}, response interface{}, statusCode ...int) {
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}

	go func() {
		bt, _ := json.Marshal(response)
		btreq, _ := json.Marshal(request)
		xlog.Response(ctx, "response writer", btreq, bt)
	}()

	json.NewEncoder(w).Encode(response)
}
