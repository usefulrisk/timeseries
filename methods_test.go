package timeseries

import (
	"math"
	"testing"
	"time"
)

// -------- helpers ------------------------------------------------------------

func almostEq(a, b float64, eps float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Abs(a-b) <= eps
}
func almostDurEq(a, b time.Duration, eps time.Duration) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= eps
}

// construit une série simple avec 4 points (dont un NaN) et ordre désordonné
// Chrons: t0, t0+25s, t0+10s, t0+12s  (exprès dans le désordre)
// Meas:   10,   20,     NaN,    5
func buildSampleSeries(t0 time.Time) *TimeSeries {
	ts := &TimeSeries{Name: "sample"}
	ts.AddData(t0, 10)
	ts.AddData(t0.Add(25*time.Second), 20)
	ts.AddData(t0.Add(10*time.Second), math.NaN())
	ts.AddData(t0.Add(12*time.Second), 5)
	return ts
}

// -------- tests --------------------------------------------------------------

func TestAddDataAndReset(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := &TimeSeries{Name: "addreset"}
	ts.AddData(t0, 1)
	ts.AddData(t0.Add(time.Second), 2)

	if got := len(ts.DataSeries); got != 2 {
		t.Fatalf("AddData len=2 expected, got %d", got)
	}
	if ts.Name != "addreset" {
		t.Fatalf("metadata (Name) should be preserved")
	}

	ts.Reset()
	if got := len(ts.DataSeries); got != 0 {
		t.Fatalf("Reset should empty DataSeries, got %d", got)
	}
	if ts.Name != "addreset" {
		t.Fatalf("Reset must keep metadata at TimeSeries level")
	}
}

func TestSortChronAscDesc(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := buildSampleSeries(t0)

	ts.SortChronAsc()
	for i := 1; i < len(ts.DataSeries); i++ {
		if !ts.DataSeries[i-1].Chron.Before(ts.DataSeries[i].Chron) &&
			!ts.DataSeries[i-1].Chron.Equal(ts.DataSeries[i].Chron) {
			t.Fatalf("SortChronAsc not sorted at %d", i)
		}
	}

	ts.SortChronDesc()
	for i := 1; i < len(ts.DataSeries); i++ {
		if !ts.DataSeries[i-1].Chron.After(ts.DataSeries[i].Chron) &&
			!ts.DataSeries[i-1].Chron.Equal(ts.DataSeries[i].Chron) {
			t.Fatalf("SortChronDesc not sorted at %d", i)
		}
	}
}

func TestSortMeasAscDesc(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := buildSampleSeries(t0)
	// remplace NaN par une valeur pour tester l’ordre total
	for i := range ts.DataSeries {
		if math.IsNaN(ts.DataSeries[i].Meas) {
			ts.DataSeries[i].Meas = 15
		}
	}

	ts.SortMeasAsc()
	for i := 1; i < len(ts.DataSeries); i++ {
		if !(ts.DataSeries[i-1].Meas <= ts.DataSeries[i].Meas) {
			t.Fatalf("SortMeasAsc not sorted at %d", i)
		}
	}

	ts.SortMeasDesc()
	for i := 1; i < len(ts.DataSeries); i++ {
		if !(ts.DataSeries[i-1].Meas >= ts.DataSeries[i].Meas) {
			t.Fatalf("SortMeasDesc not sorted at %d", i)
		}
	}
}

// --- helpers -----------------------------------------------------------------

func buildUnsortedSeries(t0 time.Time) *TimeSeries {
	ts := &TimeSeries{Name: "orig"}
	// ordre volontairement désordonné
	ts.AddData(t0.Add(20*time.Second), 2) // A
	ts.AddData(t0, 1)                     // B
	ts.AddData(t0.Add(10*time.Second), 3) // C
	return ts
}

func findIdxByChron(v []DataUnit, ch time.Time) int {
	for i := range v {
		if v[i].Chron.Equal(ch) {
			return i
		}
	}
	return -1
}

// --- tests -------------------------------------------------------------------

func TestCopy_DeepCopyAndSortedAndStats(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := buildUnsortedSeries(t0)

	// Mémorise l'ordre initial de l'original (doit rester inchangé)
	o0, o1, o2 := ts.DataSeries[0].Chron, ts.DataSeries[1].Chron, ts.DataSeries[2].Chron

	cp := ts.Copy()

	t.Run("original_not_mutated_order", func(t *testing.T) {
		// L'original ne doit pas être trié par Copy()
		if !ts.DataSeries[0].Chron.Equal(o0) ||
			!ts.DataSeries[1].Chron.Equal(o1) ||
			!ts.DataSeries[2].Chron.Equal(o2) {
			t.Fatalf("original order changed by Copy()")
		}
	})

	t.Run("copy_is_sorted_chron_asc", func(t *testing.T) {
		for i := 1; i < len(cp.DataSeries); i++ {
			prev, cur := cp.DataSeries[i-1].Chron, cp.DataSeries[i].Chron
			if !(prev.Before(cur) || prev.Equal(cur)) {
				t.Fatalf("copy not sorted asc at %d: %v !<= %v", i, prev, cur)
			}
		}
	})

	t.Run("deep_copy_no_alias_both_ways", func(t *testing.T) {
		// 1) modifier l'original ne doit pas impacter la copie (au même Chron)
		targetChron := ts.DataSeries[0].Chron
		ts.DataSeries[0].Meas = 999

		j := findIdxByChron(cp.DataSeries, targetChron)
		if j == -1 {
			t.Fatalf("same Chron not found in copy")
		}
		if cp.DataSeries[j].Meas == 999 {
			t.Fatalf("mutation leak: copy reflects change from original (alias?)")
		}

		// 2) modifier la copie ne doit pas impacter l'original (au même Chron)
		cp.DataSeries[j].Meas = 555

		i := findIdxByChron(ts.DataSeries, targetChron)
		if i == -1 {
			t.Fatalf("same Chron not found in original")
		}
		if ts.DataSeries[i].Meas == 555 {
			t.Fatalf("mutation leak: original reflects change from copy (alias?)")
		}
	})

	t.Run("name_preserved_comment_recomputed", func(t *testing.T) {
		// Name doit être recopié ; Comment peut être réécrit par Sort_Deltas_Stats()
		if cp.Name != ts.Name {
			t.Fatalf("Name not preserved in copy: got %q want %q", cp.Name, ts.Name)
		}
		// On ne teste pas l'égalité de Comment car ComputeBasicStats peut le modifier.
	})

	t.Run("stats_computed_on_copy", func(t *testing.T) {
		// La copie a subi Sort_Deltas_Stats() : vérif d'invariants simples
		if cp.Len != len(cp.DataSeries) {
			t.Fatalf("cp.Len = %d, want %d", cp.Len, len(cp.DataSeries))
		}
		if cp.Chmin.After(cp.Chmax) {
			t.Fatalf("cp.Chmin must be <= cp.Chmax")
		}
		// Si ComputeBasicStats fixe un message, il doit être non vide
		if cp.Comment == "" {
			t.Fatalf("cp.Comment should be set by ComputeBasicStats")
		}
	})
}

