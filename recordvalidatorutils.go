package main

import (
	"fmt"
	"math/rand"
	"sort"
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
	switch sc.Order {
	case pb.Scheme_RANDOM:
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(sc.InstanceIds), func(i, j int) { sc.InstanceIds[i], sc.InstanceIds[j] = sc.InstanceIds[j], sc.InstanceIds[i] })
	case pb.Scheme_ORDER:
		sort.SliceStable(sc.InstanceIds, func(i, j int) bool {
			return sc.InstanceIds[i] < sc.InstanceIds[j]
		})
	case pb.Scheme_REVERSE_ORDER:
		sort.SliceStable(sc.InstanceIds, func(i, j int) bool {
			return sc.InstanceIds[i] < sc.InstanceIds[j]
		})
	}

	var scheme schemeGenerator
	for _, sch := range s.sgs {
		if sch.name() == sc.GetName() {
			scheme = sch
		}
	}

	found := false
	for i, iid := range sc.InstanceIds {
		if iid == 19866960 {
			s.Log(fmt.Sprintf("Found %v at %v", 19866960, i))
			found = true
		}
	}
	s.Log(fmt.Sprintf("Did find %v: %v", 19866960, found))

	if scheme != nil {
		// Find the first instance that is still relevant
		for i, iid := range sc.InstanceIds {
			rec, err := s.getRecord(ctx, iid)
			if err != nil {
				s.Log(fmt.Sprintf("Repick load failed: %v", err))
			}

			_, invalid := scheme.filter(rec)
			s.Log(fmt.Sprintf("%v is %v", iid, invalid))
			if invalid {
				s.Log(fmt.Sprintf("19866960 invalid %v from %v instead", iid, i))
				in := []int32{}
				for _, tg := range sc.GetInstanceIds() {
					if tg != iid {
						in = append(in, tg)
					}
				}
				sc.InstanceIds = in
				sc.CompletedIds = append(sc.CompletedIds, iid)
			} else {
				s.Log(fmt.Sprintf("19866960 picking %v from %v instead", iid, i))
				err := s.update(ctx, iid, sc.GetUnbox())
				if err == nil {
					sc.CurrentPick = iid
					return
				}
			}
		}
	}
}

func (s *Server) initScheme(ctx context.Context, sg schemeGenerator) (*pb.Scheme, error) {
	s.Log(fmt.Sprintf("Init shceme: %v", sg.name()))
	defer s.Log(fmt.Sprintf("Init of %v complete", sg.name()))
	scheme := &pb.Scheme{Name: sg.name(), StartTime: time.Now().Unix()}

	if sg.name() == "old_age" {
		scheme.Unbox = true
		scheme.Order = pb.Scheme_ORDER
	}

	if sg.name() == "new_age" {
		scheme.Unbox = true
		scheme.Order = pb.Scheme_REVERSE_ORDER
	}

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
