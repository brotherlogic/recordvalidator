package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	pb "github.com/brotherlogic/recordvalidator/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
)

var (
	resets = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_reset",
		Help: "The size of the print queue",
	}, []string{"scheme"})
)

func (s *Server) validateScheme(sc *pb.Scheme) {

	//Ensure we don't have something in togo and complete at the same time
	var nc []int32
	for _, c := range sc.GetCompletedIds() {
		found := false
		for _, c2 := range sc.GetInstanceIds() {
			if c2 == c {
				found = true
			}
		}

		if !found {
			nc = append(nc, c)
		}
	}
	sc.CompletedIds = nc

	for k, _ := range sc.GetCompleteDate() {
		found := false
		for _, c := range sc.GetCompletedIds() {
			if c == k {
				found = true
			}
		}

		if !found {
			delete(sc.CompleteDate, k)
		}
	}

	for _, c := range sc.GetCompletedIds() {
		found := false
		for cd := range sc.GetCompleteDate() {
			if cd == c {
				found = true
			}
		}

		if !found {
			sc.CompleteDate[c] = time.Now().Unix()
			resets.With(prometheus.Labels{"scheme": sc.GetName()}).Inc()
		}

		for key, _ := range sc.GetOrdering() {
			found := false
			for _, k := range sc.GetInstanceIds() {
				if k == key {
					found = true
				}
			}

			if !found {
				delete(sc.Ordering, key)
			}
		}
	}
}

func (s *Server) repick(ctx context.Context, sc *pb.Scheme) {
	s.CtxLog(ctx, fmt.Sprintf("Running repick with %v for %v", len(sc.InstanceIds), sc.GetName()))

	// Don't repick if there's nothing to pick
	if len(sc.InstanceIds) == 0 {
		return
	}

	// Don't repick if the scheme is not active
	if !sc.GetActive() {
		return
	}

	sc.CurrentPick = 0
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
	}

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
				s.CtxLog(ctx, fmt.Sprintf("Repick load failed: %v", err))
			}

			stillMatch, invalid, _ := scheme.filter(rec)
			s.CtxLog(ctx, fmt.Sprintf("%v is %v (%v) for %v", iid, invalid, stillMatch, scheme.name()))
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
				err := s.update(ctx, iid, sc.GetSoft(), sc.GetUnbox(), scheme.name())
				if err == nil {
					sc.CurrentPick = iid
					return
				}
			} else {
				//This record no longer applies
				in := []int32{}
				for _, tg := range sc.GetInstanceIds() {
					if tg != iid {
						in = append(in, tg)
					}
				}
				sc.InstanceIds = in
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
