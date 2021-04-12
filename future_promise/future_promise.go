package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// 原文： http://labs.strava.com/blog/futures-in-golang/
// promise/future是一个非常重要的异步编程模型，它可以让我们摆脱传统的回调陷阱，从而使用更加优雅、清晰的方式进行异步编程

func RequestFuture(url string) <-chan []byte {
	c := make(chan []byte, 1)

	go func() {
		var body []byte
		defer func() {
			c <- body
		}()
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	}()

	return c
}

func RequestFutureFunction(url string) func() ([]byte, error) {
	var body []byte
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	}()
	return func() ([]byte, error) {
		<-c
		return body, err
	}
}

func Future(f func() (interface{}, error)) func() (interface{}, error) {
	var result interface{}
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()
	return func() (interface{}, error) {
		<-c
		return result, err
	}
}

func future_demo() {
	url := "http://labs.strava.com"
	future := Future(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
	// do many other things
	b, err := future()
	body, _ := b.([]byte)
	log.Printf("response length: %d", len(body))
	log.Printf("request error: %v", err)
}

func main() {
	future := RequestFuture("http://labs.strava.com")
	// do many other things, maybe create other futures
	body := <-future
	log.Printf("response length: %d", len(body))
}
