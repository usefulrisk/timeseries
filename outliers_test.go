package timeseries

import (
	"math"
	"reflect"
	"testing"
	"time"
)

// --- Helpers ---------------------------------------------------------------

func mkTS(vals ...float64) TimeSeries {
	ts := TimeSeries{Name: "T", Comment: "fixture"}
	t0 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	for i, v := range vals {
		du := DataUnit{
			Chron:  t0.Add(time.Minute * time.Duration(i)),
			Meas:   v,
			Status: StOK,
		}
		ts.AddDataUnit(du)
	}
	return ts
}

func measSlice(ts TimeSeries) []float64 {
	out := make([]float64, len(ts.DataSeries))
	for i, du := range ts.DataSeries {
		out[i] = du.Meas
	}
	return out
}

func chronIsSortedAsc(ts TimeSeries) bool {
	for i := 1; i < len(ts.DataSeries); i++ {
		if ts.DataSeries[i-1].Chron.After(ts.DataSeries[i].Chron) {
			return false
		}
	}
	return true
}

// --- Tests -----------------------------------------------------------------

func TestRemoveOutbounds_Basic(t *testing.T) {
	// Input: 1,2,3,4,5  -> keep [2..4], reject 1 and 5
	in := mkTS(1, 2, 3, 4, 5)

	clean, rej := in.RemoveOutbounds(2, 4, "test bounds")

	if !chronIsSortedAsc(in) {
		t.Fatalf("original TimeSeries should be restored in chronological order")
	}
	if !chronIsSortedAsc(clean) || !chronIsSortedAsc(rej) {
		t.Fatalf("outputs must be chronologically sorted")
	}

	gotClean := measSlice(clean)
	gotRej := measSlice(rej)

	wantClean := []float64{2, 3, 4}
	wantRej := []float64{1, 5}

	if !reflect.DeepEqual(gotClean, wantClean) {
		t.Fatalf("cleaned Meas = %v, want %v", gotClean, wantClean)
	}
	if !reflect.DeepEqual(gotRej, wantRej) {
		t.Fatalf("rejected Meas = %v, want %v", gotRej, wantRej)
	}

	// NOTE: Dans le code actuel, Status est forcé à StOutlier
	// aussi bien pour les gardés que les rejetés. On vérifie juste
	// que le split est correct, pas le Status.
}

func TestZscoreCleaning_Level1(t *testing.T) {
	// Mean ~= 0, Std ~= 1 ; lvl=1 -> garder [-1, +1] (approx)
	in := mkTS(-2, -0.5, 0, 0.5, 2)

	clean, rej := in.ZscoreCleaning(1.0)

	if !chronIsSortedAsc(clean) || !chronIsSortedAsc(rej) {
		t.Fatalf("outputs must be chronologically sorted")
	}

	gotClean := measSlice(clean)
	gotRej := measSlice(rej)

	// On tolère les arrondis : on s'attend à garder {-0.5, 0, 0.5}
	wantClean := []float64{-0.5, 0, 0.5}
	wantRej := []float64{-2, 2}

	if !reflect.DeepEqual(gotClean, wantClean) {
		t.Fatalf("cleaned Meas = %v, want %v", gotClean, wantClean)
	}
	if !reflect.DeepEqual(gotRej, wantRej) {
		t.Fatalf("rejected Meas = %v, want %v", gotRej, wantRej)
	}
}

func TestPeirceOutlierRemoval_Simple(t *testing.T) {
	// Données avec un outlier évident (100)
	in := mkTS(10, 11, 9, 10.5, 100, 9.5, 10.2)
	clean, rej := in.PeirceOutlierRemoval()

	if len(rej.DataSeries) == 0 {
		t.Fatalf("expected at least one rejected data point")
	}

	// Le plus probable est que 100 soit rejeté
	found100 := false
	for _, du := range rej.DataSeries {
		if du.Meas == 100 {
			found100 = true
			break
		}
	}
	if !found100 {
		t.Fatalf("expected 100 to be rejected by Peirce criterion, got rejected=%v", measSlice(rej))
	}

	// Sanity: cleaned + rejected = total
	if len(clean.DataSeries)+len(rej.DataSeries) != len(in.DataSeries) {
		t.Fatalf("sizes don't add up: clean=%d rej=%d in=%d",
			len(clean.DataSeries), len(rej.DataSeries), len(in.DataSeries))
	}
}

func TestMerge(t *testing.T) {
	a := mkTS(1, 2, 3)
	b := mkTS(10, 20)
	m := Merge(&a, &b)

	if len(m.DataSeries) != len(a.DataSeries)+len(b.DataSeries) {
		t.Fatalf("merge length = %d, want %d", len(m.DataSeries), len(a.DataSeries)+len(b.DataSeries))
	}
	// Concat exact dans l'ordre
	got := measSlice(m)
	want := []float64{1, 2, 3, 10, 20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged values = %v, want %v", got, want)
	}
}

func TestRemovedNonValid(t *testing.T) {
	ts := TimeSeries{Name: "T"}
	t0 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	ts.AddDataUnit(DataUnit{Chron: t0, Meas: 1, Status: StOK})
	ts.AddDataUnit(DataUnit{Chron: t0.Add(time.Minute), Meas: 2, Status: StOutlier})
	ts.AddDataUnit(DataUnit{Chron: t0.Add(2 * time.Minute), Meas: 3, Status: StOK})

	clean, rej := ts.RemovedNonValid()

	if got := measSlice(clean); !reflect.DeepEqual(got, []float64{1, 3}) {
		t.Fatalf("cleaned = %v, want [1 3]", got)
	}
	if got := measSlice(rej); !reflect.DeepEqual(got, []float64{2}) {
		t.Fatalf("rejected = %v, want [2]", got)
	}
}

