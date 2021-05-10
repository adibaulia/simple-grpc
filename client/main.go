package main

import (
	"context"
	"fmt"
	"io"

	prsn "grpc/person"

	"google.golang.org/grpc"
)

const (
	address = "localhost:3333"
)

// createPerson calls the RPC method CreatePerson of PersonServer
func createPerson(client prsn.PersonClient, person *prsn.PersonRequest) {
	resp, err := client.CreatePerson(context.Background(), person)
	if err != nil {
		fmt.Println("Could not create Person: ", err)
		return
	}
	if resp.Success {
		fmt.Println("A new Person has been added with id: ", resp.Id)
	}
}

// getPersons calls the RPC method GetPersons of PersonServer
func getPersons(client prsn.PersonClient, filter *prsn.PersonFilter) {
	// calling the streaming API
	stream, err := client.GetPerson(context.Background(), filter)
	if err != nil {
		fmt.Println("Error on get persons: ", err)
		return
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("%v.GetPersons(_) = _, %v", client, err)
		}
		fmt.Println("Person: ", person)
	}
}
func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect: ", err)
		return
	}
	defer conn.Close()
	// Creates a new PersonClient
	client := prsn.NewPersonClient(conn)
	person := &prsn.PersonRequest{
		Id:    1001,
		Name:  "Reddy",
		Email: "reddy@xyz.com",
		Phone: "9898982929",
		Address: []*prsn.PersonRequest_Address{
			{
				Street:            "Tripilcane",
				City:              "Chennai",
				State:             "TN",
				Zip:               "600002",
				IsShippingAddress: false,
			},
			{
				Street:            "Balaji colony",
				City:              "Tirupati",
				State:             "AP",
				Zip:               "517501",
				IsShippingAddress: true,
			},
		},
	}
	// Create a new person
	createPerson(client, person)
	person = &prsn.PersonRequest{
		Id:    1002,
		Name:  "Raj",
		Email: "raj@xyz.com",
		Phone: "5000510001",
		Address: []*prsn.PersonRequest_Address{
			{
				Street:            "Marathahalli",
				City:              "Bangalore",
				State:             "KS",
				Zip:               "560037",
				IsShippingAddress: true,
			},
		},
	}
	// Create a new person
	createPerson(client, person)
	// Filter with an empty Keyword
	filter := &prsn.PersonFilter{Keyword: ""}
	getPersons(client, filter)
}
