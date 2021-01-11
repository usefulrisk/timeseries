package timeseries

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"math"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// DataSeries is the slice of data that is the main data series element. The type is defined so that it can be used
// as a composition without defining a variable in the TimeSeries.
type DataSeries []DataUnit
// a TimeSeries is a collection of DataUnit, a composition of a slice of DataUnit.
// A map of strings is used as metadata. It should be, in first approximation, a superset of the metadata of the DataUnit
// in the TimeSeries
type TimeSeries struct {
	Name string
	Comment string
	DataSeries
	TsStats
}
//TsInArrays is a Time Series that is used for external communication. The structure is "transposed"
//in that sense that the data are in one array. It is easier to read with other softwares (d3, R and so on)
//but the data are less easy to work with. First use: charts
type TsInArrays struct{
	Name            string
	Chron           []time.Time
	ChronInFloat64  []float64
	Meas            []float64
	Dchron          []time.Duration
	DchronInFloat64 []float64
	Dmeas           []float64
}
// A Summary Statistics field can be stored at the timeseries level to spare some further computation
type TsStats struct {
	Len     int
	Chmin   time.Time
	Chmax   time.Time
	Chmed   time.Time
	Chmean time.Time
	Chstd   time.Time
	Msmin   float64
	Msmax   float64
	Msmean  float64
	Msmed float64
	Msstd   float64
	DChmin  time.Duration
	DChmax  time.Duration
	DChmean time.Duration
	DChmed time.Duration
	DChstd  time.Duration
	DMsmin  float64
	DMsmax float64
	DMsmed float64
	DMsmean float64
	DMsstd float64
	NbreOfNaN int

}
// SortChronAsc sort a TimeSeries in chronological ascending order in-place
func (ts *TimeSeries)SortChronAsc(){
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Chron.Before(ts.DataSeries[j].Chron)
	})
}
// SortChronDesc sort a TimeSeries in chronological descending order in-place
func (ts *TimeSeries)SortChronDesc(){
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Chron.After(ts.DataSeries[j].Chron)
	})
}
// Sort a TimeSeries in ascending order of measure in-place
func (ts *TimeSeries) SortMeasAsc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Meas< ts.DataSeries[j].Meas
	})
}
// Sort a TimeSeries in descending order of measure in-place
func (ts *TimeSeries) SortMeasDesc() {
	sort.Slice(ts.DataSeries, func(i, j int) bool {
		return ts.DataSeries[i].Meas>ts.DataSeries[j].Meas
	})
}
// ResetTimeSeries a TimeSeries to a zero length TimeSeries but keep metadata information at the TimeSeries level
func (ts *TimeSeries) ResetTimeSeries(){
	ts.DataSeries=ts.DataSeries[0:0]
}

