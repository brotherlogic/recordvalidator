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
	s.Log(fmt.Sprintf("Running repick"))
	ntg := []int32{}

	for _, tg := range sc.GetInstanceIds() {
		if tg != sc.GetCurrentPick() {
			ntg = append(ntg, tg)
		}
	}
	sc.InstanceIds = ntg
	sc.CompletedIds = append(sc.CompletedIds, sc.GetCurrentPick())
	sc.CompleteDate[sc.GetCurrentPick()] = time.Now().Unix()

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
			return sc.InstanceIds[i] > sc.InstanceIds[j]
		})
	case pb.Scheme_GIVEN_ORDER:
		sort.SliceStable(sc.InstanceIds, func(i, j int) bool {
			return sc.Ordering[sc.InstanceIds[i]] < sc.Ordering[sc.InstanceIds[j]]
		})

		s.Log(fmt.Sprintf("GIVEN PICK %v -> %v with %v", sc.InstanceIds[0], sc.InstanceIds[len(sc.InstanceIds)-1], sc.Ordering))
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
	s.Log(fmt.Sprintf("Did find %v: %v -> %v", 19866960, found, scheme))

	if scheme != nil {
		// Find the first instance that is still relevant
		for _, iid := range sc.InstanceIds {
			rec, err := s.getRecord(ctx, iid)
			if err != nil {
				s.Log(fmt.Sprintf("Repick load failed: %v", err))
			}

			stillMatch, invalid, _ := scheme.filter(rec)
			s.Log(fmt.Sprintf("%v is %v", iid, invalid))
			if invalid {
				in := []int32{}
				for _, tg := range sc.GetInstanceIds() {
					if tg != iid {
						in = append(in, tg)
					}
				}
				sc.InstanceIds = in
				sc.CompleteDate[iid] = time.Now().Unix()
				sc.CompletedIds = append(sc.CompletedIds, iid)
			} else if stillMatch {
				s.Log(fmt.Sprintf("Updating %v -> %v", iid, scheme.name()))
				err := s.update(ctx, iid, sc.GetSoft(), sc.GetUnbox(), scheme.name())
				if err == nil {
					sc.CurrentPick = iid
					return
				}
			}
		}
	}
}

func (s *Server) initScheme(ctx context.Context, sg schemeGenerator) (*pb.Scheme, error) {
	var scheme *pb.Scheme
	s.CtxLog(ctx, fmt.Sprintf("Init shceme: %v", sg.name()))
	defer s.CtxLog(ctx, fmt.Sprintf("Init of %v complete -> %v", sg.name(), len(scheme.GetInstanceIds())))

	schemes, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	var scs []*pb.Scheme
	var seen []string
	for _, sc := range schemes.GetSchemes() {
		seens := false
		for _, se := range seen {
			if se == sc.GetName() {
				seens = true
			}
		}

		if sc.GetName() == sg.name() {
			scheme = sc
		}

		if !seens {
			scs = append(scs)
			seen = append(seen, sc.GetName())
		}
	}

	if scheme == nil {
		scheme = &pb.Scheme{Name: sg.name(), StartTime: time.Now().Unix(), Ordering: make(map[int32]float32)}
	}

	if sg.name() == "old_age" || sg.name() == "old_age_twelves" {
		scheme.Unbox = true
		scheme.Order = pb.Scheme_ORDER
	}

	if sg.name() == "new_age" || sg.name() == "new_age_twelves" {
		scheme.Unbox = true
		scheme.Order = pb.Scheme_REVERSE_ORDER
	}

	if sg.name() == "seven_inch_sleeves" {
		scheme.Soft = true
	}

	if sg.name() == "bad_ones" {
		scheme.Order = pb.Scheme_GIVEN_ORDER
	}

	if sg.name() == "bad_ones_twelves" {
		scheme.Unbox = true
		scheme.Order = pb.Scheme_GIVEN_ORDER
	}

	//Init everything empty
	iids, err := s.getAllRecords(ctx)
	if err != nil {
		return nil, err
	}
	for _, iid := range iids {
		r, err := s.loadRecord(ctx, iid)
		if err != nil {
			return nil, err
		}

		f, p, o := sg.filter(r)
		if scheme.Ordering == nil {
			scheme.Ordering = make(map[int32]float32)
		}
		scheme.Ordering[iid] = o
		s.CtxLog(ctx, fmt.Sprintf("Found %v -> %v,%v", r.GetRelease().GetInstanceId(), f, p))
		if f {
			scheme.InstanceIds = append(scheme.InstanceIds, iid)
			if p {
				scheme.CompletedIds = append(scheme.CompletedIds, iid)
			}
		}
	}
	return scheme, nil
}
