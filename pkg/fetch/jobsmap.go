package fetch

import (
	"sync"

	"github.com/tierpod/go-osm/metatile"
)

// jobsMap incapsulate mutex and map of jobs.
type jobsMap struct {
	mu   sync.RWMutex
	jobs map[metatile.Metatile]interface{}
}

// newJobsMap creates new jobsMap
func newJobsMap() *jobsMap {
	jobs := make(map[metatile.Metatile]interface{})
	return &jobsMap{jobs: jobs}
}

// add adds metatile in jobsMap if does not exists (otherwise do nothing).
func (j *jobsMap) add(mt metatile.Metatile) {
	if j.exists(mt) {
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()
	j.jobs[mt] = nil
}

// delete deletes metatile from jobsMap.
func (j *jobsMap) delete(mt metatile.Metatile) {
	j.mu.Lock()
	defer j.mu.Unlock()

	delete(j.jobs, mt)
}

// exists checks if metatile inside jobsMap.
func (j *jobsMap) exists(mt metatile.Metatile) bool {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if _, found := j.jobs[mt]; found {
		return true
	}

	return false
}

// items returns list of jobs inside jobsMap.
func (j *jobsMap) items() []metatile.Metatile {
	j.mu.RLock()
	defer j.mu.RUnlock()

	var items []metatile.Metatile
	for i := range j.jobs {
		items = append(items, i)
	}

	return items
}
