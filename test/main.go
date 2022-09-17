package main

import (
	"fmt"
	"github.com/qiwik/golru"
	"time"
)

func main() {
	fmt.Println(0.5 * 10)
	c, _ := golru.NewCache(4, 0.5)
	c.Watch()

	c.Add("first", 1)
	c.Add("second", 2)
	fmt.Println(c.Len(), "first")

	time.Sleep(2 * time.Second)

	fmt.Println(c.Len(), "second")

	c.Add("third", 3)
	fmt.Println(c.Len(), "third")

	time.Sleep(3 * time.Second)
	fmt.Println(c.Len(), "fourth")
}
