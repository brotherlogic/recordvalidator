package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordvalidator/proto"
)

var (
	current = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_currentpick",
		Help: "The size of the print queue",
	}, []string{"scheme"})
)

// ClientUpdate forces a move
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	schemes, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	picked := false
	for _, scheme := range schemes.GetSchemes() {
		var sg schemeGenerator
		for _, schemegen := range s.sgs {
			if schemegen.name() == scheme.GetName() && len(scheme.GetInstanceIds()) > 0 {
				sg = schemegen
			}
		}

		if sg == nil {
			continue
		}

		if scheme.GetName() == "twelve_inch_sleeves" || scheme.GetName() == "seven_inch_sleeves" {
			scheme.Soft = true
		}

		current.With(prometheus.Labels{"scheme": scheme.GetName()}).Set(float64(scheme.GetCurrentPick()))

		if scheme.GetCurrentPick() == in.GetInstanceId() || scheme.GetCurrentPick() == 0 {
			r, err := s.loadRecord(ctx, in.GetInstanceId())
			if err != nil {
				if status.Convert(err).Code() == codes.OutOfRange {
					s.repick(ctx, scheme)
				}
				return nil, err
			}

			marked, k, _ := sg.filter(r)
			s.CtxLog(ctx, fmt.Sprintf("[%v]: for %v -> %v,%v", in.GetInstanceId(), scheme.GetName(), marked, k))

			if (!marked || k) || scheme.GetCurrentPick() == 0 {
				if marked && scheme.GetSoft() && scheme.GetActive() {
					s.softValidate(ctx, in.GetInstanceId(), scheme.GetName())
				}
				s.repick(ctx, scheme)
				picked = true
			} else if r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE &&
				r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_SOFT_VALIDATE &&
				r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_VALIDATE && time.Since(time.Unix(r.GetMetadata().GetLastValidate(), 0)) > time.Hour*24*3 {
				//This should be in pre-valid unless it's just been validatd
				s.RaiseIssue(fmt.Sprintf("%v is a Valid miss", r.GetRelease().GetTitle()), fmt.Sprintf("%v should be in prevalidate but it's actually in %v (%v, %v and %v)", r.GetRelease().GetInstanceId(), r.GetMetadata().GetCategory(), marked, k, r.GetMetadata().GetLastValidate()))
			}

		}
	}

	for _, sg := range s.sgs {
		found := false

		for _, scheme := range schemes.GetSchemes() {
			if scheme.GetName() == sg.name() { //} && len(scheme.GetInstanceIds()) > 1 {
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
		if status.Convert(err).Code() == codes.OutOfRange {
			return &rcpb.ClientUpdateResponse{}, nil
		}
		return nil, err
	}
	mapper := ""
	for _, sg := range schemes.GetSchemes() {
		inS := false
		for _, id := range sg.GetInstanceIds() {
			if id == in.GetInstanceId() {
				inS = true
			}
		}

		mapper += fmt.Sprintf(" %v -> %v [%v]", sg.GetName(), inS, len(sg.GetInstanceIds()))
		if !inS {
			for _, scheme := range s.sgs {
				if scheme.name() == sg.GetName() {
					app, done, order := scheme.filter(rec)
					if sg.Ordering == nil {
						sg.Ordering = make(map[int32]float32)
					}

					if app {
						sg.Ordering[in.GetInstanceId()] = order
						if done {
							sg.CompletedIds = append(sg.CompletedIds, in.GetInstanceId())
						} else {
							sg.InstanceIds = append(sg.InstanceIds, in.GetInstanceId())
						}
						picked = true
					}
				}
			}
		} else {
			for _, scheme := range s.sgs {
				if scheme.name() == sg.GetName() {
					app, done, _ := scheme.filter(rec)
					if app && done {
						s.CtxLog(ctx, fmt.Sprintf("Removing %v from todo list (%v): %v,%v", in.GetInstanceId(), sg.GetName(), app, done))
						nc := []int32{}
						for _, iid := range sg.GetInstanceIds() {
							if iid != in.GetInstanceId() {
								nc = append(nc, iid)
							} else {
								sg.CompletedIds = append(sg.CompletedIds, iid)
								sg.CompleteDate[iid] = time.Now().Unix()
							}
						}
						sg.InstanceIds = nc
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

func (s *Server) GetScheme(ctx context.Context, req *pb.GetSchemeRequest) (*pb.GetSchemeResponse, error) {
	schemes, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	for _, scheme := range schemes.GetSchemes() {
		if req.GetInstanceId() > 0 {
			if scheme.GetCurrentPick() == req.GetInstanceId() {
				return &pb.GetSchemeResponse{Scheme: scheme}, nil
			}
		} else {
			if scheme.GetName() == req.GetName() {
				return &pb.GetSchemeResponse{Scheme: scheme}, nil
			}
		}
	}

	return nil, fmt.Errorf("Cannot find: %v", req.GetName())
}
