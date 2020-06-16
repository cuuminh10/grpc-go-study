package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"maín/calculatorpb"
	"runtime"
	"time"
)



func main() {
	cc,err := grpc.Dial("localhost:50089",grpc.WithInsecure())

	if err != nil{
		log.Fatal("loi ket noi ",err)
	}



	client := calculatorpb.NewCalculatorServiceClient(cc)

	log.Println("day la client ",client)

	//callSum(client)
	//callPnd(client)
	//reqServer(client)
	// FindMax(client)

	Sqare(client)

	defer cc.Close()
}

func callSum(client calculatorpb.CalculatorServiceClient)  {
	resp,err := client.Sum(context.Background(),&calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})

	if err != nil{
		log.Println("loi jet noi")
	}

	log.Println("day la ket qua ",resp.Result)
}

func callPnd(client calculatorpb.CalculatorServiceClient)  {
	stream ,err := client.Pbd(context.Background(),&calculatorpb.PndRequest{
		Num: 120,
	})

	if err != nil{
		log.Println("loi jet noi")
	}

	for{
		resp,errRecv := stream.Recv()

		if (errRecv == io.EOF){
			log.Fatal("Cham dut ket noi")
		}

		log.Println("ket qua tra ve stream ",resp)
	}
}

func reqServer(client calculatorpb.CalculatorServiceClient){

	stream,err := client.Avg(context.Background())

	if err != nil{
		log.Println("loi jet noi")
	}

	arrAvg := []calculatorpb.AvgRequest{
		calculatorpb.AvgRequest{
			Num: 2,
		},
		calculatorpb.AvgRequest{
			Num: 3,
		},
		calculatorpb.AvgRequest{
			Num: 1,
		},
	}

	for _, request := range arrAvg {
		stream.Send(&request)
		if (err != nil){
			log.Fatal("loi ket noi 1",err)
		}
	}

	resp,errRec:=stream.CloseAndRecv()

	if (errRec != nil){
		log.Fatal("loi ket noi ",errRec)
	}

	log.Println("loi du k ",resp)

}

func FindMax(client calculatorpb.CalculatorServiceClient)  {

	stream,err := client.Max(context.Background())

	if err != nil{
		log.Println("loi jet noi")
	}

	cWait := make(chan  int)

	arrMax := []calculatorpb.MaxRequest{
		calculatorpb.MaxRequest{
			Num: 4,
		},
		calculatorpb.MaxRequest{
			Num: 5,
		},
		calculatorpb.MaxRequest{
			Num: 8,
		},
		calculatorpb.MaxRequest{
			Num: 1,
		},
	}

	go func() {
		for _, max := range arrMax {
			log.Println("client gui ",max)
			err := stream.Send(&max)
			if err != nil{
				log.Println("loi Ket noi")
			}
			time.Sleep(2*time.Second)
		}
		stream.CloseSend()
	}()

	fmt.Print(runtime.NumGoroutine())

	go func() {
		for{
			resp,err := stream.Recv()

			if err == io.EOF{
				log.Println("serve cham dut")
				break
			}

			if err != nil{
				log.Println("serve cham loi ",err)
				break
			}

			log.Println("gia tri max la ",resp.Result)

		}
		fmt.Println("day la ",len(cWait))
		close(cWait)
	}()

	fmt.Println("day la quet")
    <-cWait
}

func Sqare(client calculatorpb.CalculatorServiceClient)  {

	resp,err := client.Sqare(context.Background(),&calculatorpb.SpareRequest{
		Num: -1,
	})

	if err != nil{

		errStt,_ := status.FromError(err)

		log.Println("day la msg ",errStt.Message())

		if errStt.Code() == codes.InvalidArgument{
			log.Fatal("loi server")
		}
	}

	fmt.Println(resp.Result)

}
