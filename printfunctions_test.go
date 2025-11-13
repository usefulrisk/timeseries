package timeseries

import (
	"testing"
	"time"
)

func TestTimeSeries_PrettyPrint(t *testing.T) {
	ts := BulkSimul("yep", time.Now(), 10*time.Second, 100, 100.0, 20.0, 0)
	ts.Sort_Deltas_Stats()
	ts.PrettyPrint()
}
func TestTsC_PrettyPrint(t *testing.T) {
	var tsc = NewTsContainer()
	ts := BulkSimul("yep", time.Now(), 10*time.Second, 100, 100.0, 20.0, 0)
	ts.Sort_Deltas_Stats()
	tsc.Ts["ori"] = &ts
	tsr := tsc.Ts["ori"].Regularize(30, "sec", "avg", 0)
	tsc.Ts["reg"] = &tsr
	tsr.Sort_Deltas_Stats()
	tsc.PrettyPrint()
	tsc.ToJSON()
}
