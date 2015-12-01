package stats_test

import (
	"fmt"
	"github.com/defcube/webservice-proxy/server/internal/stats"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRecordRequest(t *testing.T) {
	records := stats.Records{}
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			records.MustRecordRequest("http://foo.bar/")
			records.MustRecordRequest(fmt.Sprintf("http://foo.bar/%v", num))
			wg.Done()
		}(i)
	}
	wg.Wait()

	// get the records list, then manipulate records more, then make sure nothing changed on our copy
	l := records.List()
	records.MustRecordRequest("http://foo.bar/")
	records.MustRecordRequest("http://foo.bar/after1")
	assert.Len(t, l, 101)
	assert.Equal(t, int64(100), l["http://foo.bar/"].TotalRequests)
	assert.Equal(t, int64(101), records.List()["http://foo.bar/"].TotalRequests)
}
