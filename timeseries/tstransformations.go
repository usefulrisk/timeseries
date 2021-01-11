package timeseries

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"time"
)

//RemoveOutbounds takes a TimeSeries and produces two new TimeSeries:One is the TimeSeries cleaned from observations outside
//the parameters min and max, the other is the data that have been rejected with "rejected" in the comment
func(tsin *TimeSeries) RemoveOutbounds(min,max float64,comment string)(TimeSeries, TimeSeries){
	var tsout,tsrej TimeSeries
	tsin.SortMeasAsc()
	for _,dataunit:=range tsin.DataSeries{
		if dataunit.Meas>max||dataunit.Meas<min{
			dataunit.State="Rejected because outside("+fmt.Sprintf("%f",min)+","+fmt.Sprintf("%f",max)+") limits"
			tsrej.AddDataUnit(dataunit)
		}else{
			dataunit.State=comment
			tsout.AddDataUnit(dataunit)
		}
	}
	//Attention on a travaillé sur la série originale il ne faut pas oublier de la re-ranger.
	tsin.SortChronAsc()
	tsout.SortChronAsc()
	tsrej.SortChronAsc()
	return tsout,tsrej
}
// Clean following indicated percentile
func (tsin *TimeSeries) PercCleaning(perc float64)(TimeSeries, TimeSeries) {
	tsin.SortMeasAsc()
	dt:= tsin.MeasToArr()
	//min:=stat.Quantile(perc,1,dt,nil)
	min,_:=stats.Percentile(dt,perc)
	//max:=stat.Quantile((1-perc),1,dt,nil)
	max,_:=stats.Percentile(dt,(100-perc))
	return tsin.RemoveOutbounds(min,max,"Percentile ("+fmt.Sprintf("%.2f",min)+","+fmt.Sprintf("%.2f",max)+") - "+time.Now().Round(0).String())
}
// Clean the cleaned timeseries1 following zscore technique. Save leftover data into leftover
func (tsin *TimeSeries) ZscoreCleaning(lvl float64)(TimeSeries, TimeSeries){
	tsin.SortMeasAsc()
	dt:= tsin.MeasToArr()
	dtmean,_:=stats.Mean(dt)
	dtstd,_:=stats.StandardDeviation(dt)
	min:=(dtmean-dtstd)*lvl
	max:=(dtmean+dtstd)*lvl
	return tsin.RemoveOutbounds(min,max,"zScore at "+fmt.Sprintf("%.2f",lvl)+ "("+fmt.Sprintf("%.2f",min)+","+fmt.Sprintf("%.2f",max)+") - "+time.Now().String())
}
