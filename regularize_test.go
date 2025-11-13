package timeseries

import (
	"math"
	"testing"
	"time"
)

func mustTime(y, m, d, hh, mm, ss int) time.Time {
	return time.Date(y, time.Month(m), d, hh, mm, ss, 0, time.UTC)
}

func du(t time.Time, v float64) DataUnit { return DataUnit{Chron: t, Meas: v} }

// compare only length, timestamps, and values (NaN-safe)
func requireSeriesEq(t *testing.T, got, want TimeSeries, eps float64) {
	t.Helper()
	if len(got.DataSeries) != len(want.DataSeries) {
		t.Fatalf("length mismatch: got %d, want %d", len(got.DataSeries), len(want.DataSeries))
	}
	for i := range want.DataSeries {
		if !got.DataSeries[i].Chron.Equal(want.DataSeries[i].Chron) {
			t.Fatalf("Chron[%d] mismatch: got %v, want %v", i, got.DataSeries[i].Chron, want.DataSeries[i].Chron)
		}
		if !almostEq(got.DataSeries[i].Meas, want.DataSeries[i].Meas, eps) {
			t.Fatalf("Meas[%d] mismatch: got %v, want %v", i, got.DataSeries[i].Meas, want.DataSeries[i].Meas)
		}
	}
}

// --------- Bounds ---------

func TestBounds(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		min, max := Bounds(nil)
		if min != 0 || max != 0 {
			t.Fatalf("empty: got (%v,%v), want (0,0)", min, max)
		}
	})
	t.Run("single element", func(t *testing.T) {
		min, max := Bounds([]float64{42})
		if min != 42 || max != 42 {
			t.Fatalf("single: got (%v,%v), want (42,42)", min, max)
		}
	})
	t.Run("mixed values", func(t *testing.T) {
		min, max := Bounds([]float64{3, -2, 7, 7, 0})
		if min != -2 || max != 7 {
			t.Fatalf("mixed: got (%v,%v), want (-2,7)", min, max)
		}
	})
}

// --------- RoundedStartTime ---------

func TestRoundedStartTime(t *testing.T) {
	base := mustTime(2025, 11, 10, 10, 23, 45) // UTC
	if got := RoundedStartTime(base, 15, "m"); !got.Equal(mustTime(2025, 11, 10, 10, 15, 0)) {
		t.Fatalf("minutes: got %v", got)
	}
	if got := RoundedStartTime(base, 30, "s"); !got.Equal(mustTime(2025, 11, 10, 10, 23, 30)) {
		t.Fatalf("seconds: got %v", got)
	}
	if got := RoundedStartTime(base, 2, "h"); !got.Equal(mustTime(2025, 11, 10, 10, 0, 0)) {
		t.Fatalf("hours: got %v", got)
	}
	// "d" path: subtract afreqq days
	if got := RoundedStartTime(base, 1, "d"); !got.Equal(mustTime(2025, 11, 9, 10, 23, 45)) {
		t.Fatalf("days: got %v", got)
	}
}

// --------- AddDuration & AddDurationTol (tol=0 path) ---------

func TestAddDuration(t *testing.T) {
	start := mustTime(2025, 11, 10, 12, 0, 0)
	if got := AddDuration(start, 45, "s"); !got.Equal(mustTime(2025, 11, 10, 12, 0, 45)) {
		t.Fatalf("seconds: got %v", got)
	}
	if got := AddDuration(start, 5, "m"); !got.Equal(mustTime(2025, 11, 10, 12, 5, 0)) {
		t.Fatalf("minutes: got %v", got)
	}
	if got := AddDuration(start, -2, "h"); !got.Equal(mustTime(2025, 11, 10, 10, 0, 0)) {
		t.Fatalf("hours: got %v", got)
	}
}

func TestAddDurationTolZero(t *testing.T) {
	start := mustTime(2025, 11, 10, 12, 0, 0)
	// tol=0 -> same as AddDuration
	if got := AddDurationTol(start, 30, "s", 0); !got.Equal(AddDuration(start, 30, "s")) {
		t.Fatalf("sec tol=0 mismatch: got %v", got)
	}
	if got := AddDurationTol(start, 10, "m", 0); !got.Equal(AddDuration(start, 10, "m")) {
		t.Fatalf("min tol=0 mismatch: got %v", got)
	}
	if got := AddDurationTol(start, 3, "h", 0); !got.Equal(AddDuration(start, 3, "h")) {
		t.Fatalf("hour tol=0 mismatch: got %v", got)
	}
}

// --------- HourlyAvg ---------

