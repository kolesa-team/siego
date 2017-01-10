package siego

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"../net"
	"../stats"

	"github.com/codegangsta/cli"
)

type Siego struct {
	url, file, log, time, userAgent, contentType string
	timeLimit                                    time.Duration
	get, post, internet, benchmark               bool
	concurrent, delay, reps                      int
	header                                       []string
	client                                       *net.Client
	stats                                        *stats.Stats
}

// Constructor
func NewSiego(c *cli.Context) *Siego {
	s := Siego{
		file:        c.String("file"),
		url:         c.String("url"),
		log:         c.String("log"),
		time:        c.String("time"),
		userAgent:   c.String("user-agent"),
		contentType: c.String("content-type"),
		header:      c.StringSlice("header"),
		get:         c.Bool("get"),
		post:        c.Bool("post"),
		internet:    c.Bool("internet"),
		benchmark:   c.Bool("benchmark"),
		concurrent:  c.Int("concurrent"),
		delay:       c.Int("delay"),
		reps:        c.Int("reps"),
		client:      net.NewClient(),
		stats:       stats.NewStats(),
	}

	return &s
}

// Validates input parametes
func (this *Siego) Validate() (err error) {
	if this.file == "" && this.url == "" {
		return fmt.Errorf("You should specify 'url' or 'file' option.")
	}

	if this.file != "" {
		f, err := os.Open(this.file)

		if err != nil {
			return fmt.Errorf("Cannot open file %s: %s", this.file, err.Error())
		}

		defer f.Close()
	}

	if this.time != "" {
		if this.timeLimit, err = time.ParseDuration(this.time); err != nil {
			return fmt.Errorf("Cannot parse 'time' parameter: %s", err.Error())
		}

		this.reps = 0
	}

	return nil
}

// Runs timed load test
func (this *Siego) Run() {
	if this.timeLimit.Seconds() > 0 {
		done := make(chan bool)

		go func() {
			this.doRun()
			done <- true
		}()

		for {
			select {
			case <-time.After(this.timeLimit):
				return
			case <-done:
				return
			}
		}
	} else {
		this.doRun()
	}
}

// Decides type of load test (file or url)
func (this *Siego) doRun() {
	if this.url != "" {
		this.runUrl()
	} else {
		this.runFile()
	}
}

// Returns stats for test and optionally writes to log
func (this *Siego) GetStats() {
	data := fmt.Sprintf("%s\r\n", this.stats)
	if this.log != "" {
		if f, e := os.OpenFile(this.log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); e == nil {
			defer f.Close()
			f.Write([]byte(data))
		}
	}

	fmt.Print(data)
}

// Executes load testing with url
func (this *Siego) runUrl() (err error) {
	var (
		req *net.Request
		ch  chan *net.Response = make(chan *net.Response)
	)

	method := "GET"
	if this.post {
		method = "POST"
	}

	req, err = net.NewRequest(method, this.url, "")
	if err != nil {
		return err
	}

	req.UserAgent(this.userAgent)
	req.ContentType(this.contentType)
	req.Headers(this.header)

	for i := 0; i < this.concurrent; i++ {
		go func(c chan *net.Response, r *net.Request) {
			if this.reps > 0 {
				for j := 0; j < this.reps; j++ {
					c <- net.NewClient().Do(r)
					this.doDelay()
				}
			} else {
				for {
					c <- net.NewClient().Do(r)
					this.doDelay()
				}
			}
		}(ch, req)
	}

	if this.reps > 0 {
		count := this.reps * this.concurrent
		for k := 0; k < count; k++ {
			this.stats.AddResponse(<-ch)
		}
	} else {
		for {
			this.stats.AddResponse(<-ch)
		}
	}

	return nil
}

// Executes load testing with urls file
func (this *Siego) runFile() error {
	var (
		ch chan *net.Response = make(chan *net.Response)
	)

	file, err := ioutil.ReadFile(this.file)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.Trim(string(file), "\r\n\t "), "\n")

	for i := 0; i < this.concurrent; i++ {
		go func(c chan *net.Response, lines []string) {
			if this.reps > 0 {
				for j := 0; j < this.reps; j++ {
					if this.internet {
						lines = this.shuffle(lines)
					}

					for _, line := range lines {
						if r, e := this.requestFromLine(line); e == nil {
							c <- net.NewClient().Do(r)
							this.doDelay()
						}
					}
				}
			} else {
				for {
					if this.internet {
						lines = this.shuffle(lines)
					}

					for _, line := range lines {
						if r, e := this.requestFromLine(line); e == nil {
							c <- net.NewClient().Do(r)
							this.doDelay()
						}
					}
				}
			}
		}(ch, lines)
	}

	if this.reps > 0 {
		count := this.reps * this.concurrent * len(lines)
		for k := 0; k < count; k++ {
			this.stats.AddResponse(<-ch)
		}
	} else {
		for {
			this.stats.AddResponse(<-ch)
		}
	}

	return nil
}

// Do a delay between requests
func (this *Siego) doDelay() {
	if this.delay > 0 && !this.benchmark {
		time.Sleep(time.Duration(this.delay) * time.Second)
	}
}

// Creates request from file line
func (this *Siego) requestFromLine(line string) (req *net.Request, err error) {
	line = strings.Trim(line, "\t ")

	parts := strings.SplitAfterN(line, " ", 3)
	url := strings.Trim(parts[0], "\t ")

	method := "GET"
	if this.post {
		method = "POST"
	}

	// Get request method from line
	if len(parts) > 1 {
		parts[1] = strings.Trim(parts[1], "\t ")

		if strings.ToUpper(parts[1]) == "GET" || strings.ToUpper(parts[1]) == "POST" {
			method = strings.ToUpper(parts[1])
		}
	}

	args := ""
	if len(parts) > 2 {
		args = parts[2]
	}

	req, err = net.NewRequest(method, url, args)
	if err != nil {
		return req, err
	}

	req.UserAgent(this.userAgent)
	req.ContentType(this.contentType)
	req.Headers(this.header)

	return req, err
}

// Shuffles slice
func (this *Siego) shuffle(slice []string) []string {
	for i, _ := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}
