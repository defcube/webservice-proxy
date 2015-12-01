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

func (rs *Records) init() {
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

// getOrCreateRecord is a non-threadsafe function that gets or creates the record for the given url
func (rs *Records) getOrCreateRecord(url string) *record {
	r, ok := rs.urlRecordMap[url]
	if !ok {
		r = newRecord(url)
		rs.urlRecordMap[url] = r
	}
	return r
}

// addUrlRecordMapRequest requests that the function f be ran in a single thread.
// This is the safe way to interact with the urlRecordMap
func (rs *Records) addUrlRecordMapRequest(f func()) error {
	rs.init()
	select {
	case rs.requests <- f:
		return nil
	default:
		return WorkBufferFullError
	}

}

// RecordRequest will create or update a record. If the
// stats logging system can't keep up with the buffer, an error will
// be returned and this function should be re-engineered.
func (rs *Records) RecordRequest(url string) error {
	return rs.addUrlRecordMapRequest(func() {
		rs.getOrCreateRecord(url).addRequest()
	})
}

func (rs *Records) RecordClientHangup(url string) error {
	return rs.addUrlRecordMapRequest(func() {
		rs.getOrCreateRecord(url).addClientHangup()
	})
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
	rs.init()
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

func (rs *Records) NumClientHangups() int64 {
	rs.init()
	done := make(chan int64)
	rs.requests <- func() {
		result := int64(0)
		for _, v := range rs.urlRecordMap {
			result += v.NumClientHangups
		}
		done <- result
	}
	return <-done
}
