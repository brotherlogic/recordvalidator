package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordvalidator/proto"
)

const (
	// SCHEMES - Where we store schemes
	SCHEMES = "/github.com/brotherlogic/recordvalidator/schemes"
)

var (
	completion = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_completion",
		Help: "The size of the print queue",
	}, []string{"scheme"})
	doneCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_complete",
		Help: "The size of the print queue",
	}, []string{"scheme"})
	toGo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_togo",
		Help: "The size of the print queue",
	}, []string{"scheme"})
	completionDate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_completion_date",
		Help: "The size of the print queue",
	}, []string{"scheme"})
	completionDateV2 = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_completion_date_v2",
		Help: "The size of the print queue",
	}, []string{"scheme"})
	perDay = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_completion_per_day",
		Help: "The size of the print queue",
	}, []string{"scheme"})
)

//Server main server type
type Server struct {
	*goserver.GoServer
	test           bool
	sgs            []schemeGenerator
	failLoad       bool
	failRecordLoad bool
	failSave       bool
	failLoadAll    bool
}

// Init builds the server
func Init() *Server {

	s := &Server{
		GoServer: &goserver.GoServer{},
		sgs:      []schemeGenerator{},
	}

	s.sgs = append(s.sgs, &keeperScheme{})
	s.sgs = append(s.sgs, &cdScheme{})
	s.sgs = append(s.sgs, &twScheme{})
	s.sgs = append(s.sgs, &fullScheme{})
	s.sgs = append(s.sgs, &allTwelves{})
	s.sgs = append(s.sgs, &fallScheme{})
	s.sgs = append(s.sgs, &ageScheme{})
	s.sgs = append(s.sgs, &tenScheme{})
	s.sgs = append(s.sgs, &libScheme{})
	s.sgs = append(s.sgs, &older{})
	s.sgs = append(s.sgs, &newer{})
	s.sgs = append(s.sgs, &nsSleeve{})
	s.sgs = append(s.sgs, &nsSevenSleeve{})
	s.sgs = append(s.sgs, &bad_ones{})
	s.sgs = append(s.sgs, &olderTwelves{})
	s.sgs = append(s.sgs, &tapeProc{})

	return s
}

func (s *Server) updateMetrics(schemes *pb.Schemes) {
	for _, sc := range schemes.GetSchemes() {
		prop := float64(len(sc.GetCompletedIds())) / float64(len(sc.GetInstanceIds())+len(sc.GetCompletedIds()))
		dur := time.Since(time.Unix(sc.GetStartTime(), 0)).Seconds()
		extraDur := dur/prop - dur
		finishTime := time.Now().Add(time.Second * time.Duration(extraDur)).Unix()

		doneCount.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(len(sc.GetCompletedIds())))
		toGo.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(len(sc.GetInstanceIds())))
		completion.With(prometheus.Labels{"scheme": sc.GetName()}).Set(prop)
		completionDate.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(finishTime))

		last14days := float64(0)
		for _, date := range sc.GetCompleteDate() {
			if time.Since(time.Unix(date, 0)) < time.Hour*24*14 {
				last14days++
			}
		}

		compPerDay := last14days / 14
		togo := float64(len(sc.GetInstanceIds()))
		days := togo / compPerDay
		ftime := time.Now().Add(time.Hour * time.Duration(24*days))
		s.Log(fmt.Sprintf("COMP %v: %v, %v, %v, %v", sc.GetName(), compPerDay, togo, days, ftime))
		perDay.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(compPerDay))
		completionDateV2.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(ftime.Unix()))
	}
}

func (s *Server) load(ctx context.Context) (*pb.Schemes, error) {
	if s.failLoad {
		return nil, fmt.Errorf("Bad load")
	}
	data, _, err := s.KSclient.Read(ctx, SCHEMES, &pb.Schemes{})

	if err != nil && status.Convert(err).Code() != codes.InvalidArgument {
		return nil, err
	}
	if err != nil {
		data = &pb.Schemes{}
	}
	schemes := data.(*pb.Schemes)

	mapper := make(map[string]*pb.Scheme)
	for _, sc := range schemes.GetSchemes() {
		mapper[sc.GetName()] = sc
	}
	schemes.Schemes = make([]*pb.Scheme, 0)
	for _, sc := range mapper {
		schemes.Schemes = append(schemes.Schemes, sc)
	}

	s.updateMetrics(schemes)

	for _, scheme := range schemes.GetSchemes() {
		if scheme.CompleteDate == nil {
			scheme.CompleteDate = make(map[int32]int64)
		}

	}

	return schemes, nil
}

