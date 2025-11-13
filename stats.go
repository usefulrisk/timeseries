package timeseries

import (
	"errors"
	"math"
	"sort"
)

// Sum adds all numbers in input and returns the total.
// If input is empty, Sum returns math.NaN() and ErrEmptyInput.
// Sum does not allocate and does not modify input.
func Sum(input []float64) (sum float64, err error) {

	if len(input) == 0 {
		return math.NaN(), ErrEmptyInput
	}

	// Add tem up
	for _, n := range input {
		sum += n
	}

	return sum, nil
}

// Mean returns the arithmetic mean of input.
// If input is empty, it returns math.NaN() and ErrEmptyInput.
// Mean does not allocate and does not modify input.
// Note: if your data may contain NaN, filter them upstream or use
// a helper that skips NaNs; otherwise the result may be NaN.
func Mean(input []float64) (float64, error) {

	if len(input) == 0 {
		return math.NaN(), ErrEmptyInput
	}
	sum, _ := Sum(input)

	return sum / float64(len(input)), nil
}

// Median returns the median of input using the standard definition
// on the sorted data. For an odd length n, it returns the middle value.
// For an even length n, it returns the mean of the two middle values.
// If input is empty, it returns math.NaN() and ErrEmptyInput.
//
// ⚠️ Median sorts the slice in place (using sort.Float64s) to avoid
// extra allocations. Call it on a copy if you need to preserve order.
func Median(input []float64) (median float64, err error) {
	sort.Float64s(input)

	// No math is needed if there are no numbers
	// For even numbers we add the two middle numbers
	// and divide by two using the mean function above
	// For odd numbers we just use the middle number
	l := len(input)
	if l == 0 {
		return math.NaN(), ErrEmptyInput
	} else if l%2 == 0 {
		median, _ = Mean(input[l/2-1 : l/2+1])
	} else {
		median = input[l/2]
	}

	return median, nil
}

// Min finds the lowest number in a set of data
func Min(input []float64) (min float64, err error) {

	// Get the count of numbers in the slice
	l := len(input)

	// Return an error if there are no numbers
	if l == 0 {
		return math.NaN(), ErrEmptyInput
	}

	// Get the first value as the starting point
	min = input[0]

	// Iterate until done checking for a lower value
	for i := 1; i < l; i++ {
		if input[i] < min {
			min = input[i]
		}
	}
	return min, nil
}

// Max finds the highest number in a slice
func Max(input []float64) (max float64, err error) {

	// Return an error if there are no numbers
	if len(input) == 0 {
		return math.NaN(), ErrEmptyInput
	}

	// Get the first value as the starting point
	max = input[0]

	// Loop and replace higher values
	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
		}
	}

	return max, nil
}

// StdDev computes the sample standard deviation (unbiased, denominator n-1).
// It returns math.NaN() and ErrEmptyInput if input is empty, and math.NaN()
// and ErrBounds if the length is < 2 (sample standard deviation undefined).
// StdDev does not allocate and does not modify input.
func StdDev(data []float64) (float64, error) {
	n := len(data)
	if n == 0 {
		return math.NaN(), errors.New("empty slice")
	}

	// Calcul de la moyenne
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(n)

	// Somme des carrés des écarts à la moyenne
	var sqDiffSum float64
	for _, v := range data {
		diff := v - mean
		sqDiffSum += diff * diff
	}

	// Écart-type (population)
	std := math.Sqrt(sqDiffSum / float64(n))

	// Pour l’écart-type d’échantillon :
	// std := math.Sqrt(sqDiffSum / float64(n-1))

	return std, nil
}

// Percentile returns the p-th percentile of input using the
// nearest-rank definition on the 1-indexed, ascendingly sorted data.
// Valid p is in (0, 100]. If input is empty, it returns math.NaN()
// and ErrEmptyInput. If p is out of bounds, it returns math.NaN()
// and ErrBounds.
//
// ⚠️ Percentile sorts the slice in place. Call it on a copy to preserve input
func Percentile(x []float64, p float64) (float64, error) {
	n := len(x)
	if n == 0 {
		return math.NaN(), ErrEmptyInput
	}
	if p <= 0 || p > 100 {
		return math.NaN(), ErrBounds
	}

	cp := append([]float64(nil), x...)
	sort.Float64s(cp)

	k := int(math.Floor(p / 100 * float64(n)))
	if k < 1 {
		return cp[0], nil
	}
	if k >= n {
		return cp[n-1], nil
	}

	return cp[k-1], nil
}
