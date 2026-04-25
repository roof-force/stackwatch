package metrics

import (
	"testing"
	"time"
)

func makeSnap(total, drifted, errors int) Snapshot {
	return Snapshot{
		Total:      total,
		Drifted:    drifted,
		Errors:     errors,
		CheckedAt:  time.Now(),
		AvgLatency: 50 * time.Millisecond,
	}
}

func TestSnapshotStore_AddAndLatest(t *testing.T) {
	store := NewSnapshotStore(5)
	if _, ok := store.Latest(); ok {
		t.Fatal("expected no snapshot on empty store")
	}
	snap := makeSnap(10, 2, 0)
	store.Add(snap)
	got, ok := store.Latest()
	if !ok {
		t.Fatal("expected snapshot after Add")
	}
	if got.Total != 10 || got.Drifted != 2 {
		t.Errorf("unexpected snapshot values: %+v", got)
	}
}

func TestSnapshotStore_EvictsOldest(t *testing.T) {
	store := NewSnapshotStore(3)
	for i := 1; i <= 4; i++ {
		store.Add(makeSnap(i, 0, 0))
	}
	all := store.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(all))
	}
	if all[0].Total != 2 {
		t.Errorf("oldest should be evicted; first item Total = %d", all[0].Total)
	}
	if all[2].Total != 4 {
		t.Errorf("latest should be 4; got %d", all[2].Total)
	}
}

func TestSnapshotStore_DriftRate(t *testing.T) {
	store := NewSnapshotStore(10)
	if store.DriftRate() != 0 {
		t.Error("expected 0 drift rate on empty store")
	}
	store.Add(makeSnap(10, 5, 0)) // 50%
	store.Add(makeSnap(10, 0, 0)) // 0%
	// aggregate: 5/20 = 0.25
	rate := store.DriftRate()
	if rate != 0.25 {
		t.Errorf("expected drift rate 0.25, got %f", rate)
	}
}

func TestSnapshotStore_DefaultCapacity(t *testing.T) {
	store := NewSnapshotStore(0)
	for i := 0; i < 12; i++ {
		store.Add(makeSnap(1, 0, 0))
	}
	if len(store.All()) != 10 {
		t.Errorf("default capacity should be 10, got %d", len(store.All()))
	}
}

func TestSnapshotStore_AllReturnsCopy(t *testing.T) {
	store := NewSnapshotStore(5)
	store.Add(makeSnap(1, 0, 0))
	all := store.All()
	all[0].Total = 999
	got, _ := store.Latest()
	if got.Total == 999 {
		t.Error("All() should return a copy, not a reference")
	}
}
