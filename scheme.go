package main

import (
	"math/rand"
	"strings"
	"time"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func rightFormatQuantity(r *rcpb.Record) bool {
	cd := false
	for _, format := range r.GetRelease().GetFormats() {
		if format.GetName() == "CD" {
			cd = true
		}
	}

	if !cd {
		return false
	}

	/*if r.GetRelease().GetFormatQuantity() == 1 {
		return true
	}*/

	return true
}

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
		marked = false
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
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_LISTED_TO_SELL,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type olderND struct{}

func (os *olderND) name() string {
	return "old_age_no_digital"
}

func (os *olderND) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type olderNDS struct{}

func (os *olderNDS) name() string {
	return "old_age_no_digital_singles"
}

func (os *olderNDS) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_DIGITAL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type olderNDSF struct{}

func (os *olderNDSF) name() string {
	return "old_age_no_digital_singles_filable"
}

func (os *olderNDSF) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetGoalFolder() == 242017 &&
			rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_DIGITAL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
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
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_LISTED_TO_SELL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newerND struct{}

func (os *newerND) name() string {
	return "new_age_no_digital"
}

func (os *newerND) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newerNDS struct{}

func (os *newerNDS) name() string {
	return "new_age_no_digital_singles"
}

func (os *newerNDS) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_DIGITAL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		-1
}

type newerNDSF struct{}

func (os *newerNDSF) name() string {
	return "new_age_no_digital_singles_filable"
}

func (os *newerNDSF) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) &&
			rec.GetMetadata().GetGoalFolder() == 242017 &&
			rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_DIGITAL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
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
		rec.GetMetadata().GetOverallScore()
}

type bad_ones_twelve_single struct{}

func (*bad_ones_twelve_single) name() string {
	return "bad_ones_twelves_single"
}

func (*bad_ones_twelve_single) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rightFormatQuantity(rec) && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.GetMetadata().GetOverallScore()
}

type random_twelves_single struct{}

func (*random_twelves_single) name() string {
	return "random_twelves_single"
}

func (*random_twelves_single) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rand.Float32()
}

type random_twelves_single_v2 struct{}

func (*random_twelves_single_v2) name() string {
	return "random_twelves_single_v2"
}

func (*random_twelves_single_v2) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rightFormatQuantity(rec) && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rand.Float32()
}

type boxsets struct{}

func (*boxsets) name() string {
	return "boxsets"
}

func (*boxsets) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rec.GetRelease().GetFormatQuantity() > 4 && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE,
		rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rand.Float32()
}

type sonimage struct{}

func (*sonimage) name() string {
	return "sonimage"
}

func (*sonimage) filter(rec *rcpb.Record) (bool, bool, float32) {
	found := false
	for _, label := range rec.GetRelease().GetLabels() {
		if label.GetName() == "Sonimage" {
			found = true
		}
	}

	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_STAGED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			found, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.GetMetadata().GetOverallScore()
}

type piecelock struct{}

func (*piecelock) name() string {
	return "piecelock"
}

func (*piecelock) filter(rec *rcpb.Record) (bool, bool, float32) {
	found := false
	for _, label := range rec.GetRelease().GetLabels() {
		if label.GetName() == "Piecelock 70" {
			found = true
		}
	}

	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_STAGED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			found, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}

type april struct{}

func (*april) name() string {
	return "april"
}

func (*april) filter(rec *rcpb.Record) (bool, bool, float32) {
	found := false
	if strings.Contains(rec.GetRelease().GetTitle(), "April") {
		found = true
	}

	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_STAGED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			found, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}

type hudson struct{}

func (*hudson) name() string {
	return "hudson"
}

func (*hudson) filter(rec *rcpb.Record) (bool, bool, float32) {
	found := false
	for _, label := range rec.GetRelease().GetLabels() {
		if strings.Contains(label.GetName(), "Hudson") {
			found = true
		}
	}

	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_STAGED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			found, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}

type fall struct{}

func (*fall) name() string {
	return "fall"
}

func (*fall) filter(rec *rcpb.Record) (bool, bool, float32) {
	found := false
	for _, artist := range rec.GetRelease().GetArtists() {
		if strings.Contains(artist.GetName(), "Fall") {
			found = true
		}
	}

	return rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_UNKNOWN &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_ARRIVED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PARENTS &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_HIGH_SCHOOL &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_STAGED &&
			rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_SOLD_ARCHIVE &&
			rec.GetMetadata().GetGoalFolder() == 242017 &&
			found, rec.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_UNKNOWN && rec.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_PRE_VALIDATE,
		rec.Metadata.GetOverallScore()
}

type oldest struct{}

func (*oldest) name() string {
	return "oldest"
}

