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

type allTwelves struct {
}

func (fs *allTwelves) filter(rec *rcpb.Record) (bool, bool) {
	// Can't process these
	if rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_PARENTS {
		return false, true
	}

	marked := false
	if rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_LISTED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_STAGED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE {
		// This is in play
		if rec.Metadata.GetGoalFolder() == 242017 ||
			rec.GetMetadata().GetGoalFolder() == 1435521 ||
			rec.GetMetadata().GetGoalFolder() == 882359 ||
			rec.GetMetadata().GetGoalFolder() == 1799161 ||
			rec.GetMetadata().GetGoalFolder() == 1281012 ||
			rec.GetMetadata().GetGoalFolder() == 2268734 ||
			rec.GetMetadata().GetGoalFolder() == 2268731 ||
			rec.GetMetadata().GetGoalFolder() == 2021660 ||
			rec.GetMetadata().GetGoalFolder() == 857451 ||
			rec.GetMetadata().GetGoalFolder() == 2259637 ||
			rec.GetMetadata().GetGoalFolder() == 1799163 ||
			rec.GetMetadata().GetGoalFolder() == 1409151 ||
			rec.GetMetadata().GetGoalFolder() == 1191108 ||
			rec.GetMetadata().GetGoalFolder() == 823501 ||
			rec.GetMetadata().GetGoalFolder() == 529723 ||
			rec.GetMetadata().GetGoalFolder() == 2307240 ||
			rec.GetMetadata().GetGoalFolder() == 857449 ||
			rec.GetMetadata().GetGoalFolder() == 2268726 ||
			rec.GetMetadata().GetGoalFolder() == 681782 ||
			rec.GetMetadata().GetGoalFolder() == 1642995 ||
			rec.GetMetadata().GetGoalFolder() == 1456851 ||
			rec.GetMetadata().GetGoalFolder() == 716318 ||
			rec.GetMetadata().GetGoalFolder() == 842724 ||
			rec.GetMetadata().GetGoalFolder() == 1607992 ||
			rec.GetMetadata().GetGoalFolder() == 681783 ||
			rec.GetMetadata().GetGoalFolder() == 466902 {
			marked = true
		}
	}

	// Run this every five years
	return marked, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (fs *allTwelves) name() string {
	return "twelves"
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
	return marked, time.Now().Sub(time.Unix(rec.GetMetadata().GetLastValidate(), 0)) < time.Hour*24*365*5 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
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

type ageScheme struct{}

func (as *ageScheme) filter(rec *rcpb.Record) (bool, bool) {
	// Sold Digital recordings should be included ehre
	marked := false
	if rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_LISTED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_STAGED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_PARENTS {
		marked = true
	}

	// Listen to everything every five years
	return marked, time.Now().Sub(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) < time.Hour*24*365*2 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (as *ageScheme) name() string {
	return "age"
}

type fallScheme struct{}

func (fs *fallScheme) filter(rec *rcpb.Record) (bool, bool) {
	//Is it a keeper?, doess it have a width?
	return rec.GetMetadata().GetGoalFolder() == 716318, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE
}

func (fs *fallScheme) name() string {
	return "fall_width"
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
