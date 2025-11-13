package timeseries

import (
	"sort"
)

func SortChronoAsc(s TimeSeries, name string) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...), // copie indépendante
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Chron.Before(sorted.DataSeries[j].Chron)
	})
	sorted.Name = name
	return sorted
}
func SortChronoDesc(s TimeSeries, name string) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...), // copie indépendante
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Chron.After(sorted.DataSeries[j].Chron)
	})
	sorted.Name = name
	return sorted
}
func SortMeasAsc(s TimeSeries, name string) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...), // copie indépendante
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Chron.Before(sorted.DataSeries[j].Chron)
	})
	sorted.Name = name
	return sorted
}
func SortMeasDesc(s TimeSeries, name string) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...), // copie indépendante
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Chron.After(sorted.DataSeries[j].Chron)
	})
	sorted.Name = name
	return sorted
}

// SortedMeasAsc retourne une nouvelle Series triée par Meas croissant
func SortedMeasAsc(s TimeSeries) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...), // copie indépendante
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Meas < sorted.DataSeries[j].Meas
	})
	return sorted
}

// SortedMeasDesc retourne une nouvelle Series triée par Meas décroissant
func SortedMeasDesc(s TimeSeries) TimeSeries {
	sorted := TimeSeries{
		Name:       s.Name,
		DataSeries: append([]DataUnit(nil), s.DataSeries...),
	}
	sort.Slice(sorted.DataSeries, func(i, j int) bool {
		return sorted.DataSeries[i].Meas > sorted.DataSeries[j].Meas
	})
	return sorted
}