func TestHourlyAvg(t *testing.T) {
	// 2 mesures à 10h, 1 à 11h, pas d'autres heures
	ts := TimeSeries{}
	ts.AddDataUnit(
		du(mustTime(2025, 11, 10, 10, 5, 0), 2),
		du(mustTime(2025, 11, 10, 10, 30, 0), 4),
		du(mustTime(2025, 11, 10, 11, 0, 0), 9),
	)
	avg := ts.HourlyAvg()
	// 10h -> (2+4)/2 = 3
	if !almostEq(avg[10], 3, 1e-12) {
		t.Fatalf("10h avg: got %v, want 3", avg[10])
	}
	// 11h -> 9/1 = 9
	if !almostEq(avg[11], 9, 1e-12) {
		t.Fatalf("11h avg: got %v, want 9", avg[11])
	}
	// heures sans données -> division 0/0 => +Inf/NaN possible si impl change,
	// mais ici hrtemp=0, lentemp=0 => hr[i]=NaN. On vérifie NaN pour une heure vide.
	if !math.IsNaN(avg[9]) {
		t.Fatalf("09h expected NaN, got %v", avg[9])
	}
}

// --------- Regularize (tolerance ignored -> set 0) ---------

// Build a simple series:
// base = 10:00:00
// Points: 10:00:05 -> 1, 10:00:20 -> 3, 10:00:45 -> 10
// freq = 30s, windows end at 10:00:30 and 10:01:00
func buildSimpleSeries() TimeSeries {
	base := mustTime(2025, 11, 10, 10, 0, 0)
	ts := TimeSeries{}
	ts.AddDataUnit(
		du(base.Add(5*time.Second), 1),
		du(base.Add(20*time.Second), 3),
		du(base.Add(45*time.Second), 10),
	)
	// Déjà triée par Chron pour ne pas dépendre de SortChronAsc()
	return ts
}

func TestRegularize_Average(t *testing.T) {
	ts := buildSimpleSeries()
	got := ts.Regularize(30, "s", "avg", 0) // tolérance ignorée
	base := mustTime(2025, 11, 10, 10, 0, 0)
	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 2),  // (1+3)/2
		du(base.Add(60*time.Second), 10), // (10)/1
	)
	requireSeriesEq(t, got, want, 1e-12)
}

func TestRegularize_Maximum(t *testing.T) {
	ts := buildSimpleSeries()
	got := ts.Regularize(30, "s", "max", 0)
	base := mustTime(2025, 11, 10, 10, 0, 0)
	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 3),
		du(base.Add(60*time.Second), 10),
	)
	requireSeriesEq(t, got, want, 0)
}

func TestRegularize_Minimum(t *testing.T) {
	ts := buildSimpleSeries()
	got := ts.Regularize(30, "s", "min", 0)
	base := mustTime(2025, 11, 10, 10, 0, 0)
	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 1),
		du(base.Add(60*time.Second), 10),
	)
	requireSeriesEq(t, got, want, 0)
}

func TestRegularize_Last(t *testing.T) {
	ts := buildSimpleSeries()
	got := ts.Regularize(30, "s", "last", 0)
	base := mustTime(2025, 11, 10, 10, 0, 0)
	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 3),  // dernier dans fenêtre 1
		du(base.Add(60*time.Second), 10), // dernier dans fenêtre 2
	)
	requireSeriesEq(t, got, want, 0)
}

func TestRegularize_Sum(t *testing.T) {
	ts := buildSimpleSeries()
	got := ts.Regularize(30, "s", "sum", 0)
	base := mustTime(2025, 11, 10, 10, 0, 0)
	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 4),  // 1+3
		du(base.Add(60*time.Second), 10), // 10
	)
	requireSeriesEq(t, got, want, 0)
}

// Vérifie la normalisation de l'argument "per" (Seconds/sec/s, Minutes/min/m, Hours/hours/h).
func TestRegularize_PerAlias(t *testing.T) {
	base := mustTime(2025, 11, 10, 10, 0, 0)
	ts := TimeSeries{}
	ts.AddDataUnit(
		du(base.Add(1*time.Second), 1),
		du(base.Add(2*time.Second), 2),
	)
	ref := ts.Regularize(2, "s", "last", 0)
	alts := []string{"sec", "Seconds", "s"}
	for _, per := range alts {
		got := ts.Regularize(2, per, "last", 0)
		requireSeriesEq(t, got, ref, 0)
	}
}

// --------- Regularize: insertion de NaN quand il manque des fenêtres ---------

