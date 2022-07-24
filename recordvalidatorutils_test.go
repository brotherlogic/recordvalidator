package main

import (
	"context"
	"testing"

	keystoreclient "github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/recordvalidator/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.test = true

	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.GoServer.KSclient.Save(context.Background(), SCHEMES, &pb.Schemes{})

	return s
}

func TestInitSchemeFailLoadRecord(t *testing.T) {
	s := InitTest()
	s.failRecordLoad = true

	_, err := s.initScheme(context.Background(), &keeperScheme{})
	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestRepick(t *testing.T) {
	s := InitTest()
	scheme := &pb.Scheme{
		InstanceIds: []int32{1, 2, 3},
		CurrentPick: 1,
	}

	s.repick(context.Background(), scheme)

	if scheme.CurrentPick == 2 {
		t.Errorf("Did not do a repick")
	}
}
