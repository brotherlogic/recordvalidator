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
	return rec.GetMetadata().GetGoalFolder() == 2259637, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (ks *keeperScheme) name() string {
	return "keeper_width"
}

type cdScheme struct{}

func (cds *cdScheme) filter(rec *rcpb.Record) (bool, bool) {
	//Is it a cd?, doess it have a width?
	return rec.GetRelease().GetFolderId() == 242018, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (cds *cdScheme) name() string {
	return "cd_width"
}
