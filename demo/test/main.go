package main

import "fmt"

func main() {
	for i := 0; i < 100; i++ {
		//go func() {
		//	for j := 0; j < 5; j++ {
		//		fmt.Printf("%d ", i)
		//	}
		//}()
		go printNum(i)
	}

}

func printNum(num int) {
	for j := 0; j < 5; j++ {
		fmt.Printf("%d ", num)
	}
}
