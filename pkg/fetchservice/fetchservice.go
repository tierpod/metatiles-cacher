// Package fetchservice starts bagroundd service for fetching tiles from remote source and saves it to disk in
// metatile format.
package fetchservice

import (
	"fmt"
	"log"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
)

// FetchService is the structure for background fetch and write service
type FetchService struct {
	WriteCh    chan Job
	FetchQueue *queue.Uniq
	logger     *log.Logger
	cw         cache.Writer
}

// NewFetchService creates new FetchService
func NewFetchService(buffer int, cw cache.Writer, logger *log.Logger) *FetchService {
	fq := queue.NewUniq()
	//wch := make(chan coords.Metatile, buffer)
	wch := make(chan Job)
	fs := &FetchService{
		WriteCh:    wch,
		FetchQueue: fq,
		logger:     logger,
		cw:         cw,
	}

	logger.Printf("FetchService: Starting background FetchService")
	go fs.start()
	return fs
}

func (fs *FetchService) start() {
	for {
		job := <-fs.WriteCh
		fs.logger.Printf("[DEBUG] FetchService: Received %v from writer channel", job)
		go fs.fetchAndWrite(job)
	}
}

func (fs *FetchService) fetchAndWrite(j Job) error {
	var result [][]byte
	var url, zxy string

	defer func() {
		fs.logger.Printf("[DEBUG] FetchService/fetchAndWrite: Done %v, delete from active queue", j.Meta)
		fs.FetchQueue.Del(j.Meta.Path())
	}()

	minX, minY := j.Meta.MinXY()
	fs.logger.Printf("FetchService: Fetch Style(%v) Z(%v) X(%v-%v) Y(%v-%v)", j.Style, j.Meta.Z, minX, minX+j.Meta.Size(), minY, minY+j.Meta.Size())

	for _, x := range j.Meta.ConvertToXYBox().X {
		for _, y := range j.Meta.ConvertToXYBox().Y {
			zxy = strconv.Itoa(j.Meta.Z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png"
			url = fmt.Sprintf(j.Source, zxy)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url)
			if err != nil {
				fs.logger.Printf("[ERROR] FetchService/fetchAndWrite: %v", err)
				return fmt.Errorf("FetchService/fetchAndWrite: %v", err)
			}
			result = append(result, res)
		}
	}

	err := fs.cw.Write(j.Meta, j.Style, result)
	if err != nil {
		return fmt.Errorf("FetchService/fetchAndWrite: %v", err)
	}

	return nil
}

// Add checks FetchQueue and add metatile to active queue if not exist
func (fs *FetchService) Add(j Job) {
	// TODO: limiter?
	if !fs.FetchQueue.Add(j.Meta.Path()) {
		log.Printf("[DEBUG] FetchService: Skip %v, in the active queue", j)
		return
	}
	select {
	case fs.WriteCh <- j:
		fs.logger.Printf("[DEBUG] FetchService/Add: Send %v to writer channel", j)
	default:
		fs.logger.Printf("[ERROR] FetchService/Add: Writer channel is full")
	}
}
