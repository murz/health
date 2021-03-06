package health

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonPollingSink(t *testing.T) {
	setNowMock("2011-09-09T23:36:13Z")
	defer resetNowMock()

	sink := NewJsonPollingSink(time.Minute, time.Minute*5)

	sink.EmitEvent("myjob", "myevent", nil)
	sink.EmitEventErr("myjob", "myevent", errors.New("myerr"), nil)
	sink.EmitTiming("myjob", "myevent", 100, nil)
	sink.EmitComplete("myjob", Success, 9, nil)

	time.Sleep(10 * time.Millisecond) // we need to make sure we process the above metrics before we get the metrics.
	intervals := sink.GetMetrics()

	sink.ShutdownServer()

	assert.Equal(t, 1, len(intervals))

	intAgg := intervals[0]
	assert.Equal(t, 1, intAgg.Events["myevent"])
	assert.Equal(t, 1, intAgg.EventErrs["myevent"].Count)
	assert.Equal(t, 1, intAgg.Timers["myevent"].Count)
	assert.Equal(t, 1, intAgg.Jobs["myjob"].Count)
}
