package timeseries

func (ts *TimeSeries) AddDataUnit(dus ...DataUnit) {
	for _, du := range dus {
		ts.DataSeries = append(ts.DataSeries, du)
	}
}