func (*oldest) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION ||
			rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE,
		time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) < time.Hour*24*365*3,
		float32(rec.GetMetadata().GetLastListenTime())

}

type oldestSingle struct{}

func (*oldestSingle) name() string {
	return "oldest_single"
}

func (*oldestSingle) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetRelease().GetFormatQuantity() == 1 && (rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION ||
			rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE) && rec.GetMetadata().GetLastListenTime() != 0,
		time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) < time.Hour*24*365*3,
		float32(rec.GetMetadata().GetLastListenTime())
}

type fastDump struct{}

func (*fastDump) name() string {
	return "fast_dump"
}

func (*fastDump) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetGoalFolder() == 565206,
		rec.GetMetadata().GetWeightInGrams() > 0,
		float32(rec.GetMetadata().GetLastListenTime())

}

type keepers struct{}

func (*keepers) name() string {
	return "keepers"
}

func (*keepers) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH &&
			(rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION || rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE) &&
			(rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_BOX_UNKNOWN || rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_OUT_OF_BOX),
		!(time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) > time.Hour*24*365 && rec.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_KEEP_UNKNOWN),
		rand.Float32()
}

type keepers_single struct{}

func (*keepers_single) name() string {
	return "keepers_single"
}

func (*keepers_single) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetRelease().GetFormatQuantity() == 1 &&
			(rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION || rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE) &&
			(rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_BOX_UNKNOWN || rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_OUT_OF_BOX),
		!(time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) > time.Hour*24*365 && rec.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_KEEP_UNKNOWN),
		rand.Float32()
}

type was_parents struct{}

func (*was_parents) name() string {
	return "was_parents"
}

func (*was_parents) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetMetadata().GetWasParents(), rec.GetMetadata().GetLastListenTime() > 0,
		float32(rec.GetRelease().GetInstanceId())
}

type was_parents_rev struct{}

func (*was_parents_rev) name() string {
	return "was_parents_rev"
}

func (*was_parents_rev) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetMetadata().GetWasParents() && rec.GetRelease().GetFormatQuantity() != 1, rec.GetMetadata().GetLastListenTime() > 0,
		float32(rec.GetRelease().GetInstanceId())
}

type was_parents_single struct{}

func (*was_parents_single) name() string {
	return "was_parents_single"
}

func (*was_parents_single) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetMetadata().GetWasParents() && rec.GetRelease().GetFormatQuantity() == 1, rec.GetMetadata().GetLastListenTime() > 0,
		float32(rec.GetRelease().GetInstanceId())
}

type was_parents_single_2 struct{}

func (*was_parents_single_2) name() string {
	return "was_parents_single_2"
}

func (*was_parents_single_2) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetMetadata().GetWasParents() && rec.GetRelease().GetFormatQuantity() == 1, rec.GetMetadata().GetLastListenTime() > 0,
		float32(rec.GetRelease().GetInstanceId())
}

type full_parents struct{}

func (*full_parents) name() string {
	return "full_parents"
}

func (*full_parents) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH && rec.GetRelease().GetFolderId() == 1727264 || rec.GetRelease().GetFolderId() == 6268933, rec.GetMetadata().GetLastListenTime() > 0,
		float32(rec.GetRelease().GetInstanceId())
}

type keepersSeven struct{}

func (*keepersSeven) name() string {
	return "keepers_seven"
}

func (*keepersSeven) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_7_INCH &&
			(rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION || rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE) &&
			(rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_BOX_UNKNOWN || rec.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_OUT_OF_BOX),
		!(time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) > time.Hour*24*365 && rec.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_KEEP_UNKNOWN),
		rand.Float32()
}

type oldTwelve struct{}

func (*oldTwelve) name() string {
	return "old_twelves"
}

func (*oldTwelve) filter(rec *rcpb.Record) (bool, bool, float32) {
	return rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH &&
			(rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION || rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE),
		time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) > time.Hour*24*365,
		float32(rec.GetMetadata().GetLastListenTime())
}

type oldFall struct{}

func (*oldFall) name() string {
	return "olf_fall"
}

func (*oldFall) filter(rec *rcpb.Record) (bool, bool, float32) {
	rand.Seed(time.Now().UnixNano())
	foundFall := false
	for _, artist := range rec.GetRelease().GetArtists() {
		if artist.GetName() == "The Fall" {
			foundFall = true
		}
	}

	return foundFall && rec.GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_12_INCH &&
			(rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_IN_COLLECTION || rec.GetMetadata().GetCategory() == rcpb.ReleaseMetadata_PRE_VALIDATE),
		time.Since(time.Unix(rec.GetMetadata().GetLastListenTime(), 0)) < time.Hour*24*365,
		rand.Float32()
}
