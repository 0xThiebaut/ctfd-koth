package monitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0xThiebaut/ctfd-koth/conf"
	"github.com/0xThiebaut/ctfd-koth/logger"
	"golang.org/x/text/encoding/unicode"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Monitor is in charge of monitoring the flags and submitting points to CTFd.
type Monitor struct {
	ticker *time.Ticker
	closer chan interface{}
	Flag   *conf.Flag
	api    *conf.API
}

// New creates a monitoring object given a flag and API configuration.
func New(flag *conf.Flag, api *conf.API) *Monitor {
	return &Monitor{Flag: flag, api: api}
}

// Start launches the monitoring agent.
// After some basic checks, which can error, the agent is launched asynchronously.
func (m *Monitor) Start() error {
	// Fail if already started
	if m.closer != nil || m.ticker != nil {
		return errors.New("monitoring agent already started")
	}

	// Parse the API URL
	api, err := url.Parse(m.api.URL)
	if err != nil {
		return err
	}
	// Deduce the award's API endpoint
	endpoint, err := url.Parse("awards")
	if err != nil {
		return err
	}
	u := api.ResolveReference(endpoint)

	client := http.Client{}
	// Initiate the state
	m.closer = make(chan interface{})
	m.ticker = time.NewTicker(m.Flag.Interval)

	// Run in new routine
	go func() {
		for {
			select {
			case <-m.ticker.C:
				// Open the flag
				file, err := os.Open(m.Flag.Path)
				if err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Get the contents
				var buf bytes.Buffer
				if _, err := buf.ReadFrom(file); err != nil {
					logger.Warn.Println(err)
					continue
				}
				c := buf.Bytes()
				// Parse the user identifier
				user, err := parse(c)
				if err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Build the award's JSON
				a := m.Flag.Award
				a.User = user
				j, err := json.Marshal(a)
				if err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Create a new POST request
				r, err := http.NewRequest("POST", u.String(), bytes.NewReader(j))
				if err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Define some required headers
				r.Header.Add("Content-Type", "application/json")
				r.Header.Add("Accept", "application/json")
				r.Header.Add("Authorization", "Bearer " + m.api.Token)
				// Execute the request
				resp, err := client.Do(r)
				if err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Get the response
				b := new(bytes.Buffer)
				if _, err := b.ReadFrom(resp.Body); err != nil {
					logger.Warn.Println(err)
					continue
				}
				if err := resp.Body.Close(); err != nil {
					logger.Warn.Println(err)
					continue
				}
				// Log the response
				logger.Info.Println(b.String())
			case <-m.closer:
				// Abort the loop
				return
			}
		}
	}()
	return nil
}

func (m *Monitor) Close() error {
	// Fail if not yet started
	if m.closer == nil || m.ticker == nil {
		return errors.New("monitoring agent needs to be started before a shut-down can be requested")
	}
	// Terminate the ticker
	m.ticker.Stop()
	// Send a closure signal
	close(m.closer)
	return nil
}

func parse(b []byte) (int, error) {
	// Try the default
	s := string(b)
	s = strings.TrimSpace(s)
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}
	// Try other encodings
	for _, e := range unicode.All {
		d := e.NewDecoder()
		o, err := d.Bytes(b)
		if err != nil {
			continue
		}
		s = string(o)
		s = strings.TrimSpace(s)
		if i, err := strconv.Atoi(s); err == nil {
			return i, nil
		}
	}
	return -1, fmt.Errorf("unable to parse %#v; either the file doesn't contain a number or the encoding is not any of ASCII, UTF8, UTF16BE, UTF16BE with BOM, UTF16LE", b)
}