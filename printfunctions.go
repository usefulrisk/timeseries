package timeseries

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type Series interface {
	PrettyPrint()
}

func (tsc *TsContainer) PrettyPrint(what ...int) {
	for k, v := range tsc.Ts {
		if v != nil {
			fmt.Printf("Container: %v\n", tsc.Name)
			fmt.Println("-------------------------------------------------")
			fmt.Printf("TimeSeries: %v\n", k)
			v.PrintTsStats()
			v.PrettyPrint(what...)
		}
	}
}

// PrettyPrint prints DataUnit as readable form to Terminal
func (du *DataUnit) PrettyPrint() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "%s|\t%v|\t%v|\t%v|\t%v|\t%v\t\n", "Measurement", "Chron", "Measure", "Dchron", "Dmeas", "Status")
	fmt.Fprintln(w, "------------|\t--------------------------------------|\t----------|\t------------|\t------------|\t------------|\t")
	fmt.Fprintf(w, "%v|\t%v|\t%v|\t%v|\t%v|\t%v|\t\n", du.Meas, du.Chron.Round(0), du.Meas, du.Dchron, du.Dmeas, du.Status)
	fmt.Fprintln(w)
	w.Flush()
}

// PrettyPrint prints a TimeSeries in a more ordered way in output terminal
func (ts *TimeSeries) PrettyPrint(what ...int) {
	fmt.Printf("Name       : %v\n", ts.Name)
	fmt.Printf("Time Series Identification: %v\n", ts.Comment)
	fmt.Println("----------------------------------------------------------------------------------------------------------------------")
	var j, k int
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
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "%v|\t%s|\t%v|\t%v|\t%v|\t%v|\t\n", "index", "Chron", "Measure", "Delta Chron", "Delta Meas", "Satus")
	fmt.Fprintln(w, "-----\t---------------------------------|\t------------|\t-------------|\t-------------|\t-----------------|\t")
	for i := j; i < k && i < len(ts.DataSeries); i++ {
		if i == 0 {
			fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|\t\n", i, ts.DataSeries[i].Chron.Format(time.RFC3339Nano), ts.DataSeries[i].Meas, " ", "")
		} else {
			fmt.Fprintf(w, "%d|\t%s|\t%v|\t%v|\t%v|\t%v|\t\n", i, ts.DataSeries[i].Chron.Format(time.RFC3339Nano), ts.DataSeries[i].Meas, ts.DataSeries[i].Dchron, ts.DataSeries[i].Dmeas, ts.DataSeries[i].Status)
		}
	}
	// Format right-aligned in space-separated columns of minimal width 5
	// and at least one blank of padding (so wider column entries do not
	// touch each other).
	fmt.Fprintln(w)
	w.Flush()
}

// PrintTsStats prints TsStats struct in a readable way in output terminal
func (ts *TimeSeries) PrintTsStats() {
	fmt.Println("------------------------------------------")
	fmt.Println(ts.Name)
	fmt.Printf("Warning: %v Missing Data\n", ts.NbreOfNaN)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Length| %v|\t\n", ts.Len)
	fmt.Fprintln(w, "\tChron|\tMeasure|\tDChron|\tDMeas|\t")
	fmt.Fprintln(w, "-\t-----------------\t------------\t------------\t------------\t")
	fmt.Fprintf(w, "Min|\t %v|\t%v|\t%v|\t%v|\t\n", ts.Chmin.Round(0), ts.Msmin, ts.DChmin, ts.DMsmin)
	fmt.Fprintf(w, "Max|\t %v|\t%v|\t%v|\t%v|\t\n", ts.Chmax.Round(0), ts.Msmax, ts.DChmax, ts.DMsmax)
	fmt.Fprintf(w, "Mean|\t %v|\t%v|\t%v|\t%v|\t\n", ts.Chmean, ts.Msmean, ts.DChmean, ts.DMsmean)
	fmt.Fprintf(w, "Median|\t %v|\t%v|\t%v|\t%v|\t\n", ts.Chmed, ts.Msmean, ts.DChmed, ts.DMsmed)
	fmt.Fprintf(w, "StdDev|\t %v|\t%v|\t%v|\t%v|\t\n", " ", ts.Msstd, ts.DChstd, ts.DMsstd)

	fmt.Fprintln(w)
	w.Flush()

}
func (ts *TimeSeries) PrettyPrintAll() {
	ts.PrettyPrint()
	ts.PrintTsStats()
}
