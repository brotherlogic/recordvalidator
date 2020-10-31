package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
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
			return nil, fmt.Errorf("Cannot locate scheme %v", scheme.GetName())
		}

		if scheme.GetCurrentPick() == in.GetInstanceId() || scheme.GetCurrentPick() == 0 {
			r, err := s.loadRecord(ctx, in.GetInstanceId())
			if err != nil {
				return nil, err
			}

			_, k := sg.filter(r)
			if k || scheme.GetCurrentPick() == 0 {
				s.repick(scheme)
				picked = true
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

	var rerr error
	if picked {
		rerr = s.save(ctx, schemes)
	}
	return &rcpb.ClientUpdateResponse{}, rerr
}
