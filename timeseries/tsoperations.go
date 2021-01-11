package timeseries

import (
	"math"
	"time"
)

// ConvertTimeSeries2TsInArrays converts TimeSeries into list of Arrays that are easier for some
// computations and necessary for graphics. Warning: In order to compute stats
// with the montanaflynn package, missing data (math.NaN() ) are not imported in float64 arrays. Vectors have therefore
// different length. It is not correct to assume that Chron.length and Meas.length have the same length.
func(ts *TimeSeries) ConvertTimeSeries2TsInArrays() TsInArrays {
	var tscom TsInArrays
	tscom.Name=ts.Name
	for _,v:=range ts.DataSeries{
		//if v.Chron!=0{
			tscom.Chron=append(tscom.Chron,v.Chron)
			tscom.ChronInFloat64 =append(tscom.ChronInFloat64,float64(v.Chron.UnixNano()))
		//}
		if math.IsNaN(v.Meas)==false{
			tscom.Meas=append(tscom.Meas,v.Meas)
		}
		tscom.Dchron=append(tscom.Dchron,v.Dchron)
		tscom.DchronInFloat64=append(tscom.DchronInFloat64,float64(v.Dchron.Nanoseconds()))
		if math.IsNaN(v.Dmeas)==false{
			tscom.Dmeas=append(tscom.Dmeas,v.Dmeas)
		}
	}
	return tscom
}
// CompleteTs Sort in ascending order and complete a raw data series with the delta times Dchron, Dmeas interval and summary statistics
func (ts *TimeSeries) CompleteTs(){
	ts.SortChronAsc()
	ts.ComputeDeltas()
	ts.ComputeStats()
}
// ChronToArr is a Utilities producing a fixed slice of Datetimes from a TimeSeries. Useful for computation requiring arrays. Typically math libraries
func (ts TimeSeries) ChronToArr() []time.Time{
	var measarr []time.Time
	for index,_:=range ts.DataSeries{
		measarr=append(measarr, ts.DataSeries[index].Chron)
	}
	return measarr
}
// MeasToArr Utilities producing a fixed slice of Measures from a TimeSeries. Useful for computation requiring arrays. Typically math libraries
func (ts TimeSeries) MeasToArr() []float64{
	var measarr []float64
	for index,_:=range ts.DataSeries{
		measarr=append(measarr, ts.DataSeries[index].Meas)
	}
	return measarr
}



