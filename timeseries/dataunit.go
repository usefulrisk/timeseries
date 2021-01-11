package timeseries

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

// a DataUnit is the smallest unity of communication. The DataUnit is supposed to be transferred through any context,
// internet amongst others. It is supposed to be informative enough to live on its own, a bit like a packet in TCP/IP protocol.
// It is easily transformed into text/json object. More specifically, all the metadata are located in the "Measurement" field.
// A map of strings or a json can be used as a more sophisticated identity device, but is not used here for simplicity. Using a map requires
// specific initialization. Differential fields are created from the start because any time series work needs them sooner or later.
type DataUnit struct {
	Measurement string // Where to classify the Datum. This can be used as a table name in a SQL environment.
	Chron time.Time // time stamp according to golang time
	Meas float64 // measure. In this setting measure is float.
	Dchron  time.Duration // Duration time from preceding Observation
	Dmeas   float64       // Delta measure from preceding observation
	State string // State of measuring device ("ok","nok"....) Could be used as non floating measuring device if any
}

const(
	NotNihil float64 =2.2250738585072014e-308 //NotNihil==math.SmallestNonzeroFloat64
)
// NewDataUnit creates a DataUnit from input. Warning: 0 as a Meas(ure) is significant, notwithstanding the fact that
// golang creates a 0 value field by default. If the data is not 0, use math.NaN() function.
func NewDataUnit(table string,chr time.Time,meas float64,st string)DataUnit{
	var du DataUnit
	du.Measurement=table
	du.Chron=chr
	du.Meas=meas
	du.State=st
	return du
}
// AddDataUnit adds well formed DataUnit (s) to a TimeSeries. Variadic form
func (ts *TimeSeries) AddDataUnit(dus ...DataUnit){
	for _,du :=range dus{
		ts.DataSeries =append(ts.DataSeries,du)
	}
}
// AddData adds DataUnit with data inputs. Empty fields are golang default values
func (ts *TimeSeries) AddData(table string,chr time.Time,meas float64,st string){
	du:=NewDataUnit(table ,chr ,meas ,st)
	ts.DataSeries =append(ts.DataSeries,du)
}
// PrettyPrint prints DataUnit as readable form to Terminal
func (du *DataUnit)PrettyPrint(){
	w:=new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "%s|\t%v|\t%v|\t%v|\t%v|\t%v|\t\n","Measurement","Chron","Measure","State","Dchron","Dmeas")
	fmt.Fprintln(w, "------------|\t--------------------------------------|\t----------|\t------------|\t------------|\t------------|\t")
	fmt.Fprintf(w, "%s|\t%v|\t%v|\t%v|\t%v|\t%v|\t\n",du.Measurement,du.Chron.Round(0), du.Meas,du.State,du.Dchron,du.Dmeas)
	fmt.Fprintln(w)
	w.Flush()
}

// ParseRFC3339toTime parse a RFC3339 date type string to golang datetime type
func ParseRFC3339toTime(nimp string, location *time.Location)time.Time{
	var dd time.Time
	lastc:=nimp[len(nimp)-1:]
	if lastc=="S" || lastc=="s" {
		yy,_:=strconv.Atoi(nimp[:2])
		yy=yy+2000
		mm,_:=strconv.Atoi(nimp[2:4])
		dda,_:=strconv.Atoi(nimp[4:6])
		hh,_:=strconv.Atoi(nimp[6:8])
		min,_:=strconv.Atoi(nimp[8:10])
		ss,_:=strconv.Atoi(nimp[10:12])
		dd=time.Date(yy,time.Month(mm),dda,hh,min,ss,0.0,location)
		return dd
	}
	layout := "2006-01-02T15:04:05.000Z"
	dd, _ = time.Parse(layout, nimp)
	return dd
}


