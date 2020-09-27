package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")

	demo3()
	fmt.Println("process main")

}

func demo1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	panic("demo1 unknown err")

	fmt.Println("demo1 process")
}

func demo0() {
	panic("demo0 unknown err")

	fmt.Println("demo0 process")
}

func handpanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func demo2() {
	defer handpanic()

	panic("demo2 unknown err")

	fmt.Println("demo2 process")
}

func handle() {
	//直接调用,无法被recover
	handpanic()
	//在通过defer 调用，也无法recover
	// defer handpanic()
}

func demo3() {
	defer handle()

	panic("demo3 unknown err")

	fmt.Println("demo3 process")
}
