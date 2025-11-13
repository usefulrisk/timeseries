package timeseries

import (
	"math/rand"
	"time"
)

func BulkSimul(name string, from time.Time, period time.Duration, samplesize int, moy float64, stdev float64, jitter time.Duration) (ts TimeSeries) {
	ts.Name = name
	timecreated := from
	r0 := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < samplesize; i++ {
		lotto := r0.NormFloat64() * float64(jitter.Nanoseconds())
		peradd := time.Duration(lotto) + period
		timecreated = timecreated.Add(peradd)
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		du := DataUnit{
			Chron: timecreated,
			Meas:  r1.NormFloat64()*stdev + moy,
		}
		ts.AddDataUnit(du)
	}
	return ts
}
