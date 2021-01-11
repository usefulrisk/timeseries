# timeseries
A library for managing time series
Most of statistical software (R, pandas,...) require the data to be evenly distanced in time. This is a fiction when working in IoT or when one processes data from source in a more general way. When receiving data from sources, before processing data in a statistics oriented software, one must:

1. organize reception of data (from source to database. Not in this library)
2. structure them into time series
3. clean data (missing data, impossible data, out-of-bounds,...)
4. regularize time frequency (seconds, minutes, hours,...)
5. keep track of the processes applied to original data (easily identify)

This requires mastering of time library and missing data. Mastering time and time statistics is especially tricky. So much that I could not find any library for our purpose.
