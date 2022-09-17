package main

import (
	"fmt"
	"time"

	"github.com/qiwik/golru"
)

func main() {
	c, _ := golru.NewCache(4, golru.WithTTL(0.8))
	c.Expire()

	c.Add("first", 1)
	c.Add("second", 2)
	fmt.Println(c.Len(), "first")

	time.Sleep(1 * time.Second)

	fmt.Println(c.Len(), "second")

	c.Add("third", 3)
	fmt.Println(c.Len(), "third")

	time.Sleep(1 * time.Second)
	fmt.Println(c.Len(), "fourth")
	fmt.Println(c.Keys(), c.Values())
}
