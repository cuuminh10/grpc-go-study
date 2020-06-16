package main

import (
	"fmt"
)

func main() {

	fmt.Println("chay")
	cWait := make(chan struct{})
	go func() {

		for  i := 0 ; i < 10 ; i++ {
			fmt.Println("chay ",i)


		}

	}()


	go func() {

		for  i := 0 ; i < 10 ; i++ {
			fmt.Println("chay2 ",i)
		}
		close(cWait)
	}()

	<-cWait


}