package timeseries

import (
	"fmt"
	"time"
)

// TsContainer is a container for TimeSeries. Container is used to keep track of different versions of a TimeSeries.
type TsContainer struct {
	Name string
	Ts   map[string] *TimeSeries
}

// NewTsContainer creates new TsContainer
func NewTsContainer(name string)TsContainer{
	var tsc TsContainer
	tsc.Name=name
	tsc.Ts =make(map[string]*TimeSeries)
	return tsc
}

// PrettyPrintContainer prints TsContainer in readable form on terminal
func (tsc *TsContainer)PrettyPrintContainer(what...int){
	for k,v:=range tsc.Ts {
		fmt.Printf("Container: %v\n",tsc.Name)
		fmt.Println("-------------------------------------------------")
		fmt.Printf("TimeSeries: %v\n",k)
		v.PrintTsStats()
		v.PrettyPrint(what...)
	}
}
//ResetContainer resest TimeSeries in Container without changing the name
func (tsc *TsContainer)ResetContainer(){
	for _,v :=range tsc.Ts{
		v.ResetTimeSeries()
	}
}
//PrepareForComm inverses time series in a time series container in order to make it easily graphable by
// javascript library.
func (tsc *TsContainer) PrepareForComm()CommContainer{
	commContainer := NewTsCommContainer(tsc.Name)
	for _,v:=range tsc.Ts{
		commContainer.CommTs=append(commContainer.CommTs,v.ConvertTimeSeries2TsForComm())
	}
	return commContainer
}
// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// ancillary to CommContainer
type CommTs []TsForComm

// CommContainer is a Container with TsInArrays to be sent through the internets.
type CommContainer struct {
	Name string
	CommTs
}

//TsForComm is a Time Series that is used to be sent through the internets. The structure is "inversed"
//in that sense that the data are in one array. It is easier to read with other softwares (d3, R and so on)
//but the data are less easy to work with. First use: to be used to chart.
//
// There is a difference with the array version of TimeSeries, DChron is a string instead of a time.Duration
type TsForComm struct{
	Name string
	Comment string
	TsTags map[string]string
	Chron []time.Time
	Meas []float64
	Dchron []string
	Dmeas []float64
	SummaryStats SummaryStatForExport
}

type SummaryStatForExport struct {
	Len    int
	Chmin  time.Time
	Chmax  time.Time
	Chmean time.Time
	Chstd  time.Time
	Msmin  float64
	Msmax  float64
	Msmean float64
	Msstd  float64
	DChmin string
	DChmax string
	DChmean string
	DChstd string
	DMsmin float64
	DMsmax float64
	DMsmean float64
	DMsstd float64

}

func NewTsCommContainer(name string)CommContainer{
	var cct CommContainer
	return cct
}
func (ts *TimeSeries)ConvertTimeSeries2TsForComm()TsForComm{
	var tscom TsForComm
	tscom.Comment=ts.Comment
	for _,v:=range ts.DataSeries{
		timein,_:=TimeIn(v.Chron,"Local")
		tscom.Chron=append(tscom.Chron,timein)
		tscom.Meas=append(tscom.Meas,v.Meas)
		tscom.Dchron=append(tscom.Dchron,v.Dchron.String())
		tscom.Dmeas=append(tscom.Dmeas,v.Dmeas)
	}
	tscom.SummaryStats.Msmin=ts.Msmin
	tscom.SummaryStats.Len=ts.Len
	tscom.SummaryStats.Chmin,_=TimeIn(ts.Chmin,"Local")
	tscom.SummaryStats.Chmax,_=TimeIn(ts.Chmax,"Local")
	tscom.SummaryStats.Chmean,_=TimeIn(ts.Chmean,"Local")
	tscom.SummaryStats.Chstd=ts.Chstd
	tscom.SummaryStats.Msmin=ts.Msmin
	tscom.SummaryStats.Msmax=ts.Msmax
	tscom.SummaryStats.Msmean=ts.Msmean
	tscom.SummaryStats.Msstd=ts.Msstd
	tscom.SummaryStats.DChmin=ts.DChmin.String()
	tscom.SummaryStats.DChmax=ts.DChmax.String()
	tscom.SummaryStats.DChmean=ts.DChmean.String()
	tscom.SummaryStats.DChstd=ts.DChstd.String()
	tscom.SummaryStats.DMsmin=ts.DMsmin
	tscom.SummaryStats.DMsmax=ts.DMsmax
	tscom.SummaryStats.DMsmean=ts.DMsmean
	tscom.SummaryStats.DMsstd=ts.DMsstd
	return tscom
}