func TestRemoveElements_IntendedSemantics(t *testing.T) {
	// Intention (d’après le nom): "RemoveElements(idx...)" devrait
	// retirer ces indices de la série principale -> ils doivent aller
	// dans "rejected". Le reste dans "cleaned".
	in := mkTS(10, 20, 30, 40, 50)
	idx := []int{1, 3} // valeurs 20 et 40

	clean, rej := in.RemoveElements(idx...)

	// NOTE: Le code actuel fait l’inverse (les idx vont dans cleaned).
	// Ce test documente le comportement attendu. Il échouera si la
	// fonction n’est pas corrigée.
	wantClean := []float64{10, 30, 50}
	wantRej := []float64{20, 40}

	if got := measSlice(clean); !reflect.DeepEqual(got, wantClean) {
		t.Fatalf("cleaned = %v, want %v (les éléments aux indices idx ne doivent PAS y être)", got, wantClean)
	}
	if got := measSlice(rej); !reflect.DeepEqual(got, wantRej) {
		t.Fatalf("rejected = %v, want %v (les éléments aux indices idx doivent être rejetés)", got, wantRej)
	}
}

func TestFindInIntSlice(t *testing.T) {
	idx := []int{4, 7, 9}
	if pos := FindInIntSlice(idx, 7); pos == -1 {
		t.Fatalf("expected to find 7 in %v", idx)
	}
	if pos := FindInIntSlice(idx, 3); pos != -1 {
		t.Fatalf("did not expect to find 3 in %v, got pos=%d", idx, pos)
	}
}

func TestFindInFloatSlice(t *testing.T) {
	vals := []float64{math.NaN(), 1.5, 2.0}
	if pos := FindInFloatSlice(vals, 2.0); pos == -1 {
		t.Fatalf("expected to find 2.0 in %v", vals)
	}
	// Attention: NaN != NaN, donc la recherche d’un NaN renverra -1
	if pos := FindInFloatSlice(vals, math.NaN()); pos != -1 {
		t.Fatalf("should not find NaN by equality, got pos=%d", pos)
	}
}

func TestChronToArr_MeasToArr(t *testing.T) {
	ts := mkTS(1, 2, 3)
	ch := ts.ChronToArr()
	ms := ts.MeasToArr()

	if len(ch) != 3 || len(ms) != 3 {
		t.Fatalf("unexpected lengths: chron=%d meas=%d", len(ch), len(ms))
	}
	if !ch[0].Before(ch[1]) || !ch[1].Before(ch[2]) {
		t.Fatalf("chronological order expected")
	}
	if !reflect.DeepEqual(ms, []float64{1, 2, 3}) {
		t.Fatalf("measures = %v, want [1 2 3]", ms)
	}
}

// --- Sanity tests around Percentile/Mean/StdDev if present -----------------
// Ces tests valident un minimum pour éviter des surprises en amont de PercCleaning/ZscoreCleaning.

func TestMeanStdDev_Sanity(t *testing.T) {
	in := []float64{2, 4, 4, 4, 5, 5, 7, 9} // mean=5, std=2 (population std si c'est l'habituel)
	m, err := Mean(append([]float64{}, in...))
	if err != nil {
		t.Fatalf("Mean error: %v", err)
	}
	if math.Abs(m-5) > 1e-9 {
		t.Fatalf("Mean=%v, want 5", m)
	}
	s, err := StdDev(append([]float64{}, in...))
	if err != nil {
		t.Fatalf("StdDev error: %v", err)
	}
	// On ne bloque pas si l’implémentation est sample vs population,
	// mais on vérifie juste que c’est raisonnable (proche de 2).
	if math.Abs(s-2) > 0.25 && math.Abs(s-2.138089935) > 0.25 {
		t.Fatalf("StdDev=%v, expected around 2 (population) or 2.138 (sample)", s)
	}
}

func TestPercentile_Sanity(t *testing.T) {
	// Définition testée :
	//  - On renvoie la plus grande valeur x telle que F_n(x) <= p
	//  - Équivalent : index = floor(p/100 * n), puis on prend l’élément d’indice index-1 (1-based)
	//  - Sur [10,20,30,40,50], p=50% => floor(2.5)=2 => 20.
	in := []float64{10, 20, 30, 40, 50}

	type tc struct {
		p    float64
		want float64
	}
	cases := []tc{
		{p: 1, want: 10},  // floor(0.05)=0 -> borne basse → selon impl, on attend le 1er
		{p: 10, want: 10}, // floor(0.5)=0 -> 10
		{p: 20, want: 10}, // floor(1.0)=1 -> 10
		{p: 25, want: 10}, // floor(1.25)=1 -> 10
		{p: 40, want: 20}, // floor(2.0)=2 -> 20
		{p: 50, want: 20}, // floor(2.5)=2 -> 20  (ton attente)
		{p: 60, want: 30}, // floor(3.0)=3 -> 30
		{p: 75, want: 30}, // floor(3.75)=3 -> 30
		{p: 99, want: 40}, // proche du max
	}

	for _, c := range cases {
		p := append([]float64{}, in...) // certaines implémentations trient in-place
		got, err := Percentile(p, c.p)
		if err != nil {
			t.Fatalf("Percentile error for p=%.2f: %v", c.p, err)
		}
		if got != c.want {
			t.Fatalf("Percentile(%.2f)=%v, want %v (définition left-continuous ≤p%%)", c.p, got, c.want)
		}
	}
}
