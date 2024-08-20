package main

import (
	"errors"
	"fmt"
	"time"

	async "github.com/mgperkowski/goasyncawait/async"
)

func main() {

	promise := async.NewPromise(func(resolve func(interface{}), reject func(error)) {

		time.Sleep(2 * time.Second)
		if time.Now().Unix()%2 == 0 {
			resolve("Success: The promise was resolved")
		} else {
			reject(errors.New("Failure: The promise was rejected"))
		}
	})

	result, err := promise.Await()
	if err != nil {
		fmt.Println("Promise rejected with error:", err)
	} else {
		fmt.Println("Promise resolved with result:", result)
	}

	p1 := async.NewPromise(func(res func(interface{}), rej func(error)) {
		res("Result from p1")
	})

	p2 := async.NewPromise(func(res func(interface{}), rej func(error)) {
		res("Result from p2")
	})

	p3 := async.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve("Result from p3")
	})

	results, err := async.AwaitAll([]*async.Promise{p1, p2, p3})

	if err != nil {
		fmt.Println("Error in one of the promises:", err)
	} else {
		fmt.Println("All promises resolved with results:", results)
	}

	p4 := async.NewPromise(func(res func(interface{}), rej func(error)) {
		time.Sleep(2 * time.Second)
		res("Result from p4")
	})

	p5 := async.NewPromise(func(res func(interface{}), rej func(error)) {
		time.Sleep(3 * time.Second)
		res("Result from p5")
	})

	p6 := async.NewPromise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(1 * time.Second)
		resolve("Result from p6")
	})

	result, err = async.AwaitRace([]*async.Promise{p4, p5, p6})

	if err != nil {
		fmt.Println("Error in one of the promises:", err)
	} else {
		fmt.Println("The first promise to resolve was:", result)
	}

}
