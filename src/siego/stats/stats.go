package stats

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"time"

	"../net"
	"strings"
)

// Stats - Siego stats structure
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

// NewStats - Creates statistics object
func NewStats() *Stats {
	s := Stats{
		codes: make(map[int]int),
		start: time.Now(),
	}

	return &s
}

// AddResponse - Track response
func (s *Stats) AddResponse(r *net.Response) {
	s.total = s.total + 1
	s.totalTime = s.totalTime + r.Duration
	s.times = append(s.times, r.Duration.Seconds())

	// track request
	if r.Error == nil && r.HttpResponse != nil && r.HttpResponse.StatusCode < 500 {
		s.codes[r.HttpResponse.StatusCode] = s.codes[r.HttpResponse.StatusCode] + 1
		s.success = s.success + 1
	}

	if r.Error != nil || (r.HttpResponse != nil && r.HttpResponse.StatusCode >= 500) {
		s.fail = s.fail + 1
	}

	// track duration
	if r.Duration > s.longest || s.longest == 0.0 {
		s.longest = r.Duration
	}

	if r.Duration < s.shortest || s.shortest == 0.0 {
		s.shortest = r.Duration
	}

	if r.HttpResponse != nil {
		if bytes, err := ioutil.ReadAll(r.HttpResponse.Body); err == nil {
			defer r.HttpResponse.Body.Close()

			s.bytes = s.bytes + len(bytes)
		}
	}
}

// String - Converts object to string
func (s *Stats) String() string {
	return s.getMainTable(false) + s.getResponseCodesTable(false) + s.getPercentilesTable(false)
}

// Xxml - Converts object to XML string
func (s *Stats) Xml() string {
	format := "<result>%s<response_codes>%s</response_codes><percentiles>%s</percentiles></result>"

	return fmt.Sprintf(format, s.getMainTable(true), s.getResponseCodesTable(true), s.getPercentilesTable(true))
}

// Get main data table
func (s *Stats) getMainTable(isXml bool) (result string) {
	elapsed := time.Since(s.start)

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
		"Transactions":            fmt.Sprintf("%0.0d", s.total),
		"Availability":            fmt.Sprintf("%0.2f%%", float64(s.total)/(float64(s.total)+float64(s.fail))*100),
		"Elapsed time":            fmt.Sprintf("%0.4fs", elapsed.Seconds()),
		"Data transferred":        fmt.Sprintf("%0.4fMb", float64(s.bytes)/(1024*1024)),
		"Response time":           fmt.Sprintf("%0.4fs", elapsed.Seconds()/float64(s.total)),
		"Transaction rate":        fmt.Sprintf("%0.4f/s", float64(s.total)/elapsed.Seconds()),
		"Throughput":              fmt.Sprintf("%0.4fMb/s", (float64(s.bytes)/elapsed.Seconds())/(1024.0*1024.0)),
		"Concurrency":             fmt.Sprintf("%0.4f", s.totalTime.Seconds()/elapsed.Seconds()),
		"Successful transactions": fmt.Sprintf("%0.0d", s.success),
		"Failed transactions":     fmt.Sprintf("%0.0d", s.fail),
		"Longest transaction":     fmt.Sprintf("%0.4fs", s.longest.Seconds()),
		"Shortest transaction":    fmt.Sprintf("%0.4fs\r\n", s.shortest.Seconds()),
	}

	for _, title := range sorts {
		if isXml {
			result = result + s.xmlRow(title, rows[title])
		} else {
			result = result + s.stringRow(title, rows[title])
		}
	}

	return result
}

// Get response codes table
func (s *Stats) getResponseCodesTable(isXml bool) (result string) {
	if !isXml {
		result = result + "\r\n" + s.stringRow("Response codes", "")
	}

	for key, value := range s.codes {
		if isXml {
			result = result + s.xmlRow(fmt.Sprintf("HTTP_%d", key), value)
		} else {
			result = result + s.stringRow(fmt.Sprintf("HTTP_%d", key), value)
		}
	}

	return result
}

// Get response time percentiles table
func (s *Stats) getPercentilesTable(isXml bool) (result string) {
	if !isXml {
		result = result + "\r\n" + s.stringRow("Response time percentiles", "")
	}

	sort.Float64s(s.times)
	for i := 10; i < 100; i += 10 {
		index := math.Floor(float64(len(s.times)) * (float64(i) / float64(100)))

		if len(s.times) > int(index) {
			value := fmt.Sprintf("%0.4fs", s.times[int(index)])
			title := fmt.Sprintf("p%d", i)

			if isXml {
				result = result + s.xmlRow(title, value)
			} else {
				result = result + s.stringRow(title, value)
			}
		}
	}

	return result
}

// Get pretty printed row
func (s *Stats) stringRow(title string, value interface{}) string {
	return fmt.Sprintf("%30s: %v\r\n", title, value)
}

// Get xml-formatted row
func (s *Stats) xmlRow(title string, value interface{}) string {
	title = strings.Replace(title, " ", "_", -1)
	title = strings.ToLower(title)

	return fmt.Sprintf("<%s>%v</%s>", title, value, title)
}
