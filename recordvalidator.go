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
	completionDate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordvalidator_completion_date",
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

	return s
}

func (s *Server) updateMetrics(schemes *pb.Schemes) {
	time.Sleep(time.Second * 2)
	for _, sc := range schemes.GetSchemes() {
		if len(sc.GetCompletedIds()) > 0 {
			s.Log(fmt.Sprintf("%v from %v for %v (e.g. %v) [%v]", len(sc.GetCompletedIds()), len(sc.GetInstanceIds()), sc.GetName(), sc.GetCompletedIds()[0], sc.GetCurrentPick()))
		} else {
			s.Log(fmt.Sprintf("No completes for %v -> %v", sc.GetName(), len(sc.GetInstanceIds())))
		}
		prop := float64(len(sc.GetCompletedIds())) / float64(len(sc.GetInstanceIds())+len(sc.GetCompletedIds()))
		dur := time.Now().Sub(time.Unix(sc.GetStartTime(), 0)).Seconds()
		extraDur := dur/prop - dur
		finishTime := time.Now().Add(time.Second * time.Duration(extraDur)).Unix()

		doneCount.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(len(sc.GetCompletedIds())))
		completion.With(prometheus.Labels{"scheme": sc.GetName()}).Set(prop)
		completionDate.With(prometheus.Labels{"scheme": sc.GetName()}).Set(float64(finishTime))
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
	s.updateMetrics(schemes)
	return schemes, nil
}

func (s *Server) save(ctx context.Context, schemes *pb.Schemes) error {
	if s.failSave {
		return fmt.Errorf("Built to fail")
	}
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
	r, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_All{true}})
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

func (s *Server) update(ctx context.Context, iid int32) error {
	if s.test {
		return nil
	}
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateRecord(ctx, &rcpb.UpdateRecordRequest{
		Reason: "Update for validation",
		Update: &rcpb.Record{
			Release: &gdpb.Release{
				InstanceId: iid,
			},
			Metadata: &rcpb.ReleaseMetadata{Category: rcpb.ReleaseMetadata_PRE_VALIDATE},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	rcpb.RegisterClientUpdateServiceServer(server, s)
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
	return []*pbg.State{
		&pbg.State{Key: "magic", Value: int64(12345)},
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
	ctx, cancel := utils.ManualContext("rvsu", "rvsu", time.Minute, false)
	if _, err := server.load(ctx); err != nil {
		server.Log(fmt.Sprintf("Unable to load: %v", err))
		time.Sleep(time.Second * 5)
		return
	}
	cancel()

	fmt.Printf("%v", server.Serve())
}
