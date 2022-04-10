package main

import (
	"BookService/protobuff"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//creating pool of connections to db
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db = NewDB()
	defer db.Close()
	//initializing logger and servers
	log := logrus.New()
	gs := grpc.NewServer()
	ms := NewMyServer(log)
	protobuff.RegisterBookServiceServer(gs, ms)
	//reflection for inspecting purposes
	reflection.Register(gs)
	//starting server
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gs.Serve(l)

}