// Série avec trous: deux points dans la première fenêtre, puis un point très loin.
// freq=30s -> on s’attend à des NaN pour chaque fenêtre vide entre les deux zones.
func TestRegularize_InsertNaNs_WhenGaps_Average(t *testing.T) {
	base := mustTime(2025, 11, 10, 10, 0, 0)

	ts := TimeSeries{}
	// Fenêtre 1: [10:00:00, 10:00:30] -> valeurs 1 et 3
	ts.AddDataUnit(
		du(base.Add(5*time.Second), 1),
		du(base.Add(20*time.Second), 3),
	)
	// Grand trou, puis un point plus tard (≈ 10:02:10) pour déclencher l'insertion de NaN.
	// IMPORTANT: on ajoute encore un point plus tard pour satisfaire la condition
	// i < len(ts)-2 dans la boucle d'insertion de NaN du code existant.
	ts.AddDataUnit(
		du(base.Add(2*time.Minute+10*time.Second), 10), // 10:02:10
		du(base.Add(2*time.Minute+40*time.Second), 11), // 10:02:40
	)

	got := ts.Regularize(30, "s", "avg", 0) // tolérance ignorée -> 0

	// Fenêtres attendues:
	// 10:00:30 -> moyenne (1+3)/2 = 2
	// 10:01:00 -> NaN
	// 10:01:30 -> NaN
	// 10:02:00 -> NaN
	// 10:02:30 -> "avg" sur fenêtre [10:02:00,10:02:30] : seul 10:02:10 -> 10
	// 10:03:00 -> "avg" sur [10:02:30,10:03:00] : seul 10:02:40 -> 11

	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 2),           // 10:00:30
		du(base.Add(60*time.Second), math.NaN()),  // 10:01:00
		du(base.Add(90*time.Second), math.NaN()),  // 10:01:30
		du(base.Add(120*time.Second), math.NaN()), // 10:02:00
		du(base.Add(150*time.Second), 10),         // 10:02:30 (moyenne = 10)
		du(base.Add(180*time.Second), 11),         // 10:03:00 (moyenne = 11)
	)

	requireSeriesEq(t, got, want, 1e-12)
}

func TestRegularize_InsertNaNs_WhenGaps_Last(t *testing.T) {
	base := mustTime(2025, 11, 10, 10, 0, 0)

	ts := TimeSeries{}
	// Toujours deux points initiaux dans première fenêtre
	ts.AddDataUnit(
		du(base.Add(3*time.Second), 5),
		du(base.Add(25*time.Second), 7),
	)
	// Puis un point très tard + un autre pour satisfaire la condition i < len-2
	ts.AddDataUnit(
		du(base.Add(2*time.Minute+1*time.Second), 100),  // 10:02:01
		du(base.Add(2*time.Minute+50*time.Second), 101), // 10:02:50
	)

	got := ts.Regularize(30, "s", "last", 0)

	// Fenêtres attendues:
	// 10:00:30 -> last des points dans [10:00:00,10:00:30] => 7
	// 10:01:00 -> NaN
	// 10:01:30 -> NaN
	// 10:02:00 -> NaN
	// 10:02:30 -> last dans [10:02:00,10:02:30] => 100 (à 10:02:01)
	// 10:03:00 -> last dans [10:02:30,10:03:00] => 101 (à 10:02:50)

	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(30*time.Second), 7),
		du(base.Add(60*time.Second), math.NaN()),
		du(base.Add(90*time.Second), math.NaN()),
		du(base.Add(120*time.Second), math.NaN()),
		du(base.Add(150*time.Second), 100),
		du(base.Add(180*time.Second), 101),
	)

	requireSeriesEq(t, got, want, 1e-12)
}

// Cas où plusieurs fenêtres initiales sont vides (la série commence “tard”).
func TestRegularize_LeadingEmptyWindows_InsertNaNs(t *testing.T) {
	base := mustTime(2025, 11, 10, 10, 0, 0)

	// On commence la série à 10:01:05, rien en [10:00:00,10:01:00]
	ts := TimeSeries{}
	ts.AddDataUnit(
		du(base.Add(65*time.Second), 1),               // 10:01:05
		du(base.Add(85*time.Second), 3),               // 10:01:25
		du(base.Add(2*time.Minute+5*time.Second), 8),  // 10:02:05
		du(base.Add(2*time.Minute+40*time.Second), 9), // 10:02:40
	)

	got := ts.Regularize(30, "s", "max", 0)

	// D’après l’algo: RoundedStartTime(10:01:05,30s) -> 10:01:00,
	// comme datemin != start, on commence à 10:01:00.
	// Fenêtres:
	// 10:01:00..10:01:30 -> max(1,3)=3 -> point à 10:01:30
	// 10:01:30..10:02:00 -> (vide) -> NaN à 10:02:00
	// 10:02:00..10:02:30 -> 8 -> point à 10:02:30
	// 10:02:30..10:03:00 -> 9 -> point à 10:03:00

	want := TimeSeries{}
	want.AddDataUnit(
		du(base.Add(90*time.Second), 3),           // 10:01:30
		du(base.Add(120*time.Second), math.NaN()), // 10:02:00
		du(base.Add(150*time.Second), 8),          // 10:02:30
		du(base.Add(180*time.Second), 9),          // 10:03:00
	)

	requireSeriesEq(t, got, want, 1e-12)
}
