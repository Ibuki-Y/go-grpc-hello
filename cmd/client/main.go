package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	hellopb "mygrpc/pkg/grpc"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC Client")

	//  標準入力から文字列を受け取るスキャナ
	scanner = bufio.NewScanner(os.Stdin)

	// gRPCサーバーとのコネクションを確立
	address := os.Getenv("LOCAL_PORT")
	conn, err := grpc.Dial(
		address,
		grpc.WithUnaryInterceptor(unaryClientInterceptor),
		grpc.WithStreamInterceptor(streamClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // コネクションでSSL/TLSを使用しない
		grpc.WithBlock(), // コネクションが確立されるまで待機
	)
	if err != nil {
		log.Fatalln(err)
		log.Fatal("Connection failed")
		return
	}
	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	fmt.Println("come here")
	for {
		fmt.Println("1: send Request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStreams")
		fmt.Println("5: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()

		case "2":
			HelloServerStream()

		case "3":
			HelloClientStream()

		case "4":
			HelloBiStreams()

		case "5":
			fmt.Println("bye")
			goto M
		}
	}
M:
}

func Hello() {
	fmt.Println("Please enter your name")
	scanner.Scan()
	name := scanner.Text()

	// リクエストに使うHelloRequest型の生成
	req := &hellopb.HelloRequest{
		Name: name,
	}
	// Helloメソッドの実行 -> HelloResponse型のレスポンスを入手
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %s\n", stat.Code())
			fmt.Printf("message: %s\n", stat.Message())
			fmt.Printf("details: %s\n", stat.Details())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Println("Please enter your name")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received")
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	cnt := 5
	fmt.Printf("Please enter %d names\n", cnt)
	for i := 0; i < cnt; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStreams() {
	stream, err := client.HelloBiStreams(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	cnt := 5
	fmt.Printf("Please enter %d names\n", cnt)

	var sendEnd, recvEnd bool
	sendCnt := 0
	for !(sendEnd && recvEnd) {
		// 送信処理
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCnt++
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Println(err)
				sendEnd = true
			}

			if sendCnt == cnt {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
				}
			}
		}

		// 受信処理
		if !recvEnd {
			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				recvEnd = true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}
}
