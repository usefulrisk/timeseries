package timeseries

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"
)

// --- Helpers ---------------------------------------------------------------

func mustNaNPtr(p *float64) bool {
	return p == nil
}
func mustFloatPtrEq(p *float64, v float64) bool {
	if p == nil {
		return false
	}
	return *p == v
}
func durNS(d time.Duration) int64 { return int64(d) }

// --- Tests: toPtrOrNil -----------------------------------------------------

func Test_toPtrOrNil_NaNBecomesNil(t *testing.T) {
	if got := toPtrOrNil(math.NaN()); got != nil {
		t.Fatalf("expected nil for NaN, got %v", got)
	}
}

func Test_toPtrOrNil_NumberBecomesPtr(t *testing.T) {
	val := 42.5
	got := toPtrOrNil(val)
	if got == nil || *got != val {
		t.Fatalf("expected pointer to %v, got %#v", val, got)
	}
}

// --- Tests: BasicStats.ToJSON ----------------------------------------------

func TestBasicStats_ToJSON_NilReceiver(t *testing.T) {
	var s *BasicStats = nil
	if s.ToJSON() != nil {
		t.Fatalf("nil receiver should return nil DTO")
	}
}

func TestBasicStats_ToJSON_ValuesCopied(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 13, 14, 15, 123456789, time.UTC)
	t1 := t0.Add(10 * time.Minute)
	t2 := t0.Add(20 * time.Minute)
	t3 := t0.Add(30 * time.Minute)
	t4 := t0.Add(40 * time.Minute)
	t5 := t0.Add(50 * time.Minute)

	in := &BasicStats{
		Len:        12,
		Chmin:      t0,
		ValAtChmin: -1.25,
		Chmax:      t5,
		ValAtChmax: 9.75,
		Chmed:      t2,
		Chmean:     t3,
		Chstd:      t4,
		Msmin:      -3.5,
		ChAtMsmin:  t1,
		Msmax:      7.25,
		ChAtMsmax:  t5,
		Msmean:     2.0,
		Msmed:      1.5,
		Msstd:      0.75,

		DChmin:     111 * time.Nanosecond,
		ChAtDChmin: t1,
		DChmax:     222 * time.Nanosecond,
		ChAtDchmax: t2,
		DChmean:    333 * time.Nanosecond,
		DChmed:     444 * time.Nanosecond,
		DChstd:     555 * time.Nanosecond,

		DMsmin:    -10.5,
		DMsmax:    10.5,
		DMsmed:    1.25,
		DMsmean:   2.25,
		DMsstd:    3.25,
		NbreOfNaN: 2,
	}

	out := in.ToJSON()
	if out == nil {
		t.Fatalf("expected non-nil DTO")
	}

	if out.Len != in.Len {
		t.Errorf("Len: got %d, want %d", out.Len, in.Len)
	}
	if !out.Chmin.Equal(in.Chmin) || !out.Chmax.Equal(in.Chmax) {
		t.Errorf("Chmin/Chmax mismatch")
	}
	if out.ValAtChmin != in.ValAtChmin || out.ValAtChmax != in.ValAtChmax {
		t.Errorf("ValAtCh* mismatch")
	}
	if !out.Chmed.Equal(in.Chmed) || !out.Chmean.Equal(in.Chmean) || !out.Chstd.Equal(in.Chstd) {
		t.Errorf("Chmed/Chmean/Chstd mismatch")
	}
	if out.Msmin != in.Msmin || out.Msmax != in.Msmax || out.Msmean != in.Msmean || out.Msmed != in.Msmed || out.Msstd != in.Msstd {
		t.Errorf("Ms* mismatch")
	}

	if out.DChminNS != int64(in.DChmin) ||
		out.DChmaxNS != int64(in.DChmax) ||
		out.DChmeanNS != int64(in.DChmean) ||
		out.DChmedNS != int64(in.DChmed) ||
		out.DChstdNS != int64(in.DChstd) {
		t.Errorf("DCh*NS mismatch")
	}
	if !out.ChAtDChmin.Equal(in.ChAtDChmin) || !out.ChAtDChmax.Equal(in.ChAtDchmax) {
		t.Errorf("ChAtDCh* mismatch")
	}

	if out.DMsmin != in.DMsmin || out.DMsmax != in.DMsmax || out.DMsmed != in.DMsmed || out.DMsmean != in.DMsmean || out.DMsstd != in.DMsstd {
		t.Errorf("DMs* mismatch")
	}
	if out.NbreOfNaN != in.NbreOfNaN {
		t.Errorf("NbreOfNaN: got %d, want %d", out.NbreOfNaN, in.NbreOfNaN)
	}
}

// --- Tests: TimeSeries.ToJSON ----------------------------------------------

