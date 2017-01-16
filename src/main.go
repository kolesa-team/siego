package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./siego/siego"

	"github.com/codegangsta/cli"
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func main() {
	app := cli.NewApp()

	app.Name = "siego"
	app.Usage = "Regression test and benchmark utility"
	app.Version = "0.0.3"
	app.Author = "Igor Borodikhin"
	app.Email = "iborodikhin@gmail.com"
	app.Action = actionRun
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "concurrent, c",
			Value: 10,
			Usage: "CONCURRENT users.",
		},
		cli.IntFlag{
			Name:  "delay, d",
			Value: 1,
			Usage: "Time DELAY, random delay before each requst between 1 and NUM. (NOT COUNTED IN STATS)",
		},
		cli.IntFlag{
			Name:  "reps, r",
			Usage: "REPS, number of times to run the test.",
		},
		cli.StringFlag{
			Name:  "url, u",
			Usage: "URL to test.",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "FILE, select a specific URLS FILE.",
		},
		cli.StringFlag{
			Name:  "log, l",
			Value: "/var/siege.log",
			Usage: "LOG to FILE.",
		},
		cli.StringFlag{
			Name:  "time, t",
			Usage: "TIMED testing where \"m\" is modifier s, m, or h. Ex: --time=1h, one hour test.",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "Add a header to request (can be many)",
		},
		cli.StringFlag{
			Name:  "user-agent, A",
			Usage: "Sets User-Agent in request",
		},
		cli.StringFlag{
			Name:  "content-type, T",
			Usage: "Sets Content-Type in request",
		},
		cli.BoolFlag{
			Name:  "get, g",
			Usage: "Use GET method.",
		},
		cli.BoolFlag{
			Name:  "post, p",
			Usage: "Use POST method.",
		},
		cli.BoolFlag{
			Name:  "internet, i",
			Usage: "INTERNET user simulation, hits URLs randomly.",
		},
		cli.BoolFlag{
			Name:  "benchmark, b",
			Usage: "BENCHMARK: no delays between requests.",
		},
		cli.BoolFlag{
			Name:  "xml, x",
			Usage: "Use XML output.",
		},
		cli.IntFlag{
			Name:  "timeout",
			Value: 1,
			Usage: "Request timeout in seconds.",
		},
	}

	app.Run(os.Args)
}

func actionRun(c *cli.Context) error {
	s := siego.NewSiego(c)

	err := s.Validate()
	if err != nil {
		return err
	}

	fmt.Printf("Server now under siego...\r\n")

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		s.GetStats()
		os.Exit(1)
	}()

	s.Run()
	s.GetStats()

	return nil
}
