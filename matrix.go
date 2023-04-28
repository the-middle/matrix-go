package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	gc "github.com/rthornton128/goncurses"
)

var DEBUG int = 0

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom() int {
	rnd := rand.Intn(ir.max-ir.min+1) + ir.min
	// log.Println("RND: ", rnd)
	return rnd
}

func main() {
	os.Rename("logs.log", "logs.old.log")
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	var wg sync.WaitGroup

	log.SetOutput(f)

	s, err := gc.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer gc.End()

	gc.Cursor(0)
	gc.Raw(false)
	gc.Echo(false)
	// s.Standend()
	// s.Print("\033[?9l")
	

	y, x := s.MaxYX()
	tx := x / 5
	if DEBUG == 1 {
		s.Move(5, 2)
		s.Println("height: ", y, " width: ", x, "tx: ", tx)
		printAllChars(y, x, s)
		// for i := 1; i <= 5; i++ {
		// 	if i == 1 {
		// 		s.Println("makeRange 1: ", makeRange(0, tx))
		// 		time.Sleep(time.Millisecond * 15)
		// 	} else {
		// 		s.Println("makeRange 1+: ", makeRange(tx*(i-1), tx*i))
		// 		time.Sleep(time.Millisecond * 15)
		// 	}
		// }
	} else {
		for {
			for i := 1; i <= 5; i++ {
				wg.Add(1)
				i := i
				if i == 1 {
					go func() {
						defer wg.Done()
						printRandom(y, makeRange(0, tx), s, i)
						time.Sleep(time.Millisecond * 15)
					}()
				} else {
					go func() {
						defer wg.Done()
						printRandom(y, makeRange(tx*(i-1), tx*i), s, i)
						time.Sleep(time.Millisecond * 15)
					}()
				}

			}
			wg.Wait()
			input := s.GetChar()
			if input == gc.KEY_ESC {
				break
			} else {
				s.Erase()
				continue
			}
		}
	}

	// s.Move(5, 2)
	// s.Println("Hello, height: ", y, " width: ", x)

}

func printAllChars(y int, x int, s *gc.Window) {
	chr_range := makeRange(32, 126)
	for _, chr := range chr_range {
		s.Move(1, 15)
		s.Refresh()
		s.Printf("Start: '%c'", chr)
	}
}

func printRandom(y int, x []int, s *gc.Window, n int) {
	var nrch int
	rand.Seed(time.Now().UnixNano())

	ir := IntRange{x[1], x[len(x)-1]}
	// rch := IntRange{32, 126}
	rch := IntRange{65, 122} //IntRange{65, 122}

	for i := x[0]; i <= x[len(x)-1]; i++ {
		rx := ir.NextRandom()
		for j := 0; j < y; j++ {
			nrch = rch.NextRandom()
			// if nrch == 58 || nrch == 59 || nrch == 60 || nrch == 62 || nrch == 63 || nrch == 64 || nrch == 91 || nrch == 92 || nrch == 93 || nrch == 94 || nrch == 107 || nrch == 108 || nrch == 109 || nrch == 110 {
			// 	continue
			// } else {
			gc_nrch := gc.Char(nrch)
			s.Standend()

			s.Move(j, rx)
			// s.AddChar(gc.Char(nrch) | gc.A_NORMAL)
			// s.Printf("%c", nrch)
			// s.AttrSet(gc.A_NORMAL)
			s.AttrOff(gc_nrch)
			s.MoveAddChar(j, rx, gc_nrch)
			log.Printf("[%v]CHAR: '%s'", n, fmt.Sprintf("%c", nrch))
			s.Refresh()
			// gc.FlushInput()
			time.Sleep(time.Millisecond * 10)
			// }
		}
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
