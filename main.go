package main

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/ruandao/micro-shippy-consignment-service-ser/consignmentMongo"
	pb "github.com/ruandao/micro-shippy-consignment-service-ser/proto/consignment"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {

	var callOptions microclient.CallOption = func(options *microclient.CallOptions) {
		options.RequestTimeout = time.Second * 5
		options.DialTimeout = time.Second * 6
	}
	service := micro.NewService(
		micro.Name("go.micro.srv.consignment.cli-create"),
		)
	service.Init()
	client := pb.NewShippingServiceClient(consignmentMongo.CONST_SER_NAME_CONSIGNMENT, service.Client())

	// Contact the server and print out its response.
	file := defaultFilename

	if len(os.Args) > 1 {
		file = os.Args[1]
	}


	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment, callOptions)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
