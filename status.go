package cmdio

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

// A progress bar with multiple counters should look like this:
//
// [P%][##########..............][C1][C2]...[Cn]
//
// where the current progress is based on:
// completed/total
//
// The display of the bar should be calculated by:
// border + leftover space
//
// The percentage, error count, and counter display can all be
// made optional via a ProgressOptions int (which is just a bitfield)
//
// The progress bar should self terminate once all
// counters reach their total (in essence, total
// complete == C1 + C2 + ... Cn total)
//
// The progress bar function should take an n number of
// channels that indicate the current progress of the n functions
// in int form. Examples:
//   - If something is at stage 2/5, the progress bar should recieve 2
//   - If something is at the 250th n of n, the progress bar should recieve 250
//
// The progress bar function will continuously range over every single channel,
// until the channel is closed.
//
// When a progress bar is created, a progressCounter is created for every channel
// given. Initially, the channels should send the amount of work expected to be done,
// so that the progress bar can calculate its width per update.
//
// Afterwards, the progress bar should accept all updates once it is initialized.
// For every update, one is added to complete. Once complete == total, of course,
// the progress bar should terminate.
//
// Done: a lot of this is implicitly done when channels are closed.
//
// TODO: Error display

type progressBar struct {
	bar struct {
		m *sync.Mutex
		s string
	}
	counters        []*progressCounter
	messages        chan string
	complete, total float64
	options         ProgressOptions
	oldLogger       io.Writer
}

type progressCounter struct {
	c chan int
	l int
	t int
}

// ProgressOptions represents the options you can pass
// to NewProgressBar for what to display while it runs.
//
// TODO: Errors
type ProgressOptions struct {
	Percentage bool
	Counters   bool

}

// NewProgressBar creates a progress bar and starts its
// main loop in the current terminal that the program
// that uses this calls it in.
//
// Takes a ProgressOptions struct
// and an n amount of integer channels for
// progress bar counters.
//
// Any and all log messages should be output via
// the "log" package to the default logger. Otherwise,
// your messages will end up interleaving with the
// progress bar - or wosre, not appearing at all.
//
// Your terminal/command line should support at least
// either tput (*nix, Linux), or should be able to
// access PowerShell (Windows)
func NewProgressBar(p ProgressOptions, c ...chan int) {
	b := new(progressBar)
	b.bar.m = new(sync.Mutex)
	b.messages = make(chan string, 5)
	b.options = p

	for _, v := range c {
		a := new(progressCounter)
		a.c = v                            // set the counter's channel
		a.t = <-a.c                        // get the total
		b.total = b.total + float64(a.t)   // add it to the total amount of work
		b.counters = append(b.counters, a) // add to counters
	}

	b.oldLogger = log.Writer()
	log.SetOutput(b)

	go b.loop()
}

// TODO: Counter and bar display options (e.g., make the bar
// display look different)
func (p *progressBar) makeBar() error {
	c, err := GetCols()
	if err != nil {
		return err // higher up, the bar should either ignore or make a fuss about this
	}

	var t string
	if p.options.Percentage {
		p := math.Floor((p.complete / p.total) * 100)
		t = "[" + strconv.FormatFloat(p, 'f', -1, 32) + "]"
	}

	var o string
	if p.options.Counters {
		for _, c := range p.counters {
			o = o + "[" + strconv.Itoa(c.l) + "/" + strconv.Itoa(c.t) + "]"
		}
	}

	// calculate the bar length based on the column length minus whatever options currently exist
	// and also the border as well
	l := c - len(t) - len(o) - 2

	r := "["
	for i := 0.0; i < float64(l); i++ {
		if p.complete / p.total > i / float64(l) {
			r = r + "#"
		} else {
			r = r + "."
		}
	}
	r = r + "]"

	p.bar.s = t + r + o + fmt.Sprintf("\033[%dD", c)
	return nil
}

func (p *progressBar) loop() {
	var wg sync.WaitGroup
	for _, c := range p.counters {
		o := c
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for i := range o.c {
				o.l = o.l + i
				p.complete = p.complete + float64(i)
				p.update()
			}
		}(&wg)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for m := range p.messages {
			p.bar.m.Lock()
			ClearLine(m)
			fmt.Printf(p.bar.s)
			p.bar.m.Unlock()
		}
	}(&wg)

	wg.Wait()

	log.SetOutput(os.Stdout)
	close(p.messages)
}

func (p *progressBar) update() {
	p.bar.m.Lock()
	p.makeBar()
	fmt.Printf(p.bar.s)
	p.bar.m.Unlock()
}

// Write provides io.Writer compability, so that
// a message can be written to while the progress bar
// works.
func (b *progressBar) Write(p []byte) (n int, err error) {
	b.messages <- string(p)

	return len(p), err
}
