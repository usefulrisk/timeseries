package timeseries

import (
	"math"
	"math/rand"
	"time"
)

// Simulate instantaneously a TimeSeries following the parameters:
//
// name: obvious
//
// Id: must be of format map[string]string. The simplest way is to write a literal: 'map[string]string{"Unit":"C","room":"kitchen"}'
//
//from: first date in the simulation. Must be in time.Time format.
//
//lag: lag between two simulation points. Must be in time.Duration For example, time.Second*2. If an integer is used it will be
// understood as nanosecond: so 1000000000 (1,000,000,000) = 1 sec. It may be the easiest if one forget the special golang way to note
// durations
//
//jitter: move of clock around the correct new time. Simulate errors and delays. Warning, in time.Duration
//
//iterations: length of the time series simulated.
func  Simul(name string,from time.Time, lag time.Duration,jitter time.Duration,iterations int)TimeSeries{
	var tss TimeSeries
	tss.Name=name
	for i := 0; i < iterations; i++ {
		rand.Seed(time.Now().UnixNano())
		var dum DataUnit
		dum.Chron=from.Add(time.Duration(i)*lag+time.Duration(rand.NormFloat64())* jitter)
		dum.Meas = rand.NormFloat64()*20 + 100
		tss.AddDataUnit(dum)
	}
	return tss
}
func  SimulWithNaN(name string,from time.Time, lag time.Duration,jitter time.Duration,iterations int)TimeSeries{
	var tss TimeSeries
	tss.Name=name
	for i := 0; i < iterations; i++ {
		rand.Seed(time.Now().UnixNano())
		var dum DataUnit
		dum.Chron=from.Add(time.Duration(i)*lag+time.Duration(rand.NormFloat64())* jitter)
		dum.Meas = rand.NormFloat64()*20 + 100
		if i==5 || i==10{
			dum.Meas=math.NaN()
		}
		tss.AddDataUnit(dum)
	}
	return tss
}