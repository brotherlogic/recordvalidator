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

	_, err = s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{InstanceId: 12})

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestFailSave(t *testing.T) {
	s := InitTest()
	s.failSave = true
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestFailInit(t *testing.T) {
	s := InitTest()
	s.failLoadAll = true
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestFailLoadUpdate(t *testing.T) {
	s := InitTest()
	s.failLoad = true
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestMissingSG(t *testing.T) {
	s := InitTest()
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})
	if err != nil {
		t.Errorf("Bad init update: %v", err)
	}
	s.sgs = []schemeGenerator{}
	_, err = s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})
	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestFailRecordLooad(t *testing.T) {
	s := InitTest()
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}

	s.failRecordLoad = true

	_, err = s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{InstanceId: 12})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}
