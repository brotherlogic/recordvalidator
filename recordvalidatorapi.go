package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	current = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_currentpick",
		Help: "The size of the print queue",
	}, []string{"scheme", "value"})
)

//ClientUpdate forces a move
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	schemes, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	picked := false
	for _, scheme := range schemes.GetSchemes() {
		var sg schemeGenerator
		for _, schemegen := range s.sgs {
			if schemegen.name() == scheme.GetName() {
				sg = schemegen
			}
		}

		if sg == nil {
			continue
		}

		current.With(prometheus.Labels{"scheme": scheme.GetName(), "value": fmt.Sprintf("%v", scheme.GetCurrentPick())}).Set(1)

		if scheme.GetCurrentPick() == in.GetInstanceId() || scheme.GetCurrentPick() == 0 {
			r, err := s.loadRecord(ctx, in.GetInstanceId())
			if err != nil {
				return nil, err
			}

			_, k := sg.filter(r)

			s.Log(fmt.Sprintf("Found pick (%v - %v) and activation is %v", in.GetInstanceId(), scheme.GetName(), k))
			if k || scheme.GetCurrentPick() == 0 {
				s.repick(ctx, scheme)
				picked = true
			} else if r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE {
				//This should be in pre-valid
				s.RaiseIssue("Valid miss", fmt.Sprintf("%v should be in prevalidate but it's actually in %v", r.GetRelease().GetInstanceId(), r.GetMetadata().GetCategory()))
			}

		}
	}

	for _, sg := range s.sgs {
		found := false

		for _, scheme := range schemes.GetSchemes() {
			if scheme.GetName() == sg.name() {
				found = true
			}
		}

		if !found {
			scheme, err := s.initScheme(ctx, sg)
			if err != nil {
				return nil, err
			}
			schemes.Schemes = append(schemes.Schemes, scheme)
			picked = true
		}

	}

	// See if this needs to be added
	rec, err := s.getRecord(ctx, in.GetInstanceId())
	if err != nil {
		return nil, err
	}
	for _, sg := range schemes.GetSchemes() {
		inS := false
		for _, id := range sg.GetInstanceIds() {
			if id == in.GetInstanceId() {
				inS = true
			}
		}

		if !inS {
			for _, scheme := range s.sgs {
				if scheme.name() == sg.GetName() {
					app, done := scheme.filter(rec)
					if app {
						if done {
							sg.CompletedIds = append(sg.CompletedIds, in.GetInstanceId())
						} else {
							sg.InstanceIds = append(sg.InstanceIds, in.GetInstanceId())
						}
						picked = true
					}
				}
			}
		}
	}

	var rerr error
	if picked {
		rerr = s.save(ctx, schemes)
	}
	return &rcpb.ClientUpdateResponse{}, rerr
}
