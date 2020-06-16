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
	"net"
	"time"
)

type server struct {}

func (s *server)  Sum(ct context.Context,req *calculatorpb.SumRequest) (* calculatorpb.SumResponse, error){
	log.Println("sum da call")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp,nil
}

func (s *server) Pbd(req *calculatorpb.PndRequest,stream calculatorpb.CalculatorService_PbdServer) error  {

	i := int32(2)
	sum := req.GetNum();

	for sum > 1{
		if (sum % i == 0){

			sum = sum / i
			fmt.Println("day la y ",i,sum)
			stream.Send(&calculatorpb.PndResponse{
				Result: i,
			})
			time.Sleep(5 * time.Second)
		}else {
			i++
		}
	}

	return nil;
}

func  (*server)	Avg(stream calculatorpb.CalculatorService_AvgServer) error{
	n := int32(1);
	for  {
		req,err := stream.Recv()
		if err == io.EOF{
			resp := &calculatorpb.AvgResponse{
				Result: int32(n),
			}
			log.Println("server tra ve ket qua ",resp)
			return stream.SendAndClose(resp)
		}
		n+=req.GetNum()
	}
}

func  (*server) Max(stream calculatorpb.CalculatorService_MaxServer) error {

	n:= int32(1)

	for{
		req,err :=stream.Recv()

		if err == io.EOF{

			log.Fatal("server het value tra ve")

		}

		if n < req.GetNum(){
			n = req.GetNum()
		}
		resp := &calculatorpb.MaxResponse{
			Result: n,
		}

		log.Println("server  tra ve ",n)

		errSend := stream.Send(resp)
		if errSend != nil{
			log.Println("server loi dung keyt noi ",errSend)
		}
	}



}

func  (*server)	Sqare(cont context.Context,req *calculatorpb.SpareRequest) (*calculatorpb.SpareResponse, error) {

	if req.GetNum() < 0{
		return nil,status.Errorf(codes.InvalidArgument,"so khong duoc la ",req.GetNum())
	}

	return  &calculatorpb.SpareResponse{
		Result: req.GetNum(),
	},nil
}



type student struct {
	name string
}

func main() {
	lis,err := net.Listen("tcp","0.0.0.0:50089")

	if (err != nil){
		log.Fatal("loi connect ",err)
	}

	e := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(e,&server{})
	fmt.Print("servee day running")
	err2 := e.Serve(lis)
	if (err2 != nil){
		log.Fatal("loi connect ",err2)
	}
}
