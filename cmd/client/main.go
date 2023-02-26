package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	hellopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		fmt.Println("2: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()

		case "2":
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
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}
