package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"google.golang.org/grpc"

	"github.com/flan6/microservice-messaging/internal/channel"
	"github.com/flan6/microservice-messaging/internal/consumer/tasks"
	proto "github.com/flan6/microservice-messaging/internal/generator/api/rpc/pb"
	"github.com/flan6/microservice-messaging/internal/generator/api/rpc/server"
)

const redisAddr = "127.0.0.1:6379" // TODO: get redis addrs from config yaml

func main() {
	sigint := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
	}()

	go runApi()
	go runWorker()

	<-sigint
	log.Println("shut down")
}

func runApi() { // TODO: add context for graceful stop
	// TODO: use grpc-gateway for http aswell
	api := server.NewApi(
		channel.NewChannel(
			asynq.NewClient(
				asynq.RedisClientOpt{Addr: redisAddr},
			),
		),
	)
	grpcServer := grpc.NewServer()

	proto.RegisterApiServer(grpcServer, api)

	listener, err := net.Listen("tcp", ":9100") // TODO: get port from config yaml
	if err != nil {
		log.Fatalf("could not listen to port:  %v", err)
	}

	log.Println("running rpc server")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("could not run rpc server: %v", err)
	}
}

func runWorker() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.Handle(tasks.TypeEmail, tasks.NewSendEmailProcessor())

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
