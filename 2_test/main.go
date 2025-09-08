package main

import (
	"fmt"
	"sync"
	// "sync/atomic"
	// "time"
	// "time"
)

// func fn(str string) {
// 	for i := range 10 {
// 		fmt.Println(str, " - ", i);
// 	}
// }

var counter = 0;

func increment (ws *sync.WaitGroup){
	defer ws.Done();
	for range 10{
		counter += 1;
	}
}

func main() {

	// c := make(chan int);

	// go func () {
	// 	temp := <-c;
	// 	// for i := range 9 {
	// 	// 	fmt.Println(i);
	// 	// }
	// 	fmt.Println(temp);
	// 	c <- 1;
	// }()

	// c <- 10;
	// <-c

	// time.Sleep(time.Second * 2); 

	// fmt.Println("process complete");


	//buffered channel of size n;
	// n := 40;
	// c := make(chan int, n);


	// go func() {
	// 	for i := range 50 {
	// 		fmt.Println(i, " sent !");
	// 		c <- i
	// 	}
	// }()

	// // time.Sleep(time.Second * 1);
	// for i := range n {
	// 	fmt.Println(i, ": ", <-c);
	// }


	//waitgrps
	// var wg sync.WaitGroup;
	// n := 10;
	// var factorial atomic.Uint64;
	// factorial.Store(1);

	// // adding the sum from 0 to 1000000

	// for i := range 5 {
	// 	smallUnit := n/5;
	// 	wg.Add(1);
	// 	go func () {
	// 		defer wg.Done();
	// 		for j := range smallUnit {
	// 			factorial *= j + 1 + (smallUnit * i);
	// 		}
	// 	}()
	// }

	// wg.Wait();
	// fmt.Println(factorial);


	expectedCounter := 100;

	var ws sync.WaitGroup;

	for range 10 {
		ws.Add(1);
		go increment(&ws);
	}

	ws.Wait()

	if(expectedCounter == counter){
		fmt.Println("race condition not detected", expectedCounter, " ", counter);
	}else{
		fmt.Println("race condition detected", expectedCounter, " ", counter);
	}
}