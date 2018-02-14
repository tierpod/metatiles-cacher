package fetch

import (
	"log"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type workerJob struct {
	index int
	url   string
}

type workerResult struct {
	index int
	data  []byte
	err   error
}

func worker(jobs <-chan workerJob, results chan<- workerResult, shutdown <-chan interface{}, cfg config.HTTPClient, logger *log.Logger) {
	wid := util.RandomString(4)
	logger.Printf("[DEBUG] (fetch) <%v> start worker", wid)

	// create new httpclient with internal connection pool
	httpc := httpclient.New(cfg.Headers, cfg.Timeout)

	for j := range jobs {
		select {
		case <-shutdown:
			logger.Printf("[ERROR] (fetch) <%v> shutdown worker (via shutdown channel)", wid)
			return
		default:
			logger.Printf("[DEBUG] (fetch) <%v> download %v", wid, j.url)
			body, err := httpc.GetBody(j.url)

			// test error
			// if j.index == 5 {
			// 	results <- workerResult{index: j.index, data: nil, err: fmt.Errorf("test error")}
			// 	return
			// }

			// test slow connections
			// time.Sleep(1 * time.Second)

			results <- workerResult{index: j.index, data: body, err: err}
		}
	}

	logger.Printf("[DEBUG] (fetch) <%v> shutdown worker", wid)
}
