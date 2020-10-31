package main

import (
	"math/rand"
	"time"

	pb "github.com/brotherlogic/recordvalidator/proto"
	"golang.org/x/net/context"
)

func (s *Server) repick(sc *pb.Scheme) {
	ntg := []int32{}

	for _, tg := range sc.GetInstanceIds() {
		if tg != sc.GetCurrentPick() {
			ntg = append(ntg, tg)
		}
	}
	sc.InstanceIds = ntg
	sc.CompletedIds = append(sc.CompletedIds, sc.GetCurrentPick())

	// Shuffle the instance ids
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(sc.InstanceIds), func(i, j int) { sc.InstanceIds[i], sc.InstanceIds[j] = sc.InstanceIds[j], sc.InstanceIds[i] })

	if len(sc.InstanceIds) > 0 {
		sc.CurrentPick = sc.InstanceIds[0]
	}
}

func (s *Server) initScheme(ctx context.Context, sg schemeGenerator) (*pb.Scheme, error) {
	scheme := &pb.Scheme{Name: sg.name()}
	iids, err := s.getAllRecords(ctx)
	if err != nil {
		return nil, err
	}
	for _, iid := range iids {
		r, err := s.loadRecord(ctx, iid)
		if err != nil {
			return nil, err
		}

		f, p := sg.filter(r)
		if f {
			if p {
				scheme.CompletedIds = append(scheme.CompletedIds, iid)
			} else {
				scheme.InstanceIds = append(scheme.InstanceIds, iid)
			}
		}
	}
	return scheme, nil
}
