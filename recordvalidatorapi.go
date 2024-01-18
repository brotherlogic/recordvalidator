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
	ropb "github.com/brotherlogic/recordsorganiser/proto"
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

	// If was_parents is empty, let's go off to the storage locker
	doneParents := false
	for _, scheme := range schemes.GetSchemes() {
		if scheme.GetName() == "was_parents" && len(scheme.GetInstanceIds()) == 0 {
			conn, err := s.FDialServer(ctx, "recordsorganiser")
			if err != nil {
				return nil, err
			}
			client := ropb.NewOrganiserServiceClient(conn)
			org, err := client.GetOrganisation(ctx, &ropb.GetOrganisationRequest{Locations: []*ropb.Location{{Name: "Sale 12 Inches"}}})
			if err != nil {
				return nil, err
			}
			found := false
			for _, loc := range org.GetLocations() {
				for _, entry := range loc.GetReleasesLocation() {
					if entry.GetSlot() > 1 {
						found = true
					}
				}
			}

			if !found {
				s.RaiseIssue("Trip To The Storage Locker", "Need to update")
			}
			doneParents = true
		}
	}
	for _, scheme := range schemes.GetSchemes() {
		if scheme.GetName() == "keepers" || scheme.GetName() == "keepers_single" {
			s.CtxLog(ctx, fmt.Sprintf("%v is being set to active: %v", scheme.GetName(), doneParents))
			scheme.Active = doneParents
		}
	}
	err = s.save(ctx, schemes)
	if err != nil {
		return nil, err
	}

	// Don't validate records until they've arrived
	r, rerr := s.loadRecord(ctx, in.GetInstanceId())
	if status.Code(rerr) == codes.OutOfRange { // Skip a deleted record
		return &rcpb.ClientUpdateResponse{}, nil
	}

	if rerr == nil {
		if r.GetMetadata().GetDateArrived() == 0 {
			return &rcpb.ClientUpdateResponse{}, nil
		}
	}

	if rerr != nil {
		return nil, rerr
	}

	s.CtxLog(ctx, "Running through schemes")
	picked := false
	for _, scheme := range schemes.GetSchemes() {
		var sg schemeGenerator
		for _, schemegen := range s.sgs {
			if schemegen.name() == scheme.GetName() && len(scheme.GetInstanceIds()) > 0 {
				sg = schemegen
			}
		}

		if sg == nil || !scheme.GetActive() {
			continue
		}

		if scheme.GetName() == "twelve_inch_sleeves" || scheme.GetName() == "seven_inch_sleeves" {
			scheme.Soft = true
		}

		current.With(prometheus.Labels{"scheme": scheme.GetName()}).Set(float64(scheme.GetCurrentPick()))

		if scheme.GetCurrentPick() == in.GetInstanceId() || scheme.GetCurrentPick() == 0 {
			if rerr != nil {
				if status.Convert(rerr).Code() == codes.OutOfRange {
					s.repick(ctx, scheme)
					picked = true
				} else {
					return nil, rerr
				}
			}

			marked, k, _ := sg.filter(r)

			s.CtxLog(ctx, fmt.Sprintf("[%v]: for %v -> %v,%v", in.GetInstanceId(), scheme.GetName(), marked, k))

			if (!marked || k) || scheme.GetCurrentPick() == 0 || r.GetMetadata().GetDateArrived() == 0 {
				if marked && scheme.GetSoft() && scheme.GetActive() {
					s.softValidate(ctx, in.GetInstanceId(), scheme.GetName())
				}
				s.repick(ctx, scheme)
				picked = true
			} else if r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE &&
				r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_SOFT_VALIDATE &&
				r.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_VALIDATE && time.Since(time.Unix(r.GetMetadata().GetLastValidate(), 0)) > time.Hour*24*3 {
			}

		}
	}

	s.CtxLog(ctx, "Searching for init")
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

	s.CtxLog(ctx, "Adding")
	// See if this needs to be added
	if rerr != nil {
		if status.Convert(rerr).Code() == codes.OutOfRange {
			return &rcpb.ClientUpdateResponse{}, s.save(ctx, schemes)
		}
		return nil, rerr
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
					app, done, order := scheme.filter(r)
					if sg.Ordering == nil {
						sg.Ordering = make(map[int32]float32)
					}

					s.CtxLog(ctx, fmt.Sprintf("Trying to add %v for %v -> %v, %v, %v", r.GetRelease().GetInstanceId(), sg.GetName(), app, done, order))

					if app {
						sg.Ordering[in.GetInstanceId()] = order
						if done {
							sg.CompletedIds = append(sg.CompletedIds, in.GetInstanceId())
							s.CtxLog(ctx, fmt.Sprintf("Added to complete %v", sg.GetName()))
						} else {
							sg.InstanceIds = append(sg.InstanceIds, in.GetInstanceId())
							s.CtxLog(ctx, fmt.Sprintf("Added to go %v", sg.GetName()))
						}
						picked = true
					}
				}
			}
		} else {
			for _, scheme := range s.sgs {
				if scheme.name() == sg.GetName() {
					app, done, _ := scheme.filter(r)
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

	// Clean existing schemes
	for _, sg := range schemes.GetSchemes() {
		for _, scheme := range s.sgs {
			if scheme.name() == sg.GetName() {
				var niids []int32
				for _, iid := range sg.GetInstanceIds() {
					if iid == in.GetInstanceId() {
						app, done, _ := scheme.filter(r)
						if !app || done {
							s.CtxLog(ctx, fmt.Sprintf("HARD REMOVE %v from %v", iid, sg.GetName()))
							sg.CompletedIds = append(sg.CompletedIds, iid)
							sg.CompleteDate[iid] = time.Now().Unix()
							picked = true
						} else {
							niids = append(niids, iid)
						}
					} else {
						niids = append(niids, iid)
					}
				}
				sg.InstanceIds = niids
			}
		}
	}

	var nerr error
	if picked {
		nerr = s.save(ctx, schemes)
	}
	return &rcpb.ClientUpdateResponse{}, nerr
}

func (s *Server) Force(ctx context.Context, req *pb.ForceRequest) (*pb.ForceResponse, error) {
	schemes, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	for _, scheme := range schemes.GetSchemes() {
		if scheme.GetName() == req.GetName() {
			scheme.CurrentPick = 0
			return &pb.ForceResponse{}, s.save(ctx, schemes)
		}
	}

	return nil, status.Errorf(codes.FailedPrecondition, "Not found")
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
