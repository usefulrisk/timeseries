# timeseries

A Go library for **building, cleaning, and analyzing time series** with explicit handling of **missing data**, **outliers**, and **irregular intervals**.

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/timeseries.svg)](https://pkg.go.dev/github.com/yourusername/timeseries)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 🧭 Overview

`timeseries` provides lightweight yet powerful tools to work with temporal data where gaps and anomalies are common — such as environmental, energy, or IoT sensor datasets.

Unlike most numeric libraries, **missing and invalid data are treated as first-class citizens**.  
Every statistical and transformation function in `timeseries` safely handles incomplete or flagged entries.

The library focuses on **clarity, composability, and correctness**, making it ideal for scientific and engineering workflows.

---

## ✨ Key Features

- 🧩 **Explicit status codes** (`OK`, `Missing`, `Outlier`, `Invalid`) for each data point
- 🕒 **Regularization** — resample or rebuild irregular time series at a fixed interval
- 🚫 **Robust outlier detection** using IQR, z-score, rolling deviation, or custom filters
- 📊 **Statistical summaries** (mean, median, stddev, percentiles) that ignore missing data
- 🔄 **Simple data model**:
  ```go
  type DataUnit  // a single timestamped observation
  type TimeSeries // an ordered collection of DataUnits