// ComputeDeltas is used in ComputeStats method. No need to use it alone.
func (ts *TimeSeries)ComputeDeltas(){
	ts.SortChronAsc()
	if len(ts.DataSeries)!=0 {
		ts.DataSeries[0].Dchron =0*time.Nanosecond
		ts.DataSeries[0].Dmeas = math.NaN()
		for i := 1; i < len(ts.DataSeries); i++ {
			ts.DataSeries[i].Dchron = ts.DataSeries[i].Chron.Sub(ts.DataSeries[i-1].Chron)
			ts.DataSeries[i].Dmeas = ts.DataSeries[i].Meas - ts.DataSeries[i-1].Meas
		}
	}
}
// Deep copy of a time series
func (ts *TimeSeries) Copy() TimeSeries {
	var copyofts TimeSeries
	copyofts.Name=ts.Name
	for _,value:=range ts.DataSeries{
		copyofts.DataSeries=append(copyofts.DataSeries,value)
	}
	copyofts.ComputeStats()
	return copyofts
}
// Compute a series of usual and basic stats for a time series. Should be used intensively for monitoring
// Note for myself: Converting Unix nanoseconds to time.Time is done with function Unix(sec int64,nsec int64) Time.
// Converting time.Time to Unix nanoseconds is done with Method (t Time) UnixNano() int64 or Method (t Time) Unix() int64 if
// precision is not needed. This is done in the ConvertTimeSeries2TsInArrays function. Reminder: In order to compute stats
// with the montanaflynn package, missing data (math.NaN() ) are not imported in float64 arrays. Vectors have therefore
// different length. It is not correct to assume that Chron.length and Meas.length have the same length.
func (ts *TimeSeries) ComputeStats(){
	if len(ts.DataSeries) >0{
		// Preparation -----------------------------------
		ts.ComputeDeltas()
		arrts:=ts.ConvertTimeSeries2TsInArrays()
		ts.Len= len(arrts.Chron)
		// stats on Chron ------------------------------------------
		ts.Chmin= ts.DataSeries[0].Chron
		ts.Chmax= ts.DataSeries[ts.Len-1].Chron
		for _,v :=range ts.DataSeries{
			if math.IsNaN(v.Meas)==true{
				ts.NbreOfNaN+=1
			}
		}
		med,_:=stats.Median(arrts.ChronInFloat64)
		timenanosec :=int64(med)
		nanosec:= timenanosec %1000000000
		sec:= timenanosec /1000000000
		ts.Chmed =time.Unix(sec,nanosec)
		moy,_:=stats.Mean(arrts.ChronInFloat64)
		timenanosec=int64(moy)
		nanosec=timenanosec%1000000000
		sec= timenanosec /1000000000
		ts.Chmean=time.Unix(0,timenanosec)

		/*std,_:=stats.StandardDeviation(float64time)
		timenanosec=int64(std)
		nanosec=timenanosec%1000000000
		sec= timenanosec /1000000000
		ts.Chstd=time.Unix(0,timenanosec)
		 */
		// stats on Meas --------------------------------------------

		ts.Msmin,_=stats.Min(arrts.Meas)
		ts.Msmax,_=stats.Max(arrts.Meas)
		ts.Msmean,_=stats.Mean(arrts.Meas)
		ts.Msmed,_=stats.Median(arrts.Meas)
		ts.Msstd,_=stats.StandardDeviation(arrts.Meas)
		// stats on Dchron -------------------------------------------
		mindch,_:=stats.Min(arrts.ChronInFloat64)
		maxdch,_:=stats.Max(arrts.ChronInFloat64)
		meandch,_:=stats.Mean(arrts.ChronInFloat64)
		meddch,_:=stats.Median(arrts.ChronInFloat64)
		devdch,_:=stats.StandardDeviation(arrts.ChronInFloat64)
		ts.DChmin=time.Duration(int64(mindch/1000000000))
		ts.DChmax=time.Duration(int64(maxdch/1000000000))
		ts.DChmed=time.Duration(int64(meddch/1000000000))
		ts.DChmean=time.Duration(int64(meandch/1000000000))
		ts.DChstd=time.Duration(int64(devdch/1000000000))
		// stats on DMeas ---------------------------------------------
		ts.DMsmin,_=stats.Min(arrts.Dmeas)
		ts.DMsmax,_=stats.Max(arrts.Dmeas)
		ts.DMsmean,_=stats.Mean(arrts.Dmeas)
		ts.DMsmed,_=stats.Median(arrts.Dmeas)
		ts.DMsstd,_=stats.StandardDeviation(arrts.Dmeas)
		// end of process ---------------------------------
		ts.Comment=" Time Series ok."
	}else {
		ts.Comment="Warning: Empty Time Series"
	}
}

