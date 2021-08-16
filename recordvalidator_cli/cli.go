package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pbrc "github.com/brotherlogic/recordcollection/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
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
		id, err := strconv.Atoi(os.Args[2])
		sclient := pbrc.NewClientUpdateServiceClient(conn)
		_, err = sclient.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(id)})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	}
}
