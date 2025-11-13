package timeseries

import (
	"reflect"
	"testing"
	"time"
)

func TestAddDataUnit(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)

	tests := []struct {
		name     string
		initial  []DataUnit
		toAdd    []DataUnit
		expected []DataUnit
	}{
		{
			name:     "Add single DataUnit to empty series",
			initial:  []DataUnit{},
			toAdd:    []DataUnit{{Chron: t1, Meas: 42.0}},
			expected: []DataUnit{{Chron: t1, Meas: 42.0}},
		},
		{
			name: "Add multiple DataUnits to existing series",
			initial: []DataUnit{
				{Chron: t1, Meas: 10.0},
			},
			toAdd: []DataUnit{
				{Chron: t2, Meas: 20.0},
				{Chron: t2.Add(time.Hour), Meas: 30.0},
			},
			expected: []DataUnit{
				{Chron: t1, Meas: 10.0},
				{Chron: t2, Meas: 20.0},
				{Chron: t2.Add(time.Hour), Meas: 30.0},
			},
		},
		{
			name:     "No DataUnit added (empty input)",
			initial:  []DataUnit{{Chron: t1, Meas: 5.5}},
			toAdd:    []DataUnit{},
			expected: []DataUnit{{Chron: t1, Meas: 5.5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TimeSeries{DataSeries: append([]DataUnit{}, tt.initial...)}
			ts.AddDataUnit(tt.toAdd...)

			if !reflect.DeepEqual(ts.DataSeries, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, ts.DataSeries)
			}
		})
	}
}
