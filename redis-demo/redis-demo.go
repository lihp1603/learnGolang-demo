package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var (
	server   string = "47.106.86.35:6400"
	password string = "aliyun-redis"
)

var pool *redis.Pool

func test(i int) {

	c := pool.Get()
	defer c.Close()

	t := strconv.Itoa(i)
	// 设置key对应字符串value，并且设置key在给定的seconds时间之后超时过期
	// reply, err := c.Do("SETEX", "foo"+t, 20, i)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(reply)
	// }

	reply, err := redis.Int(c.Do("GET", "foo"+t))
	if err == nil {
		fmt.Println(reply)
	} else {
		fmt.Println(err)
		fmt.Println(reply)
	}
	time.Sleep(1 * time.Second)
}

func poolInit() *redis.Pool {
	//redis pool
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}

			//if _, err := c.Do("SELECT",1); err != nil {
			// c.Close()
			// return nil, err
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main() {

	pool = poolInit()

	for i := 0; i < 1000; i++ {
		test(i)
	}
}
