package main

import (
	"time"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

type schemeGenerator interface {
	// Does this record apply to the this filter
	filter(rec *rcpb.Record) (bool, bool)
	name() string
}

type fullScheme struct{}

func (fs *fullScheme) filter(rec *rcpb.Record) (bool, bool) {
	// Can't process these
	if rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_PARENTS {
		return false, true
	}

	// Sold Digital recordings should be included ehre
	marked := true
	if rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_LISTED_TO_SELL ||
		rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_STAGED_TO_SELL ||
		rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_SOLD_ARCHIVE {
		for _, f := range rec.Release.GetFormats() {
			if f.Name == "CD" || f.Name == "File" || f.Name == "CDr" {
				marked = true
				break
			}

			marked = false
		}
	}

	// Run this every five years
	return marked, time.Now().Sub(time.Unix(rec.GetMetadata().GetLastValidate(), 0)) < time.Hour*24*365*5
}

func (fs *fullScheme) name() string {
	return "full_validate_correct"
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

type twScheme struct{}

func (tw *twScheme) filter(rec *rcpb.Record) (bool, bool) {
	//Is it a cd?, doess it have a width?
	return rec.GetRelease().GetFolderId() == 242017, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (tw *twScheme) name() string {
	return "twelve_width"
}
