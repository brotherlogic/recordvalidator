package main

import rcpb "github.com/brotherlogic/recordcollection/proto"

type schemeGenerator interface {
	// Does this record apply to the this filter
	filter(rec *rcpb.Record) (bool, bool)
	name() string
}

type keeperScheme struct{}

func (ks *keeperScheme) filter(rec *rcpb.Record) (bool, bool) {
	//Is it a keeper?, doess it have a width?
	return rec.GetMetadata().GetGoalFolder() == 466902, rec.GetMetadata().GetRecordWidth() > 0
}

func (ks *keeperScheme) name() string {
	return "keeper_width"
}
