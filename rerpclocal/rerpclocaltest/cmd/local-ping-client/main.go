package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pingpb "github.com/rerpc/rerpc/internal/ping/v1test"
	"github.com/rerpc/rerpc/rerpclocal"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pingpb.NewPingServiceClientReRPCV2(rerpclocal.NewClient("local-ping"))

	// /Ping
	pingRequest := &pingpb.PingRequest{Number: 42, Msg: "message"}
	pingResponse, err := client.Ping(context.Background(), pingRequest)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Failed to ping: %v\n", err))
		os.Exit(1)
	}

	_, _ = os.Stdout.WriteString(fmt.Sprintf("Got number (%d), message %q\n", pingResponse.Number, pingResponse.Msg))

	// /Sum
	sumStream := client.Sum(ctx)
	for i := 1; i <= 10; i++ {
		number := int64(i)
		fmt.Println("Sending", number)
		if err := sumStream.Send(&pingpb.SumRequest{Number: number}); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Failed to sum: %v\n", err))
			os.Exit(1)
		}
	}
	sumResponse, err := sumStream.CloseAndReceive()
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Failed to close and receive sum: %v\n", err))
		os.Exit(1)
	}

	_, _ = os.Stdout.WriteString(fmt.Sprintf("Got sum (%d)\n", sumResponse.Sum))

	return
}
