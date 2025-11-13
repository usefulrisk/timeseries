package timeseries

import (
	"math"
	"time"
)

// toPtrOrNil converts a float64 to *float64, mapping NaN to nil.
// Unexported helper, kept small and focused.
func toPtrOrNil(f float64) *float64 {
	if math.IsNaN(f) {
		return nil
	}
	v := f
	return &v
}

// ToJSON converts a TimeSeries into its JSON-friendly DTO.
// It produces:
//   - Chron      []time.Time
//   - Meas       []*float64 (NaN -> nil)
//   - DchronNS   []int64    (time.Duration in ns)
//   - Dmeas      []*float64 (NaN -> nil)
//   - Status     []StatusCode
func (ts *TimeSeries) ToJSON() *TimeSeriesJSON {
	n := len(ts.DataSeries)
	chron := make([]time.Time, n)
	meas := make([]*float64, n)
	dchron := make([]int64, n)
	dmeas := make([]*float64, n)
	status := make([]StatusCode, n)

	for i, du := range ts.DataSeries {
		chron[i] = du.Chron
		meas[i] = toPtrOrNil(du.Meas)
		dchron[i] = int64(du.Dchron)
		dmeas[i] = toPtrOrNil(du.Dmeas)
		status[i] = du.Status
	}

	return &TimeSeriesJSON{
		Name:     ts.Name,
		Comment:  ts.Comment,
		Chron:    chron,
		Meas:     meas,
		DchronNS: dchron,
		Dmeas:    dmeas,
		Status:   status,
		Stats:    ts.BasicStats.ToJSON(),
	}
}

// ToJSON converts BasicStats into a JSON-friendly DTO.
// Note: keep ns in int64 to avoid float rounding.
func (s *BasicStats) ToJSON() *BasicStatsJSON {
	if s == nil {
		return nil
	}
	return &BasicStatsJSON{
		Len:        s.Len,
		Chmin:      s.Chmin,
		ValAtChmin: s.ValAtChmin,
		Chmax:      s.Chmax,
		ValAtChmax: s.ValAtChmax,
		Chmed:      s.Chmed,
		Chmean:     s.Chmean,
		Chstd:      s.Chstd,
		Msmin:      s.Msmin,
		ChAtMsmin:  s.ChAtMsmin,
		Msmax:      s.Msmax,
		ChAtMsmax:  s.ChAtMsmax,
		Msmean:     s.Msmean,
		Msmed:      s.Msmed,
		Msstd:      s.Msstd,

		DChminNS:   int64(s.DChmin),
		ChAtDChmin: s.ChAtDChmin,
		DChmaxNS:   int64(s.DChmax),
		ChAtDChmax: s.ChAtDchmax,
		DChmeanNS:  int64(s.DChmean),
		DChmedNS:   int64(s.DChmed),
		DChstdNS:   int64(s.DChstd),

		DMsmin:    s.DMsmin,
		DMsmax:    s.DMsmax,
		DMsmed:    s.DMsmed,
		DMsmean:   s.DMsmean,
		DMsstd:    s.DMsstd,
		NbreOfNaN: s.NbreOfNaN,
	}
}

// ToJSON converts a TsContainer into its JSON-friendly DTO and returns it.
// No side effects; the caller decides what faire des bytes (Marshal, log, etc.).
func (tsc *TsContainer) ToJSON() *TsContainerJSON {
	out := &TsContainerJSON{
		Name:    tsc.Name,
		Comment: tsc.Comment,
		Series:  make(map[string]*TimeSeriesJSON, len(tsc.Ts)),
	}
	for k, v := range tsc.Ts {
		out.Series[k] = v.ToJSON()
	}
	return out
}
