package main

import pbrc "github.com/brotherlogic/recordcollection/proto"

type schemeGenerator interface {
	// Does this record apply to the this filter
	filter(rec *pbrc.Record) bool
	name() string
}
