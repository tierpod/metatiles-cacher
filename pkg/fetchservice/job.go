package fetchservice

import "github.com/tierpod/metatiles-cacher/pkg/coords"

// Job contains metatile coordinates, style and source.
type Job struct {
	Meta   coords.Metatile
	Style  string
	Source string
}

// NewJob creates new job.
func NewJob(meta coords.Metatile, style, source string) Job {
	return Job{
		Meta:   meta,
		Style:  style,
		Source: source,
	}
}
