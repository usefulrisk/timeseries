package timeseries

import (
	"time"
)

// StatusCode classifies the validity of a single observation (a DataUnit).
// It allows algorithms to treat missing and anomalous data explicitly
// instead of silently propagating NaNs or zero-values.
//
// Typical semantics:
//   - StOK:       the observation is valid and should be used in stats.
//   - StMissing:  the observation is missing (e.g., Meas is NaN or absent).
//   - StOutlier:  the observation was flagged as an outlier by a detector.
//   - StInvalid:  the observation exists but must not be used (bad sensor,
//     parse error, unit mismatch, etc.).
type StatusCode uint8

// StOK, StMissing, StOutlier and StInvalid enumerate the canonical states
// an observation can be in. Consumers should prefer StatusCode over ad-hoc
// sentinels (like NaN-only) because it’s explicit and type-safe.
const (
	StOK      StatusCode = iota // default: valid
	StMissing                   // missing value (gap)
	StOutlier                   // flagged outlier
	StInvalid                   // present but unusable
)

// DataUnit represents a single timestamped measurement and its meta-state.
//
// Fields (typical usage):
//   - Chron:  the wall-clock timestamp of the observation.
//   - Meas:   the measured value. If Meas is not a float (or if NaN is used
//     to encode missingness), the Status should carry the ground truth.
//   - Dchron: optional time delta (e.g., since previous Chron) in nanoseconds.
//   - Dmeas:  optional delta of measurement (vs previous point).
//   - Status: the StatusCode describing validity.
//
// DataUnit is kept small and contiguous to allow efficient slices (no pointers
// for the hot path). Missing values can be represented by Status=StMissing and
// optionally Meas=math.NaN().
type DataUnit struct {
	Chron  time.Time
	Meas   float64
	Dchron time.Duration
	Dmeas  float64
	Status StatusCode
}

// TimeSeries is an ordered collection of DataUnit, typically sorted by Chron.
// The type is designed for correctness and traceability: each point carries
// its own StatusCode so that statistics and transforms can honor missing data
// and outliers without ad-hoc filtering.
//
// TimeSeries interoperates well with JSON/CSV/DB sources and can be converted
// to DTOs (see TimeSeriesJSON) for transport or storage.
type TimeSeries struct {
	Name       string
	Comment    string
	DataSeries []DataUnit
	BasicStats
}

type DeltaTimeSeries struct {
	Name string
}

// BasicStats holds summary statistics for a series. All values must be defined
// in terms of valid observations only (Status=StOK), i.e., missing or invalid
// points are excluded from the aggregates.
//
// Field meaning (typical):
//   - Len:        count of valid observations used in stats (not total length).
//   - Chmin/Chmax: timestamps at which the min/max measurement occurred.
//   - ValAtChmin/ValAtChmax: the min/max measurement values.
//   - Chmed/Chmean/Chstd: representative timestamps (median/mean/“std anchor”).
//   - Msmin/Msmax/Msmean/Msmed/Msstd: scalar stats for measurement values.
//   - DChmin/DChmax/…: if present, summaries for time deltas (Dchron).
type BasicStats struct {
	Len        int
	Chmin      time.Time
	ValAtChmin float64
	Chmax      time.Time
	ValAtChmax float64
	Chmed      time.Time
	Chmean     time.Time
	Chstd      time.Time
	Msmin      float64
	ChAtMsmin  time.Time
	Msmax      float64
	ChAtMsmax  time.Time
	Msmean     float64
	Msmed      float64
	Msstd      float64
	DChmin     time.Duration
	ChAtDChmin time.Time
	DChmax     time.Duration
	ChAtDchmax time.Time
	DChmean    time.Duration
	DChmed     time.Duration
	DChstd     time.Duration
	DMsmin     float64
	DMsmax     float64
	DMsmed     float64
	DMsmean    float64
	DMsstd     float64
	NbreOfNaN  int
}
type TsContainer struct {
	Name    string
	Comment string
	Ts      map[string]*TimeSeries
}

