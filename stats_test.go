package timeseries

import (
	"math"
	"reflect"
	"testing"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSum(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		want    float64
		wantErr bool
	}{
		{"normal", []float64{1, 2, 3, 4}, 10, false},
		{"empty", []float64{}, math.NaN(), true},
	}

	for _, tt := range tests {
		got, err := Sum(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("%s: got %.4f, want %.4f", tt.name, got, tt.want)
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		want    float64
		wantErr bool
	}{
		{"normal", []float64{1, 2, 3, 4}, 2.5, false},
		{"empty", []float64{}, math.NaN(), true},
	}

	for _, tt := range tests {
		got, err := Mean(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("%s: got %.4f, want %.4f", tt.name, got, tt.want)
		}
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		want    float64
		wantErr bool
	}{
		{"odd", []float64{1, 3, 2}, 2, false},
		{"even", []float64{1, 2, 3, 4}, 2.5, false},
		{"empty", []float64{}, math.NaN(), true},
	}

	for _, tt := range tests {
		got, err := Median(append([]float64(nil), tt.input...))
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("%s: got %.4f, want %.4f", tt.name, got, tt.want)
		}
	}
}

func TestMinMax(t *testing.T) {
	data := []float64{5, 2, 9, 1, 7}

	min, err := Min(data)
	if err != nil || min != 1 {
		t.Errorf("Min() = %.2f, err=%v, want 1", min, err)
	}

	max, err := Max(data)
	if err != nil || max != 9 {
		t.Errorf("Max() = %.2f, err=%v, want 9", max, err)
	}

	_, err = Min([]float64{})
	if err == nil {
		t.Errorf("Min() expected error for empty slice")
	}

	_, err = Max([]float64{})
	if err == nil {
		t.Errorf("Max() expected error for empty slice")
	}
}

func TestStdDev(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		want    float64
		wantErr bool
	}{
		{"simple", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 2, false}, // exemple classique
		{"empty", []float64{}, math.NaN(), true},
	}

	for _, tt := range tests {
		got, err := StdDev(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("%s: got %.4f, want %.4f", tt.name, got, tt.want)
		}
	}
}

func TestPercentile(t *testing.T) {
	data := []float64{50, 20, 40, 35, 15}

	tests := []struct {
		name    string
		percent float64
		want    float64
		wantErr bool
	}{
		{"25th", 25, 15, false},
		{"50th", 50, 20, false},
		{"75th", 75, 35, false},
		{"0%", 0, math.NaN(), true},
		{">100%", 110, math.NaN(), true},
		{"empty", 50, math.NaN(), true},
	}

	for _, tt := range tests {
		input := append([]float64(nil), data...) // copy
		if tt.name == "empty" {
			input = []float64{}
		}
		got, err := Percentile(input, tt.percent)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && !almostEqual(got, tt.want, 1e-9) {
			t.Errorf("%s: got %.4f, want %.4f", tt.name, got, tt.want)
		}
	}
}

// Vérification d'intégrité sur les erreurs globales
func TestErrorVariables(t *testing.T) {
	if reflect.ValueOf(ErrEmptyInput).IsZero() {
		t.Error("ErrEmptyInput should not be nil")
	}
	if reflect.ValueOf(ErrBounds).IsZero() {
		t.Error("ErrBounds should not be nil")
	}
}
