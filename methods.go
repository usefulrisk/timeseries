package timeseries

import (
	"log"
	"math"
	"sort"
	"time"
)

// AddData adds DataUnit with data inputs. Empty fields are golang default values
func (ts *TimeSeries) AddData(chr time.Time, meas float64) {
	du := NewDataUnit(chr, meas)
	ts.DataSeries = append(ts.DataSeries, du)
}

// SortChronAsc sort a TimeSeries in chronological ascending order in-place
func (ts *TimeSeries) SortChronAsc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Chron.Before(ts.DataSeries[j].Chron)
	})
}

// SortChronDesc sort a TimeSeries in chronological descending order in-place
func (ts *TimeSeries) SortChronDesc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Chron.After(ts.DataSeries[j].Chron)
	})
}

// Sort a TimeSeries in ascending order of measure in-place
func (ts *TimeSeries) SortMeasAsc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Meas < ts.DataSeries[j].Meas
	})
}

// Sort a TimeSeries in descending order of measure in-place
func (ts *TimeSeries) SortMeasDesc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Meas > ts.DataSeries[j].Meas
	})
}

// Reset a TimeSeries to a zero length TimeSeries but keep metadata information at the TimeSeries level
func (ts *TimeSeries) Reset() {
	ts.DataSeries = ts.DataSeries[0:0]
}

// Deep copy of a time series
func (ts *TimeSeries) Copy() TimeSeries {
	var copyofts TimeSeries
	copyofts.Name = ts.Name
	copyofts.Comment = ts.Comment
	for _, value := range ts.DataSeries {
		copyofts.DataSeries = append(copyofts.DataSeries, value)
	}
	copyofts.Sort_Deltas_Stats()
	return copyofts
}

func (s *TimeSeries) SortChrono() {
	sort.Slice(s.DataSeries, func(i, j int) bool {
		return s.DataSeries[i].Chron.Before(s.DataSeries[j].Chron)
	})
}

// DeltasFiller is used in Sort_Deltas_Stats method. No need to use it alone. Warning. Does not sort data in chrnological order.
// Once Sort_Deltas_Stats function is applied, data are sorted.
func (ts *TimeSeries) DeltasFiller() {
	if len(ts.DataSeries) != 0 {
		ts.DataSeries[0].Dchron = 0 * time.Nanosecond
		ts.DataSeries[0].Dmeas = 0.0 //math.NaN()
		for i := 1; i < len(ts.DataSeries); i++ {
			ts.DataSeries[i].Dchron = ts.DataSeries[i].Chron.Sub(ts.DataSeries[i-1].Chron)
			ts.DataSeries[i].Dmeas = ts.DataSeries[i].Meas - ts.DataSeries[i-1].Meas
		}
	}
}

// Sort_Deltas_Stats computes Deltas and feeds TsStats struct
func (ts *TimeSeries) Sort_Deltas_Stats() {
	if len(ts.DataSeries) > 0 {
		// Preparation -----------------------------------
		ts.SortChronAsc()
		ts.DeltasFiller()
		ts.ComputeBasicStats()
	} else {
		ts.Comment = "Warning: Empty Time Series"
	}
}

func (ts *TimeSeries) ComputeBasicStats() {
	if len(ts.DataSeries) > 0 {
		for _, v := range ts.DataSeries {
			if math.IsNaN(v.Meas) == true {
				ts.NbreOfNaN += 1
			}
		}
		// stats on Chron ------------------------------------------
		ts.Chmin = ts.DataSeries[0].Chron
		ts.ValAtChmin = ts.DataSeries[0].Meas
		ts.Chmax = ts.DataSeries[len(ts.DataSeries)-1].Chron
		ts.ValAtChmax = ts.DataSeries[len(ts.DataSeries)-1].Meas
		var ChrVec []float64
		for _, val := range ts.DataSeries {
			ChrVec = append(ChrVec, float64(val.Chron.UnixNano()))
		}
		ts.Len = len(ChrVec)
		meanch, err := Mean(ChrVec)
		if err != nil {
			log.Println(err)
		}
		ts.Chmean = time.Unix(0, int64(meanch))
		medch, _ := Median(ChrVec)
		ts.Chmed = time.Unix(0, int64(medch))
		// stats on Dchron ------------------------------------------
		var DChrVec []float64
		for _, val := range ts.DataSeries {
			DChrVec = append(DChrVec, float64(val.Dchron))
		}
		DChrVec = DChrVec[1:]
		mindcfl, _ := Min(DChrVec)
		ts.DChmin = time.Duration(int64(mindcfl))
		maxdcfl, _ := Max(DChrVec)
		ts.DChmax = time.Duration(int64(maxdcfl))
		meandcfl, _ := Mean(DChrVec)
		ts.DChmean = time.Duration(int64(meandcfl))
		meddcfl, _ := Median(DChrVec)
		ts.DChmed = time.Duration(int64(meddcfl))
		// stats on Meas --------------------------------------------
		var MeasVec []float64
		for _, val := range ts.DataSeries {
			if math.IsNaN(val.Meas) == false {
				MeasVec = append(MeasVec, val.Meas)
			}
		}
		ts.Msmin, _ = Min(MeasVec)
		ts.Msmax, _ = Max(MeasVec)
		for _, v := range ts.DataSeries {
			if v.Meas == ts.Msmin {
				ts.ChAtMsmin = v.Chron
			}
			if v.Meas == ts.Msmax {
				ts.ChAtMsmax = v.Chron
			}
		}
		ts.Msmean, _ = Mean(MeasVec)
		ts.Msmed, _ = Median(MeasVec)
		ts.Msstd, _ = StdDev(MeasVec)
		// stats on DMeas ---------------------------------------------
		var DMeasVec []float64
		for _, val := range ts.DataSeries {
			if math.IsNaN(val.Dmeas) == false {
				DMeasVec = append(DMeasVec, val.Dmeas)
			}
		}
		if len(DMeasVec) > 1 {
			DMeasVec = DMeasVec[1:]
			ts.DMsmin, _ = Min(DMeasVec)
			ts.DMsmax, _ = Max(DMeasVec)
			ts.DMsmean, _ = Mean(DMeasVec)
			ts.DMsmed, _ = Median(DMeasVec)
			ts.DMsstd, _ = StdDev(DMeasVec)
			// end of process ---------------------------------
			ts.Comment = " Time Series ok."
		}

	} else {
		ts.Comment = "Warning: Empty Time Series"
	}
}

func NewTsContainer() (tsd TsContainer) {
	tsd.Ts = make(map[string]*TimeSeries)
	return
}