// TimeSeriesJSON is a JSON-friendly DTO for TimeSeries. It expands the series
// into columnar arrays so that clients can efficiently ingest or visualize it.
//
// Conventions:
//   - Chron:   []time.Time
//   - Meas:    []*float64 (nil encodes missing instead of NaN)
//   - Dchron:  []int64 (time.Duration in ns)
//   - Dmeas:   []*float64 (nil for missing)
//   - Status:  []StatusCode
//   - Stats:   *BasicStatsJSON
//
// The DTO keeps missing-data semantics explicit for safer cross-language usage.
type TimeSeriesJSON struct {
	Name     string          `json:"name"`
	Comment  string          `json:"comment,omitempty"`
	Chron    []time.Time     `json:"chron"` // RFC3339
	Meas     []*float64      `json:"meas"`  // nil => null
	DchronNS []int64         `json:"dchron_ns,omitempty"`
	Dmeas    []*float64      `json:"dmeas,omitempty"`
	Status   []StatusCode    `json:"status,omitempty"`
	Stats    *BasicStatsJSON `json:"stats,omitempty"`
}

// BasicStatsJSON is the serialized counterpart to BasicStats. It mirrors the
// same semantics but uses JSON-friendly field types (time.Time, float64, int).
type BasicStatsJSON struct {
	Len        int       `json:"len"`
	Chmin      time.Time `json:"chmin"`
	ValAtChmin float64   `json:"valAtChmin"`
	Chmax      time.Time `json:"chmax"`
	ValAtChmax float64   `json:"valAtChmax"`
	Chmed      time.Time `json:"chmed"`
	Chmean     time.Time `json:"chmean"`
	Chstd      time.Time `json:"chstd"`
	Msmin      float64   `json:"msmin"`
	ChAtMsmin  time.Time `json:"chAtMsmin"`
	Msmax      float64   `json:"msmax"`
	ChAtMsmax  time.Time `json:"chAtMsmax"`
	Msmean     float64   `json:"msmean"`
	Msmed      float64   `json:"msmed"`
	Msstd      float64   `json:"msstd"`
	DChminNS   int64     `json:"dChmin_ns"`
	ChAtDChmin time.Time `json:"chAtDChmin"`
	DChmaxNS   int64     `json:"dChmax_ns"`
	ChAtDChmax time.Time `json:"chAtDChmax"`
	DChmeanNS  int64     `json:"dChmean_ns"`
	DChmedNS   int64     `json:"dChmed_ns"`
	DChstdNS   int64     `json:"dChstd_ns"`
	DMsmin     float64   `json:"dMsmin"`
	DMsmax     float64   `json:"dMsmax"`
	DMsmed     float64   `json:"dMsmed"`
	DMsmean    float64   `json:"dMsmean"`
	DMsstd     float64   `json:"dMsstd"`
	NbreOfNaN  int       `json:"nbreOfNaN"`
}

type TsContainerJSON struct {
	Name    string                     `json:"name"`
	Comment string                     `json:"comment,omitempty"`
	Series  map[string]*TimeSeriesJSON `json:"series"`
}

// NewDataUnit constructs a DataUnit from a timestamp and a value,
// assigning the provided StatusCode. It is a convenience helper that
// keeps construction explicit and consistent across the codebase.
// Warning:
//
//	A value of 0.0 for Meas is significant — it will not be treated as "unset".
//	This is important because Go initializes float fields to 0 by default.
//	If your data is missing or undefined, use math.NaN() explicitly.
func NewDataUnit(chr time.Time, meas float64) DataUnit {
	return DataUnit{
		Chron: chr,
		Meas:  meas,
	}
}

// NewDataUnitWithStatus creates a DataUnit with an explicit StatusCode.
func NewDataUnitWithStatus(chr time.Time, meas float64, status StatusCode) DataUnit {
	return DataUnit{
		Chron:  chr,
		Meas:   meas,
		Status: status,
	}
}

// Agg enumerates aggregation policies used when condensing multiple samples
// into a regular time bucket (e.g., during resampling/regularization).
//
// Semantics:
//   - AggMin:  choose the minimum value in the bucket.
//   - AggMax:  choose the maximum value in the bucket.
//   - AggMean: use the arithmetic mean of the bucket.
//   - AggLast: take the last (rightmost) sample in the bucket.
//   - AggSum:  sum all samples in the bucket (useful for counters/energy).
type Agg int

// Aggregation modes for resampling/regularization. Choose the one that
// matches your signal semantics (e.g., AggSum for additive counters,
// AggMean for average power, AggLast for step-like processes).
const (
	AggMin Agg = iota
	AggMax
	AggMean
	AggLast
	AggSum
)
