package main

import (
	"context"
	"testing"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func TestBlankUpdate(t *testing.T) {
	s := InitTest()
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}
