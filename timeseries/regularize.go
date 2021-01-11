package timeseries

import (
	"log"
	"math"
	"time"
)

// Main Method for downsampling a time series. Create and return a new time series.
// There is a parameter with the frequency of new sampling, but there is a choice to be made on the processing of data inside the sampling
// Currently these parameters are available: min,max, avg, and last.
func (ts *TimeSeries) Downsampling(freq int,per string, meth string, tolerance int)TimeSeries {
	switch per {
	case "Seconds":
		per= "s"
	case "s":
		per="s"
	case "sec":
		per="s"
	case "Minutes":
		per= "m"
	case "min":
		per="m"
	case "m":
		per="m"
	case "Hours":
		per= "h"
	case "h":
		per="h"
	default:
		log.Fatal("no period inserted")
	}
	var resampled TimeSeries
	ts.SortChronAsc()
	datemin := ts.DataSeries[0].Chron
	start := RoundedStartTime(datemin, freq,per)
	if datemin.Equal(start){start=AddDuration(datemin,-freq,per)}
	nextTicker := AddDurationTol(start, freq,per,0)
	nextTickerTol:= AddDurationTol(start, freq,per,tolerance)
	for i := 0; i < (len(ts.DataSeries) - 1); {
		sum := 0.0
		var vec []float64
		var obsmin, obsmax float64
		for (ts.DataSeries[i].Chron.Before(nextTickerTol) || ts.DataSeries[i].Chron.Equal(nextTickerTol) )&& i < (len(ts.DataSeries)-1) {
			sum += ts.DataSeries[i].Meas
			vec = append(vec, ts.DataSeries[i].Meas)
			i++
		}
		obsmin, obsmax = Bounds(vec)
		var obslast float64
		var du DataUnit
		if len(vec)!=0{
			obslast = vec[len(vec)-1]
			switch meth {
			case "Average","average","avg":
				du.Chron= nextTicker
				du.Meas=  Mean(vec)
				resampled.AddDataUnit(du)
			case "Maximum","maximum","max":
				du.Chron= nextTicker
				du.Meas=obsmax
				resampled.AddDataUnit(du)
			case "Minimum","minimum","min":
				du.Chron= nextTicker
				du.Meas=  obsmin
				resampled.AddDataUnit(du)
			case "Last","last":
				du.Chron= nextTicker
				du.Meas=  obslast
				resampled.AddDataUnit(du)
			case "Sum","sum":
				du.Chron= nextTicker
				du.Meas=  sum
				resampled.AddDataUnit(du)
			default:
				du.Chron= nextTicker
				du.Meas= NotNihil
			}
		}
		oldTicker:=nextTicker
		nextTicker = AddDurationTol(oldTicker, freq,per,0)
		nextTickerTol = AddDurationTol(oldTicker, freq,per,tolerance)

		// check if there are data between this tick and the next one. If not, insert a math.NaN().
		for (i< len(ts.DataSeries)-2) && ts.DataSeries[i].Chron.After(nextTickerTol) && (nextTickerTol.Equal(ts.DataSeries[i].Chron)==false)==true {
			du.Chron= nextTicker
			du.Meas= math.NaN()
			resampled.AddDataUnit(du)
			oldTicker=nextTicker
			nextTicker = AddDurationTol(oldTicker, freq,per,0)
			nextTickerTol = AddDurationTol(oldTicker, freq,per,tolerance)
		}
	}
	resampled.ComputeStats()
	return resampled
}
// Truncate a datetime to the closest beginning of time frequence but below. Ancillary to resampling methods.
func RoundedStartTime(timetoround time.Time, afreqq int, aper string) time.Time {
	roundedtime := time.Now()
	switch aper {
	case "m":
		roundedtime = timetoround.Truncate(time.Minute * time.Duration(afreqq))
	case "s":
		roundedtime = timetoround.Truncate(time.Second * time.Duration(afreqq))
	case "h":
		roundedtime = timetoround.Truncate(time.Hour * time.Duration(afreqq))
	case "d":
		roundedtime = timetoround.AddDate(0, 0, -afreqq)
	default:
		roundedtime = timetoround
	}
	return roundedtime
}
// Add duration to a given date. The parameter is a string consisting of an integer and one letter ("s" for seconds, "m" for minute, "h" for hour)
func AddDuration(start time.Time, freq int, per string) time.Time {
	switch per {
	case "s":
		return start.Add(time.Second*time.Duration(freq))
	case "m":
		return start.Add(time.Minute * time.Duration(freq))
	case "h":
		return start.Add(time.Hour * time.Duration(freq))
	default:
		return start
	}
}
// AddDurationTol add a duration plus a tolerance. Tolerance is an int. If the period is in seconds, tolerance is expressed
// in Millisecond. If the period is in Minutes, the tolerance is expressed in Seconds. If the period is in Hours, tolerance
// is expressed in Minutes. So if the regularisation is 30 minutes with a 3 minutes tolerance, 3 minutes should be expressed
// as 180
func AddDurationTol(start time.Time, freq int, per string,tolerance int) time.Time {
	switch per {
	case "s":
		return start.Add(time.Millisecond*time.Duration(freq*1000+tolerance))
	case "m":
		return start.Add(time.Second * time.Duration(freq*60+tolerance))
	case "h":
		return start.Add(time.Minute * time.Duration(freq*60+tolerance))
	default:
		return start
	}
}
// Computation of means without help of external library
func Mean(xs []float64) float64 {
	if len(xs) == 0.00000 {
		return NotNihil
	}
	m := 0.0
	for i, x := range xs {
		m += (x - m) / float64(i+1)
	}
	return m
}
// Computation of minimum and maximum without help of external library
func Bounds(xs []float64) (min float64, max float64) {
	if len(xs) == 0.00000000 {
		return 0.0,0.0 //math.NaN(), math.NaN()
	}
	min, max = xs[0], xs[0]
	for _, x := range xs {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return
}
