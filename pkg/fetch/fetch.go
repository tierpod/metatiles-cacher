// Package fetch provides background fetch service.
package fetch

import (
	"log"
	"time"

	"github.com/tierpod/go-osm/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

// CacheWriter provides interface for writing data to metatile cache.
type CacheWriter interface {
	Write(mt metatile.Metatile, data [][]byte) error
}

type job struct {
	mt      metatile.Metatile
	urlTmpl string
}

// Service is the basic fetch service struct.
type Service struct {
	logger  *log.Logger
	cfg     *config.Config
	queue   chan (job)
	jobsMap *jobsMap
	cw      CacheWriter
}

// NewService creates new fetch Service.
func NewService(cfg *config.Config, cw CacheWriter, logger *log.Logger) *Service {
	queue := make(chan job, cfg.Fetch.Buffer)
	jm := newJobsMap()
	s := &Service{
		logger:  logger,
		cfg:     cfg,
		queue:   queue,
		jobsMap: jm,
		cw:      cw,
	}

	return s
}

// Start starts fetch Service in background.
func (s *Service) Start() {
	go func() {
		for {
			job := <-s.queue
			go s.process(job)
		}
	}()
}

// Add adds job for metatile `mt` and url template `URLTmpl` to fetching queue. Skip if item already
// in queue.
func (s *Service) Add(mt metatile.Metatile, URLTmpl string) {
	if s.jobsMap.exists(mt) {
		s.logger.Printf("[DEBUG] skip job %v: already in process", mt)
		return
	}

	select {
	case s.queue <- job{mt: mt, urlTmpl: URLTmpl}:
	// TODO: configure timeout?
	case <-time.After(10 * time.Second):
		s.logger.Printf("[ERROR] unable to add job to queue, timeout exceeded")
	}
}

// Jobs returns jobs who are currently in process.
func (s *Service) Jobs() []metatile.Metatile {
	return s.jobsMap.items()
}

// process starts processing job `j`.
func (s *Service) process(j job) error {
	start := time.Now()
	s.logger.Printf("[DEBUG] start job: %+v", j)

	s.jobsMap.add(j.mt)
	defer func() {
		s.jobsMap.delete(j.mt)
		s.logger.Printf("[DEBUG] end job: %+v", j)
	}()

	data, err := s.fetch(j.mt, j.urlTmpl)
	if err != nil {
		s.logger.Printf("[ERROR] fetch: %v", err)
		return err
	}

	err = s.cw.Write(j.mt, data)
	if err != nil {
		s.logger.Printf("[ERROR] write: %v", err)
		return err
	}

	elapsed := time.Since(start)
	s.logger.Printf("[INFO] fetch and write complete in %v", elapsed)

	return nil
}

func (s *Service) fetch(mt metatile.Metatile, URLTmpl string) ([][]byte, error) {
	xx, yy := mt.XYBox()
	s.logger.Printf("[INFO] fetch style(%v) z(%v) xx(%v-%v) yy(%v-%v)", mt.Style, mt.Zoom,
		xx[0], xx[len(xx)-1], yy[0], yy[len(yy)-1])

	count := mt.Size() * mt.Size()

	jobs := make(chan workerJob, count)
	results := make(chan workerResult, count)
	shutdown := make(chan interface{})

	for w := 0; w < s.cfg.Fetch.Workers; w++ {
		go worker(jobs, results, shutdown, s.cfg.HTTPClient, s.logger)
	}

	data := make([][]byte, metatile.Area)
	for _, x := range xx {
		for _, y := range yy {
			i := metatile.XYToIndex(x, y)
			url := util.MakeURL(URLTmpl, mt.Zoom, x, y)
			jobs <- workerJob{index: i, url: url}
		}
	}
	close(jobs)

	for w := 0; w < count; w++ {
		r := <-results
		if r.err != nil {
			close(shutdown)
			return nil, r.err
		}
		data[r.index] = r.data
	}
	close(shutdown)

	return data, nil
}