// Merge two TimeSeries. Return a new TimeSeries
func Merge(tsa *TimeSeries, tsb *TimeSeries) TimeSeries {
	var tsr TimeSeries
	for _, element := range tsa.DataSeries {
		tsr.AddDataUnit(element)
	}
	for _, element := range tsb.DataSeries {
		tsr.AddDataUnit(element)
	}
	return tsr
}
// Remove non valid dataunits from a time series based on the state of DataUnit. The output is two time series: the new one and one with rejected data.
func (ts TimeSeries) RemovedNonValid() (TimeSeries, TimeSeries){
	var cleanedts,rejectedts TimeSeries
	for _,du := range ts.DataSeries{
		if du.State==""{
			cleanedts.AddDataUnit(du)
		}else{
			rejectedts.AddDataUnit(du)
		}
	}
	return cleanedts,rejectedts
}
// PrettyPrint prints a TimeSeries in a more ordered way in output terminal
func (ts *TimeSeries) PrettyPrint(what ...int) {
	fmt.Printf("Name       : %v\n",ts.Name)
	fmt.Printf("Time Series Identification: %v\n",ts.Comment)
	fmt.Println("----------------------------------------------------------------------------------------------------------------------")
	var j,k int
	switch len(what) {
	case 0:
		{
			j = 0
			k = len(ts.DataSeries)
		}
	case 1:
		{
			j = 0
			k = what[1]
		}
	case 2:
		{
			j = what[0]
			k = what[1]
		}
	default:
		{
			j = 0
			k = len(ts.DataSeries)
		}
	}
	w:=new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "%v|\t%s|\t%v|\t%v|\t%v|\t%v|\t\n","index","Chron","Measure","Delta Chron","Delta Meas","State")
	fmt.Fprintln(w, "-----\t---------------------------------|\t------------|\t-------------|\t-------------|\t------------|\t")
	for i:=j;i<k && i<len(ts.DataSeries);i++{

		if ts.DataSeries[i].State=="" {
			if i==0{
				fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|\t\n", i,ts.DataSeries[i].Chron.Round(0), ts.DataSeries[i].Meas, " ", "")
			}else {
				fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|\t\n", i, ts.DataSeries[i].Chron.Round(0), ts.DataSeries[i].Meas, ts.DataSeries[i].Dchron, ts.DataSeries[i].Dmeas)
			}
		}else{
			if i==0{
				fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|\t%v|\t\n", i, ts.DataSeries[i].Chron.Round(0), ts.DataSeries[i].Meas, "","", ts.DataSeries[i].State)

			}else{
				fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|%v|\t\t\n", i, ts.DataSeries[i].Chron.Round(0), ts.DataSeries[i].Meas, ts.DataSeries[i].Dchron, ts.DataSeries[i].Dmeas, ts.DataSeries[i].State)
			}
		}

	}
	// Format right-aligned in space-separated columns of minimal width 5
	// and at least one blank of padding (so wider column entries do not
	// touch each other).

	fmt.Fprintln(w)
	w.Flush()
}

// PrintTsStats prints TsStats struct in a readable way on terminal
func (ts *TimeSeries) PrintTsStats(){
	fmt.Println("------------------------------------------")
	fmt.Println(ts.Name)
	fmt.Println("------------------------------------------")
	fmt.Printf("Description: %v\n",ts.Comment)
	fmt.Println("------------------------------------------")
	fmt.Printf("Warning: %v Missing Data\n",ts.NbreOfNaN)
	w:=new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Length| %v|\t\n",ts.Len)
	fmt.Fprintln(w, "\tChron|\tMeasure|\tDChron|\tDMeas|\t")
	fmt.Fprintln(w, "-\t-----------------\t------------\t------------\t------------\t")
	fmt.Fprintf(w, "Min|\t %v|\t%v|\t%v|\t%v|\t\n",ts.Chmin.Round(0),ts.Msmin,ts.DChmin,ts.DMsmin)
	fmt.Fprintf(w, "Max|\t %v|\t%v|\t%v|\t%v|\t\n", ts.Chmax.Round(0),ts.Msmax,ts.DChmax,ts.DMsmax)
	fmt.Fprintf(w, "Mean|\t %v|\t%v|\t%v|\t%v|\t\n",ts.Chmean,ts.Msmean,ts.DChmean,ts.DMsmean)
	fmt.Fprintf(w, "Median|\t %v|\t%v|\t%v|\t%v|\t\n",ts.Chmed,ts.Msmean,ts.DChmed,ts.DMsmed)
	fmt.Fprintf(w, "StdDev|\t %v|\t%v|\t%v|\t%v|\t\n" ," ",ts.Msstd,ts.DChstd,ts.DMsstd)

	fmt.Fprintln(w)
	w.Flush()

}