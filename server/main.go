package main

import (
	"context"
	prsn "grpc/person"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
)

const (
	port = ":3333"
)

type Person struct {
	savedPersons []*prsn.PersonRequest
}

func (p *Person) CreatePerson(c context.Context, input *prsn.PersonRequest) (*prsn.PersonResponse, error) {
	p.savedPersons = append(p.savedPersons, input)
	return &prsn.PersonResponse{Id: input.Id, Success: true}, nil
}

func (p *Person) GetPerson(fltr *prsn.PersonFilter, stream prsn.Person_GetPersonServer) error {

	for _, person := range p.savedPersons {
		if fltr.Keyword != "" {
			if !strings.Contains(person.Name, fltr.Keyword) {
				continue
			}
		}
		err := stream.Send(person)
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen port %v, cause: %v", port, err)
		return
	}
	s := grpc.NewServer()
	prsn.RegisterPersonServer(s, &Person{})
	s.Serve(lis)
}
