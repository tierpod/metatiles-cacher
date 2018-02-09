package fetch

import (
	"log"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
)

type fetchJob struct {
	index int
	url   string
}

type fetchResult struct {
	index int
	data  []byte
	err   error
}

func fetchWorker(jobs <-chan fetchJob, results chan<- fetchResult, shutdown <-chan interface{}, cfg config.HTTPClient, logger *log.Logger) {
	logger.Printf("[DEBUG] start fetch worker")

	// create new httpclient with internal connection pool
	httpc := httpclient.New(cfg.Headers, cfg.Timeout)

	for j := range jobs {
		select {
		case <-shutdown:
			logger.Printf("[ERROR] shutdown fetch worker (via shutdown channel)")
			return
		default:
			// logger.Printf("[DEBUG] fetch %v", j.url)
			body, err := httpc.GetBody(j.url)

			// test error
			// if j.index == 5 {
			// 	results <- fetchResult{index: j.index, data: nil, err: fmt.Errorf("test error")}
			// 	return
			// }

			if err != nil {
				results <- fetchResult{index: j.index, data: nil, err: err}
				logger.Printf("[ERROR] fetch: %v", err)
				return
			}

			// debug slow connections
			// time.Sleep(1 * time.Second)

			results <- fetchResult{index: j.index, data: body, err: nil}
		}
	}

	logger.Printf("[DEBUG] shutdown fetch worker")
}