func (s *Server) save(ctx context.Context, schemes *pb.Schemes) error {
	if s.failSave {
		return fmt.Errorf("Built to fail")
	}

	for _, scheme := range schemes.GetSchemes() {
		var nums []int32
		nmap := make(map[int32]bool)
		for _, num := range scheme.GetCompletedIds() {
			if num > 0 {
				nmap[num] = true
			}
		}
		for v := range nmap {
			nums = append(nums, v)
		}
		scheme.CompletedIds = nums

		/*if scheme.GetName() == "twelve_inch_sleeves" {
			scheme.CompletedIds = make([]int32, 0)
			scheme.CompleteDate = make(map[int32]int64)
			s.Log(fmt.Sprintf("Reset"))
		}*/
	}

	s.updateMetrics(schemes)

	return s.KSclient.Save(ctx, SCHEMES, schemes)
}

func (s *Server) loadRecord(ctx context.Context, iid int32) (*rcpb.Record, error) {
	if s.failRecordLoad {
		return nil, fmt.Errorf("Built too fail")
	}
	if s.test {
		if iid == 1 {
			return &rcpb.Record{Metadata: &rcpb.ReleaseMetadata{GoalFolder: 2259637}}, nil
		}
		if iid == 2 {
			return &rcpb.Record{Metadata: &rcpb.ReleaseMetadata{GoalFolder: 2259637, RecordWidth: 1.2}}, nil
		}
		return &rcpb.Record{}, nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	r, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: iid})
	if err != nil {
		return nil, err
	}
	return r.GetRecord(), nil
}

func (s *Server) getAllRecords(ctx context.Context) ([]int32, error) {
	if s.failLoadAll {
		return nil, fmt.Errorf("Built to fail")
	}
	if s.test {
		return []int32{1, 2, 3}, nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	r, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_FolderId{267116}})
	if err != nil {
		return nil, err
	}
	return r.GetInstanceIds(), nil
}

func (s *Server) getRecord(ctx context.Context, iid int32) (*rcpb.Record, error) {
	if s.failLoadAll {
		return nil, fmt.Errorf("Built to fail")
	}
	if s.test {
		return &rcpb.Record{}, nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	r, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: iid})
	if err != nil {
		return nil, err
	}
	return r.GetRecord(), nil
}

func (s *Server) update(ctx context.Context, iid int32, soft, unbox bool) error {
	s.Log(fmt.Sprintf("Updating for %v", iid))
	if s.test {
		return nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	req := &rcpb.UpdateRecordRequest{
		Reason: "Update for validation",
		Update: &rcpb.Record{
			Release: &gdpb.Release{
				InstanceId: iid,
			},
			Metadata: &rcpb.ReleaseMetadata{
				Category: rcpb.ReleaseMetadata_PRE_VALIDATE,
				Dirty:    true},
		}}

	if unbox {
		req.Update.Metadata.NewBoxState = rcpb.ReleaseMetadata_OUT_OF_BOX
	}

	if soft {
		req.Update.Metadata.Category = rcpb.ReleaseMetadata_PRE_SOFT_VALIDATE
	}

	_, err = client.UpdateRecord(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) softValidate(ctx context.Context, iid int32) error {
	s.Log(fmt.Sprintf("Updating for %v", iid))
	if s.test {
		return nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	req := &rcpb.UpdateRecordRequest{
		Reason: "Update for validation",
		Update: &rcpb.Record{
			Release: &gdpb.Release{
				InstanceId: iid,
			},
			Metadata: &rcpb.ReleaseMetadata{
				Category: rcpb.ReleaseMetadata_SOFT_VALIDATED,
				Dirty:    true},
		}}

	_, err = client.UpdateRecord(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	rcpb.RegisterClientUpdateServiceServer(server, s)
	pb.RegisterRecordValidatorServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	found := false
	for _, scheme := range s.sgs {
		if scheme.name() == "twelve_inch_sell" {
			found = true
		}
	}
	return []*pbg.State{
		&pbg.State{Key: "magic", Text: fmt.Sprintf("%v", found)},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("recordvalidator", false, true)
	if err != nil {
		return
	}

	//Do a load to prepopulate metrics
	ctx, cancel := utils.ManualContext("rvsu", time.Minute)
	if _, err := server.load(ctx); err != nil {
		server.Log(fmt.Sprintf("Unable to load: %v", err))
		time.Sleep(time.Second * 5)
		return
	}
	cancel()

	fmt.Printf("%v", server.Serve())
}
