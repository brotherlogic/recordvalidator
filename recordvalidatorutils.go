package main

import (
	"fmt"
	"math/rand"
	"time"

	pb "github.com/brotherlogic/recordvalidator/proto"
	"golang.org/x/net/context"
)

func (s *Server) repick(ctx context.Context, sc *pb.Scheme) {
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

	var scheme schemeGenerator
	for _, sch := range s.sgs {
		if sch.name() == sc.GetName() {
			scheme = sch
		}
	}

	if scheme != nil {
		// Find the first instance that is still relevant
		for _, iid := range sc.InstanceIds {
			rec, err := s.getRecord(ctx, iid)
			if err != nil {
				s.Log(fmt.Sprintf("Repick load failed: %v", err))
			}

			_, invalid := scheme.filter(rec)
			if invalid {
				in := []int32{}
				for _, tg := range sc.GetInstanceIds() {
					if tg != iid {
						in = append(in, tg)
					}
				}
				sc.InstanceIds = in
				sc.CompletedIds = append(sc.CompletedIds, iid)
			} else {
				err := s.update(ctx, iid)
				if err == nil {
					sc.CurrentPick = iid
					return
				}
			}
		}
	}
}

func (s *Server) initScheme(ctx context.Context, sg schemeGenerator) (*pb.Scheme, error) {
	scheme := &pb.Scheme{Name: sg.name(), StartTime: time.Now().Unix()}
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
