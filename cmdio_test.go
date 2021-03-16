package cmdio

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestProgressBar(t *testing.T) {
	var wg sync.WaitGroup

	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)

	go func() { c1 <- 5 }()
	go func() { c2 <- 10 }()
	go func() { c3 <- 15 }()
	NewProgressBar(ProgressOptions{Counters:true,Percentage:true}, c1, c2, c3)

	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			c1 <- 1
			time.Sleep(1 * time.Second)
		}
		log.Println("c1 done")
		close(c1)
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			c2 <- 1
			time.Sleep(2 * time.Second)
		}
		log.Println("c2 done")
		close(c2)
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 15; i++ {
			c3 <- 1
			time.Sleep(1 * time.Second)
		}
		log.Println("c3 done")
		close(c3)
	}()

	wg.Wait()

	time.Sleep(1 * time.Second)
	ClearLine("done")

}
