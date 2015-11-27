package stats

import (
	"math/big"
	"time"
)

var bigOne = (&big.Int{}).SetInt64(1)

func newRecord(url string) *record {
	return &record{Url: url, CreatedAt: time.Now()}
}

type record struct {
	Url           string
	CreatedAt     time.Time
	TotalRequests big.Int
}

func (r *record) copy() *record {
	c := record{
		Url:       r.Url,
		CreatedAt: r.CreatedAt,
	}
	c.TotalRequests.Set(&r.TotalRequests)
	return &c
}

func (r *record) addRequest() {
	r.TotalRequests.Add(&r.TotalRequests, bigOne)
}
