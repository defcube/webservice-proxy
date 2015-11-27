package stats

import (
	"errors"
	"sync"
)

type Records struct {
	initOnce sync.Once

	// requests is the only way to manipulate the urlRecordMap in a thread safe way
	requests     chan (func())
	urlRecordMap map[string]*record
}

func (rs *Records) Init() {
	rs.initOnce.Do(func() {
		rs.urlRecordMap = make(map[string]*record)
		rs.requests = make(chan (func()), 1000)
		go rs.processRequests()
	})
}

func (rs *Records) processRequests() {
	for r := range rs.requests {
		r()
	}
}

var WorkBufferFullError = errors.New("Work buffer is full")

// RecordRequest will create or update a record. If the
// stats logging system can't keep up with the buffer, an error will
// be returned and this function should be re-engineered.
func (rs *Records) RecordRequest(url string) error {
	rs.Init()
	addFunc := func() {
		r, ok := rs.urlRecordMap[url]
		if !ok {
			r = newRecord(url)
			rs.urlRecordMap[url] = r
		}
		r.addRequest()
	}
	select {
	case rs.requests <- addFunc:
		return nil
	default:
		return WorkBufferFullError
	}
}

func (rs *Records) MustRecordRequest(url string) {
	err := rs.RecordRequest(url)
	if err != nil {
		panic(err)
	}
}

// List provides threadsafe access to the internal list. Raises an
// error if buffers are full
func (rs *Records) List() map[string]*record {
	done := make(chan map[string]*record)
	rs.requests <- func() {
		newMap := make(map[string]*record)
		for k, v := range rs.urlRecordMap {
			newMap[k] = v.copy()
		}
		done <- newMap
	}
	return <-done
}
