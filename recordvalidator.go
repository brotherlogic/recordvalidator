package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordvalidator/proto"
)

const (
	// SCHEMES - Where we store schemes
	SCHEMES = "/github.com/brotherlogic/recordvalidator/schemes"
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

	return s
}

func (s *Server) load(ctx context.Context) (*pb.Schemes, error) {
	if s.failLoad {
		return nil, fmt.Errorf("Bad load")
	}
	data, _, err := s.KSclient.Read(ctx, SCHEMES, &pb.Schemes{})
	if err != nil {
		return nil, err
	}
	schemes := data.(*pb.Schemes)
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
			return &rcpb.Record{Metadata: &rcpb.ReleaseMetadata{GoalFolder: 466902}}, nil
		}
		if iid == 2 {
			return &rcpb.Record{Metadata: &rcpb.ReleaseMetadata{GoalFolder: 466902, RecordWidth: 1.2}}, nil
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

	fmt.Printf("%v", server.Serve())
}