func TestTimeSeries_ToJSON_MapsNaNToNilAndCopiesFields(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 9, 0, 0, 0, time.UTC)
	du := []DataUnit{
		{
			Chron:  t0,
			Meas:   1.25,
			Dchron: 5 * time.Second,
			Dmeas:  math.NaN(),
			Status: StOK,
		},
		{
			Chron:  t0.Add(10 * time.Second),
			Meas:   math.NaN(),
			Dchron: 10 * time.Second,
			Dmeas:  -2.5,
			Status: StMissing,
		},
	}

	ts := &TimeSeries{
		Name:       "tsA",
		Comment:    "demo",
		DataSeries: du,
		BasicStats: BasicStats{
			Len: 2,
		},
	}

	dto := ts.ToJSON()
	if dto == nil {
		t.Fatalf("expected non-nil DTO")
	}

	// Name / Comment
	if dto.Name != ts.Name || dto.Comment != ts.Comment {
		t.Errorf("Name/Comment mismatch: got %q/%q", dto.Name, dto.Comment)
	}

	// Chron
	if len(dto.Chron) != len(du) || !dto.Chron[0].Equal(du[0].Chron) || !dto.Chron[1].Equal(du[1].Chron) {
		t.Errorf("Chron slice mismatch")
	}

	// Meas: [1.25, nil]
	if !mustFloatPtrEq(dto.Meas[0], 1.25) {
		t.Errorf("Meas[0] expected *1.25, got %#v", dto.Meas[0])
	}
	if !mustNaNPtr(dto.Meas[1]) {
		t.Errorf("Meas[1] expected nil for NaN, got %#v", dto.Meas[1])
	}

	// DchronNS
	wantDchron := []int64{durNS(5 * time.Second), durNS(10 * time.Second)}
	if !reflect.DeepEqual(dto.DchronNS, wantDchron) {
		t.Errorf("DchronNS mismatch: got %v, want %v", dto.DchronNS, wantDchron)
	}

	// Dmeas: [nil, -2.5]
	if !mustNaNPtr(dto.Dmeas[0]) {
		t.Errorf("Dmeas[0] expected nil for NaN, got %#v", dto.Dmeas[0])
	}
	if !mustFloatPtrEq(dto.Dmeas[1], -2.5) {
		t.Errorf("Dmeas[1] expected *-2.5, got %#v", dto.Dmeas[1])
	}

	// Status
	wantStatus := []StatusCode{StOK, StMissing}
	if !reflect.DeepEqual(dto.Status, wantStatus) {
		t.Errorf("Status mismatch: got %v, want %v", dto.Status, wantStatus)
	}

	// Stats present
	if dto.Stats == nil || dto.Stats.Len != 2 {
		t.Errorf("Stats not copied properly")
	}

	// Bonus: ensure JSON marshaling doesn't choke on NaN (should become null)
	raw, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	if !jsonContains(raw, `"Meas":[{"`[:7]) { // quick sanity (string presence). We’ll do a simpler check below.
		// no-op: this guard avoids unused warning if you delete jsonContains
	}
	// simple check for "null" presence coming from NaN → nil
	if !containsBytes(raw, []byte("null")) {
		t.Errorf("expected marshaled JSON to contain null for NaN -> nil mapping, got: %s", string(raw))
	}
}

// --- Tests: TsContainer.ToJSON ---------------------------------------------

func TestTsContainer_ToJSON_CopiesSeriesAndMetadata(t *testing.T) {
	t0 := time.Date(2025, 11, 12, 10, 0, 0, 0, time.UTC)
	ts := &TimeSeries{
		Name:    "series-1",
		Comment: "ok",
		DataSeries: []DataUnit{
			{Chron: t0, Meas: 3.14, Dchron: time.Second, Dmeas: math.NaN(), Status: StOK},
		},
	}
	c := &TsContainer{
		Name:    "container-A",
		Comment: "hello",
		Ts: map[string]*TimeSeries{
			"a": ts,
		},
	}

	dto := c.ToJSON()
	if dto == nil {
		t.Fatalf("expected non-nil container DTO")
	}
	if dto.Name != c.Name || dto.Comment != c.Comment {
		t.Errorf("Name/Comment mismatch on container")
	}
	if dto.Series == nil || len(dto.Series) != 1 {
		t.Fatalf("Series map missing/size mismatch")
	}
	inner, ok := dto.Series["a"]
	if !ok || inner == nil {
		t.Fatalf("Series[\"a\"] missing")
	}
	if inner.Name != ts.Name || inner.Comment != ts.Comment {
		t.Errorf("inner series Name/Comment mismatch")
	}
	// quick spot checks
	if len(inner.Chron) != 1 || !inner.Chron[0].Equal(t0) {
		t.Errorf("inner series Chron mismatch")
	}
	if !mustFloatPtrEq(inner.Meas[0], 3.14) {
		t.Errorf("inner series Meas[0] mismatch")
	}
	if inner.DchronNS[0] != int64(time.Second) {
		t.Errorf("inner series DchronNS[0] mismatch")
	}
	if !mustNaNPtr(inner.Dmeas[0]) {
		t.Errorf("inner series Dmeas[0] should be nil (from NaN)")
	}
}

// --- Small byte helpers ----------------------------------------------------

func containsBytes(haystack, needle []byte) bool {
	return len(haystack) >= len(needle) && (reflect.DeepEqual(haystack, needle) || indexBytes(haystack, needle) >= 0)
}
func indexBytes(b, subslice []byte) int {
	// tiny naive search; good enough for tests
	for i := 0; i+len(subslice) <= len(b); i++ {
		if reflect.DeepEqual(b[i:i+len(subslice)], subslice) {
			return i
		}
	}
	return -1
}
func jsonContains(_ []byte, _ string) bool { return true } // placeholder to keep earlier guard concise
