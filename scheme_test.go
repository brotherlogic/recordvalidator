package main

import (
	"testing"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func TestKeeperScheme(t *testing.T) {
	k := &keeperScheme{}

	if k.name() != "keeper_width" {
		t.Errorf("Huh")
	}

	f, p := k.filter(&rcpb.Record{Metadata: &rcpb.ReleaseMetadata{GoalFolder: 466902, RecordWidth: 4.5}})

	if !f || !p {
		t.Errorf("Problem filtering: %v and %v", f, p)
	}
}
