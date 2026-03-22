package main

import (
	"fmt"
	// "math/rand"
)

// func isEven(n int) bool {
// 	if n%2 == 0 {
// 		return true
// 	} else {
// 		fmt.Println("false statement")
// 	}
// 	return false
// }

type rect struct {
	width int
	height int
}

func (r rect) area() int {
	return r.height * r.width;
}

func (r rect) perimeter() int {
	return 2 * (r.height + r.width);
}

// func area(r rect) int {
// 	return r.height * r.width
// }

// func perimeter(r rect) int {
// 	return 2 * (r.height + r.width);
// }

func main() {
	rectA := rect{width: 5, height: 10};

	fmt.Println(rectA.area());

	fmt.Println(rectA.perimeter());


	// n := rand.Int31n(12);

	// const a int = 1;
	// const b float32 = 134.123124;
	// fmt.Print("hello this is Dev: ", n, " ", a, " ", b);

	// const rangee = 12;
	// var ragee = 12;
	// for ; i<31; i++ {
	// 	fmt.Println(i);
	// }

	// i := 0;
	// for i := range rangee {
		// fmt.Println(i);
	// }

	// n := 51;
	// ans := isEven(n)
	// fmt.Println(ans);

	// a := [10]int{12, 5, 62, 2, 7, 23, 2, 356, 23, 5}
	// // var aPtr *[10]int = &a;
	// sum := 0

	// for i := range 10 {
	// 	// fmt.Println(aPtr[i])
	// 	sum += a[i]
	// }

	// fmt.Println(a, " sum: ", sum)


	// a := []int{123, 5,1,2,5};
	// b := make([]int, len(a));
	// copy(b, a);
	// b[0] = 123124;
	// fmt.Println(a);
	// fmt.Println(b);

	// temp := []int{}
	// fmt.Println(cap(temp))

	// a := make([]int, 3);
	// fmt.Println(a);
	// c := append(a, 1,2,3)
	// fmt.Println("a: ", a)
	// fmt.Println("c: ", c)

	// b := make([]int, 0, 6);
	// fmt.Println("b: ", b)
	// b = append(b, 1,24,1,2,51,12)
	// fmt.Println("b: ", b);


	// var mp map[string]string = map[string]string{
	// 	"email": "dev@gmail.com",
	// };
	// mp := make(map[string]string)

	// mp["name"] = "Dev"

	// fmt.Println(mp)

	// var value, exits = mp["email"] 
	// fmt.Println(value, " ", exits)


	// simple fun
	// fmt.Println(add(12, 5));

	//returns a pair

	// one, two := get12(a);

	// fmt.Println(one, " ", two)


	//anonymous func

	// get12 := func(arr []int) (a int, b int){
	// 	a = 0;
	// 	b = 0;
		
	// 	if len(arr) > 0 {	
	// 		a = arr[0];
	// 	}
	// 	if len(arr) > 1 {	
	// 		b = arr[1];
	// 	}

	// 	return;
	// }

	// one, two := call(get12);
	// fmt.Println(one, two)

}


//func getting a func as an argument
// func call(get12 func(a []int) (int,int)) (int,int) {
// 	arr := make([]int, 1)
// 	arr = []int{23};
// 	one, two := get12(arr);
// 	return one, two 
// }


// //return a pair
// func get12(arr []int) (a int, b int) {
// 	a = 0;
// 	b = 0;
// 	if len(arr) > 0 {
// 		a = arr[0];
// 	}	
// 	if len(arr) > 1 {
// 		b = arr[1];
// 	}	
// 	// a = arr[0];
// 	// b = arr[1];
// 	return;
// }

// //simple func
// func add(a int, b int) int {
// 	return a + b
// }
