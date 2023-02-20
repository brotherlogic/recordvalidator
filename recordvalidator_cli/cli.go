package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordvalidator/proto"
)

func main() {
	ctx, cancel := utils.ManualContext("recordvalidator_cli", time.Minute*30)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "recordvalidator")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "ping":
		id, _ := strconv.ParseInt(os.Args[2], 10, 32)
		sclient := pbrc.NewClientUpdateServiceClient(conn)
		_, err = sclient.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(id)})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	case "which":
		id, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			log.Fatalf("Bad parse: %v", err)
		}
		sclient := pb.NewRecordValidatorServiceClient(conn)
		scheme, err := sclient.GetScheme(ctx, &pb.GetSchemeRequest{InstanceId: int32(id)})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("Scheme: %v [%v]\n", scheme.GetScheme().GetName(), scheme.GetScheme().GetActive())
		fmt.Printf("Pick: %v\n", scheme.Scheme.GetCurrentPick())
		for id, date := range scheme.Scheme.GetCompleteDate() {
			if time.Since(time.Unix(date, 0)) < time.Hour*24 {
				fmt.Printf("%v was recorded in the last 24 hours (%v) \n", id, time.Unix(date, 0))
			}
		}
	case "get":
		sclient := pb.NewRecordValidatorServiceClient(conn)
		scheme, err := sclient.GetScheme(ctx, &pb.GetSchemeRequest{Name: os.Args[2]})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("Scheme: %v\n", scheme)
		fmt.Printf("Pick: %v\n", scheme.Scheme.GetCurrentPick())
		for id, date := range scheme.Scheme.GetCompleteDate() {
			if time.Since(time.Unix(date, 0)) < time.Hour*24 {
				fmt.Printf("%v was recorded in the last 24 hours (%v) \n", id, time.Unix(date, 0))
			}
		}
	case "fullping":
		ctx2, cancel2 := utils.ManualContext("recordcollectioncli-"+os.Args[1], time.Hour)
		defer cancel2()

		conn2, err := utils.LFDialServer(ctx2, "recordcollection")
		if err != nil {
			log.Fatalf("Cannot reach rc: %v", err)
		}
		defer conn2.Close()

		registry := pbrc.NewRecordCollectionServiceClient(conn2)
		ids, err := registry.QueryRecords(ctx2, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_All{true}})
		if err != nil {
			log.Fatalf("Bad query: %v", err)
		}

		sclient := pbrc.NewClientUpdateServiceClient(conn)
		for i, id := range ids.GetInstanceIds() {
			fmt.Printf("PING %v -> %v\n", i, id)
			ctx, cancel = utils.ManualContext("recordvalidator_cli", time.Minute*30)
			_, err = sclient.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(id)})
			cancel()
			if err != nil {
				log.Fatalf("Error on GET: %v", err)
			}
		}
	}
}
