package stats

import (
	"time"
)

func newRecord(url string) *record {
	return &record{Url: url, CreatedAt: time.Now()}
}

type record struct {
	Url              string
	CreatedAt        time.Time
	TotalRequests    int64
	NumClientHangups int64
}

func (r *record) copy() *record {
	c := record{
		Url:              r.Url,
		CreatedAt:        r.CreatedAt,
		TotalRequests:    r.TotalRequests,
		NumClientHangups: r.NumClientHangups,
	}
	return &c
}

func (r *record) addRequest() {
	r.TotalRequests += 1
}

func (r *record) addClientHangup() {
	r.NumClientHangups += 1
}
