package main

import (
	"strings"
	"time"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

type schemeGenerator interface {
	// Does this record apply to the this filter and what's the order
	filter(rec *rcpb.Record) (bool, bool, float32)
	name() string
}

type allTwelves struct {
}

func (fs *allTwelves) filter(rec *rcpb.Record) (bool, bool, float32) {
	// Can't process these
	if rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_PARENTS {
		return false, true, -1
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
	return marked,
		rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (fs *allTwelves) name() string {
	return "twelves"
}

type fullScheme struct{}

func (fs *fullScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	// Can't process these
	if rec.Metadata.GetCategory() == rcpb.ReleaseMetadata_PARENTS {
		return false, true, -1
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
	return marked,
		time.Now().Sub(time.Unix(rec.GetMetadata().GetLastValidate(), 0)) < time.Hour*24*365*5 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (fs *fullScheme) name() string {
	return "full_validate_correct"
}

type keeperScheme struct{}

func (ks *keeperScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	//Is it a keeper?, doess it have a width?
	return rec.GetMetadata().GetGoalFolder() == 2259637,
		rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (ks *keeperScheme) name() string {
	return "keeper_width"
}

type ageScheme struct{}

func (as *ageScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	// Sold Digital recordings should be included ehre
	marked := false
	if rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_LISTED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_STAGED_TO_SELL &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
		rec.Metadata.GetCategory() != rcpb.ReleaseMetadata_PARENTS {
		marked = true
	}

	// Listen to everything every five years
	return marked,
		time.Now().Sub(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) < time.Hour*24*365*2 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (as *ageScheme) name() string {
	return "age"
}

type fallScheme struct{}

func (fs *fallScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	//Is it a keeper?, doess it have a width?
	return rec.GetMetadata().GetGoalFolder() == 716318,
		rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (fs *fallScheme) name() string {
	return "fall_width"
}

type nsSleeve struct{}

func (nss *nsSleeve) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetRelease().GetFolderId() == 242017,
		rec.GetMetadata().GetSleeve() != rcpb.ReleaseMetadata_SLEEVE_UNKNOWN && rec.GetMetadata().GetSleeve() != rcpb.ReleaseMetadata_VINYL_STORAGE_NO_INNER, -1
}

func (nss *nsSleeve) name() string {
	return "twelve_inch_sleeves"
}

type tapeProc struct{}

func (_ *tapeProc) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_TAPE,
		rec.GetMetadata().GetCdPath() != "" && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE, -1
}

func (_ *tapeProc) name() string {
	return "tape_processing"
}

type nsSevenSleeve struct{}

func (nss *nsSevenSleeve) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_7_INCH && rec.GetRelease().GetFolderId() == 267116,
		rec.GetMetadata().GetSleeve() != rcpb.ReleaseMetadata_SLEEVE_UNKNOWN && rec.GetMetadata().GetSleeve() != rcpb.ReleaseMetadata_VINYL_STORAGE_NO_INNER, -1
}

func (nss *nsSevenSleeve) name() string {
	return "seven_inch_sleeves"
}

type cdScheme struct{}

func (cds *cdScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	//Is it a cd?, doess it have a width?
	return rec.GetRelease().GetFolderId() == 242018,
		rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (cds *cdScheme) name() string {
	return "cd_width"
}

type twScheme struct{}

func (tw *twScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	//Is it a cd?, doess it have a width?
	return rec.GetRelease().GetFolderId() == 242017,
		rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (tw *twScheme) name() string {
	return "twelve_width"
}

type tenScheme struct{}

func (tw *tenScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	isTen := false
	for _, format := range rec.GetRelease().GetFormats() {
		if strings.Contains(format.GetName(), "10") {
			isTen = true
		}
		for _, desc := range format.GetDescriptions() {
			if strings.Contains(desc, "10") {
				isTen = true
			}
		}
	}
	//Is it a cd?, doess it have a width?
	return isTen,
		((rec.GetMetadata().GetSaleId() > 0 || rec.GetMetadata().GetSoldDate() > 0) || rec.GetMetadata().GetLastValidate() > 0) && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (tw *tenScheme) name() string {
	return "ten_inches"
}

type libScheme struct{}

func (ls *libScheme) filter(rec *rcpb.Record) (bool, bool, float32) {
	//Is it a cd?, doess it have a width?
	return rec.GetMetadata().GetGoalFolder() == 882359 ||
			rec.GetMetadata().GetGoalFolder() == 1281012 ||
			rec.GetMetadata().GetGoalFolder() == 857451 ||
			rec.GetMetadata().GetGoalFolder() == 1409151 ||
			rec.GetMetadata().GetGoalFolder() == 823501 ||
			rec.GetMetadata().GetGoalFolder() == 842724 ||
			rec.GetMetadata().GetGoalFolder() == 681783 ||
			rec.GetMetadata().GetGoalFolder() == 529723 ||
			rec.GetMetadata().GetGoalFolder() == 1642995, rec.GetMetadata().GetRecordWidth() > 0 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

func (ls *libScheme) name() string {
	return "library_width"
}

type older struct{}

func (os *older) name() string {
	return "old_age"
}

func (os *older) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type olderTwelves struct{}

func (os *olderTwelves) name() string {
	return "old_age_twelves"
}

func (os *olderTwelves) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type olderSevens struct{}

func (os *olderSevens) name() string {
	return "old_age_sevens"
}

func (os *olderSevens) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 267116,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newer struct{}

func (os *newer) name() string {
	return "new_age"
}

func (os *newer) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newerTwelves struct{}

func (*newerTwelves) name() string {
	return "new_age_twelves"
}

func (*newerTwelves) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newerSevens struct{}

func (*newerSevens) name() string {
	return "new_age_sevens"
}

func (*newerSevens) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 267116, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type bad_ones struct{}

func (*bad_ones) name() string {
	return "bad_ones"
}

func (*bad_ones) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}

type bad_ones_twelve struct{}

func (*bad_ones_twelve) name() string {
	return "bad_ones_twelves"
}

func (*bad_ones_twelve) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}
