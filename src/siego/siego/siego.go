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

// Siego - Main Siego object
type Siego struct {
	url, file, log, time, userAgent, contentType string
	timeLimit                                    time.Duration
	get, post, internet, benchmark, xml          bool
	concurrent, delay, reps, timeout             int
	header                                       []string
	stats                                        *stats.Stats
}

// NewSiego - Constructor
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
		xml:         c.Bool("xml"),
		internet:    c.Bool("internet"),
		benchmark:   c.Bool("benchmark"),
		concurrent:  c.Int("concurrent"),
		delay:       c.Int("delay"),
		reps:        c.Int("reps"),
		timeout:     c.Int("timeout"),
		stats:       stats.NewStats(),
	}

	return &s
}

// Validate - Validates input parametes
func (s *Siego) Validate() (err error) {
	if s.file == "" && s.url == "" {
		return fmt.Errorf("You should specify 'url' or 'file' option.")
	}

	if s.file != "" {
		f, err := os.Open(s.file)

		if err != nil {
			return fmt.Errorf("Cannot open file %s: %s", s.file, err.Error())
		}

		defer f.Close()
	}

	if s.time != "" {
		if s.timeLimit, err = time.ParseDuration(s.time); err != nil {
			return fmt.Errorf("Cannot parse 'time' parameter: %s", err.Error())
		}

		s.reps = 0
	}

	return nil
}

// Run - Runs timed load test
func (s *Siego) Run() {
	if s.timeLimit.Seconds() > 0 {
		done := make(chan bool)

		go func() {
			s.doRun()
			done <- true
		}()

		for {
			select {
			case <-time.After(s.timeLimit):
				return
			case <-done:
				return
			}
		}
	} else {
		s.doRun()
	}
}

// Decides type of load test (file or url)
func (s *Siego) doRun() {
	if s.url != "" {
		s.runUrl()
	} else {
		s.runFile()
	}
}

// GetStats - Returns stats for test and optionally writes to log
func (s *Siego) GetStats() {
	data := ""

	if s.xml {
		data = s.stats.Xml()
	} else {
		data = fmt.Sprintf("%s\r\n", s.stats.String())
	}
	if s.log != "" {
		if f, e := os.OpenFile(s.log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); e == nil {
			defer f.Close()
			f.Write([]byte(data))
		}
	}

	fmt.Print(data)
}

// Executes load testing with url
func (s *Siego) runUrl() (err error) {
	var (
		req *net.Request
		ch  chan *net.Response = make(chan *net.Response)
	)

	method := "GET"
	if s.post {
		method = "POST"
	}

	req, err = net.NewRequest(method, s.url, "")
	if err != nil {
		return err
	}

	req.UserAgent(s.userAgent)
	req.ContentType(s.contentType)
	req.Headers(s.header)

	for i := 0; i < s.concurrent; i++ {
		go func(c chan *net.Response, r *net.Request, id int) {
			client := net.NewClient(s.timeout)

			if s.reps > 0 {
				for j := 0; j < s.reps; j++ {
					result := client.Do(r)
					s.doDelay()
					c <- result
				}
			} else {
				for {
					result := client.Do(r)
					s.doDelay()
					c <- result
				}
			}
		}(ch, req, i)
	}

	if s.reps > 0 {
		count := s.reps * s.concurrent
		for k := 0; k < count; k++ {
			s.stats.AddResponse(<-ch)
		}
	} else {
		for {
			s.stats.AddResponse(<-ch)
		}
	}

	return nil
}

// Executes load testing with urls file
func (s *Siego) runFile() error {
	var (
		ch chan *net.Response = make(chan *net.Response)
	)

	file, err := ioutil.ReadFile(s.file)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.Trim(string(file), "\r\n\t "), "\n")

	for i := 0; i < s.concurrent; i++ {
		go func(c chan *net.Response, lines []string) {
			client := net.NewClient(s.timeout)

			if s.reps > 0 {
				for j := 0; j < s.reps; j++ {
					if s.internet {
						lines = s.shuffle(lines)
					}

					for _, line := range lines {
						if r, e := s.requestFromLine(line); e == nil {
							result := client.Do(r)
							s.doDelay()
							c <- result
						}
					}
				}
			} else {
				for {
					if s.internet {
						lines = s.shuffle(lines)
					}

					for _, line := range lines {
						if r, e := s.requestFromLine(line); e == nil {
							result := client.Do(r)
							s.doDelay()
							c <- result
						}
					}
				}
			}
		}(ch, lines)
	}

	if s.reps > 0 {
		count := s.reps * s.concurrent * len(lines)
		for k := 0; k < count; k++ {
			s.stats.AddResponse(<-ch)
		}
	} else {
		for {
			s.stats.AddResponse(<-ch)
		}
	}

	return nil
}

// Do a delay between requests
func (s *Siego) doDelay() {
	if s.delay > 0 && !s.benchmark {
		delay := float64(rand.Intn(s.delay*1000)) / 1000

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// Creates request from file line
func (s *Siego) requestFromLine(line string) (req *net.Request, err error) {
	line = strings.Trim(line, "\t ")

	parts := strings.SplitAfterN(line, " ", 3)
	url := strings.Trim(parts[0], "\t ")

	method := "GET"
	if s.post {
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

	req.UserAgent(s.userAgent)
	req.ContentType(s.contentType)
	req.Headers(s.header)

	return req, err
}

// Shuffles slice
func (s *Siego) shuffle(slice []string) []string {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}