func TestDeltasFiller(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := &TimeSeries{Name: "deltas"}
	ts.AddData(t0, 1)
	ts.AddData(t0.Add(10*time.Second), 4)
	ts.AddData(t0.Add(22*time.Second), 10)

	ts.SortChronAsc()
	ts.DeltasFiller()

	if ts.DataSeries[0].Dchron != 0 || !almostEq(ts.DataSeries[0].Dmeas, 0, 0) {
		t.Fatalf("First delta must be zero (Dchron and Dmeas)")
	}

	if !almostDurEq(ts.DataSeries[1].Dchron, 10*time.Second, time.Nanosecond) {
		t.Fatalf("Unexpected Dchron[1], got %v", ts.DataSeries[1].Dchron)
	}
	if !almostEq(ts.DataSeries[1].Dmeas, 3, 1e-12) {
		t.Fatalf("Unexpected Dmeas[1], got %v", ts.DataSeries[1].Dmeas)
	}
	if !almostDurEq(ts.DataSeries[2].Dchron, 12*time.Second, time.Nanosecond) {
		t.Fatalf("Unexpected Dchron[2], got %v", ts.DataSeries[2].Dchron)
	}
	if !almostEq(ts.DataSeries[2].Dmeas, 6, 1e-12) {
		t.Fatalf("Unexpected Dmeas[2], got %v", ts.DataSeries[2].Dmeas)
	}
}

func TestSortDeltasStatsAndComputeBasicStats(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	ts := buildSampleSeries(t0) // contient un NaN

	ts.Sort_Deltas_Stats()

	// Chmin / Chmax / Len
	if !ts.Chmin.Equal(t0) {
		t.Fatalf("Chmin should be earliest time")
	}
	if !ts.Chmax.Equal(t0.Add(25 * time.Second)) {
		t.Fatalf("Chmax should be latest time")
	}
	if ts.Len != len(ts.DataSeries) {
		t.Fatalf("Len should equal number of points; got %d want %d", ts.Len, len(ts.DataSeries))
	}

	// NbreOfNaN (1 point NaN de mesure)
	if ts.NbreOfNaN != 1 {
		t.Fatalf("NbreOfNaN=1 expected, got %d", ts.NbreOfNaN)
	}

	// Msmin / Msmax & timestamps associés (NaN exclus des stats)
	if !almostEq(ts.Msmin, 5, 1e-12) {
		t.Fatalf("Msmin expected 5, got %v", ts.Msmin)
	}
	if !almostEq(ts.Msmax, 20, 1e-12) {
		t.Fatalf("Msmax expected 20, got %v", ts.Msmax)
	}
	if !ts.ChAtMsmin.Equal(t0.Add(12 * time.Second)) {
		t.Fatalf("ChAtMsmin wrong, got %v", ts.ChAtMsmin)
	}
	if !ts.ChAtMsmax.Equal(t0.Add(25 * time.Second)) {
		t.Fatalf("ChAtMsmax wrong, got %v", ts.ChAtMsmax)
	}

	// DChmin / DChmax / DChmean / DChmed (sur deltas > 0)
	// gaps (après tri): 10s, 2s, 13s  => min=2s, max=13s
	if !almostDurEq(ts.DChmin, 2*time.Second, time.Nanosecond) {
		t.Fatalf("DChmin expected 2s, got %v", ts.DChmin)
	}
	if !almostDurEq(ts.DChmax, 13*time.Second, time.Nanosecond) {
		t.Fatalf("DChmax expected 13s, got %v", ts.DChmax)
	}

	// DMs* existent car >= 3 points -> 2 deltas sur Meas après suppression du premier
	if math.IsNaN(ts.DMsmean) || math.IsNaN(ts.DMsmed) || math.IsNaN(ts.DMsstd) {
		t.Fatalf("DMs* stats should be defined (not NaN)")
	}

	// Comment final
	if ts.Comment != " Time Series ok." {
		t.Fatalf("Comment expected ' Time Series ok.', got %q", ts.Comment)
	}
}

func TestNewTsContainer(t *testing.T) {
	c := NewTsContainer()
	if c.Ts == nil {
		t.Fatalf("NewTsContainer should init map")
	}
	if len(c.Ts) != 0 {
		t.Fatalf("NewTsContainer map should be empty")
	}
}
