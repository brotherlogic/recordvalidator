protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/recordvalidator.proto
mv proto/github.com/brotherlogic/recordvalidator/proto/* ./proto
