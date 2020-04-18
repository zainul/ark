package xsend

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/zainul/ark/xlog"

	"sync"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Go(ctx context.Context, request interface{}, response interface{}, statusCode int) (int, interface{}) {
	go func(req interface{}, res interface{}) {
		var wg sync.WaitGroup

		wg.Add(2)

		var bt []byte
		go func(waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			bt, _ = json.Marshal(res)
		}(&wg)

		var btreq []byte
		go func(waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			btreq, _ = json.Marshal(req)
		}(&wg)

		wg.Wait()

		xlog.Response(ctx, "response writer", btreq, bt)
	}(request, response)

	return statusCode, response
}
