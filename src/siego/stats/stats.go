package stats

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"time"

	"../net"
)

type Stats struct {
	start     time.Time
	longest   time.Duration
	shortest  time.Duration
	totalTime time.Duration
	total     int
	success   int
	fail      int
	bytes     int
	codes     map[int]int
	times     []float64
}

// Constructor
func NewStats() *Stats {
	s := Stats{
		codes: make(map[int]int),
		start: time.Now(),
	}

	return &s
}

// Track response
func (this *Stats) AddResponse(r *net.Response) {
	this.total = this.total + 1
	this.totalTime = this.totalTime + r.Duration
	this.times = append(this.times, r.Duration.Seconds())

	// track request
	if r.Error == nil && r.HttpResponse != nil && r.HttpResponse.StatusCode < 500 {
		this.codes[r.HttpResponse.StatusCode] = this.codes[r.HttpResponse.StatusCode] + 1
		this.success = this.success + 1
	}

	if r.Error != nil || (r.HttpResponse != nil && r.HttpResponse.StatusCode >= 500) {
		this.fail = this.fail + 1
	}

	// track duration
	if r.Duration > this.longest || this.longest == 0.0 {
		this.longest = r.Duration
	}

	if r.Duration < this.shortest || this.shortest == 0.0 {
		this.shortest = r.Duration
	}

	if r.HttpResponse != nil {
		if bytes, err := ioutil.ReadAll(r.HttpResponse.Body); err == nil {
			defer r.HttpResponse.Body.Close()

			this.bytes = this.bytes + len(bytes)
		}
	}
}

// Converts object to string
func (this *Stats) String() string {
	return this.getMainTable() + this.getResponseCodesTable() + this.getPercentilesTable()
}

// Get main data table
func (this *Stats) getMainTable() (result string) {
	elapsed := time.Since(this.start)

	sorts := []string{
		"Transactions",
		"Availability",
		"Elapsed time",
		"Data transferred",
		"Response time",
		"Transaction rate",
		"Throughput",
		"Concurrency",
		"Successful transactions",
		"Failed transactions",
		"Longest transaction",
		"Shortest transaction",
	}
	rows := map[string]interface{}{
		"Transactions":            fmt.Sprintf("%0.0d", this.total),
		"Availability":            fmt.Sprintf("%0.4f", float64(this.total)/(float64(this.total)+float64(this.fail))),
		"Elapsed time":            fmt.Sprintf("%0.4fs", elapsed.Seconds()),
		"Data transferred":        fmt.Sprintf("%0.4fMb", float64(this.bytes)/(1024*1024)),
		"Response time":           fmt.Sprintf("%0.4fs", elapsed.Seconds()/float64(this.total)),
		"Transaction rate":        fmt.Sprintf("%0.4f/s", float64(this.total)/elapsed.Seconds()),
		"Throughput":              fmt.Sprintf("%0.4fMb/s", float64(this.bytes)/elapsed.Seconds()),
		"Concurrency":             fmt.Sprintf("%0.4f", this.totalTime.Seconds()/elapsed.Seconds()),
		"Successful transactions": fmt.Sprintf("%0.0d", this.success),
		"Failed transactions":     fmt.Sprintf("%0.0d", this.fail),
		"Longest transaction":     fmt.Sprintf("%0.4fs", this.longest.Seconds()),
		"Shortest transaction":    fmt.Sprintf("%0.4fs\r\n", this.shortest.Seconds()),
	}

	for _, title := range sorts {
		result = result + this.stringRow(title, rows[title])
	}

	return result
}

// Get response codes table
func (this *Stats) getResponseCodesTable() (result string) {
	result = result + "\r\n" + this.stringRow("Response codes", "")
	for key, value := range this.codes {
		result = result + this.stringRow(fmt.Sprintf("HTTP_%d", key), value)
	}

	return result
}

// Get response time percentiles table
func (this *Stats) getPercentilesTable() (result string) {
	result = result + "\r\n" + this.stringRow("Response time percentiles", "")
	sort.Float64s(this.times)
	for i := 10; i < 100; i += 10 {
		index := math.Floor(float64(len(this.times)) * (float64(i) / float64(100)))

		if len(this.times) > int(index) {
			value := fmt.Sprintf("%0.4fs", this.times[int(index)])
			title := fmt.Sprintf("%d%%", i)

			result = result + this.stringRow(title, value)
		}
	}

	return result
}

// Get pretty printed row
func (this *Stats) stringRow(title string, value interface{}) string {
	return fmt.Sprintf("%30s: %v\r\n", title, value)
}
