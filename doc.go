// Package timeseries offers robust tools for building, cleaning, and analyzing
// time series with explicit handling of missing data, outliers, and irregular intervals.
//
// Package timeseries provides lightweight yet powerful tools to build,
// analyze, and clean time series data with explicit handling of missing
// values and irregular sampling.
//
// Unlike most numeric libraries, timeseries treats missing data as a first-class
// citizen: all statistical and transformation functions are designed to work
// safely even when parts of the series are incomplete or flagged as invalid.
// This makes it suitable for environmental, energy, sensor and IoT datasets where
// gaps and outliers are common.
//
// Key features:
//
//   - Explicit status codes (OK, Missing, Outlier, Invalid) attached to each
//     data unit — enabling traceable and reliable data processing pipelines.
//
//   - Built-in resampling and regularization: transform an irregular sequence
//     of measurements into a regularly spaced time series, preserving timing
//     accuracy and filling gaps in a controlled way.
//
//   - Robust outlier detection and cleaning using several statistical methods
//     (IQR, z-score, rolling deviation, or custom filters).
//
//   - Statistical summaries (mean, median, standard deviation, percentiles, etc.)
//     that automatically ignore missing or invalid entries.
//
//   - Transparent and extensible data model:
//     type DataUnit  // a single timestamped observation
//     type TimeSeries // an ordered collection of DataUnits
//     designed to integrate easily with JSON, CSV, or database sources.
//
// The library focuses on correctness, clarity, and composability rather than
// on bulk numerical speed. It complements general-purpose math packages by
// adding robust temporal semantics and missing-data awareness — features
// rarely found in standard Go statistical libraries.
//
// Typical usage:
//
//	ts := timeseries.BulkSimul("demo", start, time.Hour, 24, 20.0, 3.0, time.Minute)
//	ts = ts.Regularize(time.Hour)
//	ts.CleanOutliers(timeseries.MethodIQR)
//	mean, _ := ts.Mean()
package timeseries
