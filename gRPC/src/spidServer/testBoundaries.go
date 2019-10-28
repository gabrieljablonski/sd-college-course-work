package main

import (
	"fmt"
	"math/rand"
	"spidServer/gps"
	"spidServer/grpcServer"
	"strconv"
)

func main() {
	s := grpcServer.NewServer("1111")
	for _, b := range s.Boundaries {
		fmt.Printf("%s\n", b.ToString())
	}
	for i := 0; i < 10; i++ {
		p := gps.GlobalPosition{
			Latitude:  float64(rand.Intn(181) - 90),
			Longitude: float64(rand.Intn(361) - 180),
		}
		fmt.Printf("%s\n%s\n", p.ToString(), strconv.Itoa(s.WhereIsPosition(p)))
	}
}
