package httpsum

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	defaultParallel = 10
	defaultTimeout  = 3
)

type Config struct {
	Client   HttpClient
	Parallel int
	Timeout  int
}

type HttpSum struct {
	client   HttpClient
	parallel int
	timeout  int
}

func New(c Config) (HttpSum, error) {
	if c.Client == nil {
		return HttpSum{}, errors.New("Config.Client should not be nil")
	}

	if c.Parallel > 10 {
		log.Printf("maximum number of parallel goroutine is 10; set to default 10 from %d", c.Parallel)
		c.Parallel = 10
	}
	if c.Parallel == 0 {
		c.Parallel = defaultParallel
	}

	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}

	h := HttpSum{
		client:   c.Client,
		parallel: c.Parallel,
		timeout:  c.Timeout,
	}

	return h, nil
}

func (h *HttpSum) Ping(sites []string) error {
	var wg sync.WaitGroup

	jobs := make(chan string, h.parallel)
	results := make(chan siteResponse, h.parallel)

	log.Printf("creating %d goroutine", h.parallel)
	for i := 0; i < h.parallel; i++ {
		wg.Add(1)
		go func(jobs <-chan string, resp chan<- siteResponse) {
			defer wg.Done()
			for site := range jobs {
				r := siteResponse{
					site: site,
				}

				md5, err := h.get(site)
				if err != nil {
					r.err = err.Error()
					resp <- r
				} else {
					r.success = true
					r.md5 = md5
					resp <- r
				}
			}
		}(jobs, results)
	}

	for _, site := range sites {
		jobs <- site
	}

	for i := 0; i < len(sites); i++ {
		r := <-results
		if r.success {
			fmt.Printf("%s \t\t %x\n", r.site, r.md5)
		} else {
			fmt.Printf("%s \t\t %s\n", r.site, r.err)
		}
	}
	close(jobs)
	wg.Wait()
	return nil
}

func (h *HttpSum) get(site string) ([md5.Size]byte, error) {
	var result [md5.Size]byte

	u, err := url.Parse(site)
	if err != nil {
		return result, fmt.Errorf("url parsing error %v: %w", site, err)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return result, fmt.Errorf("request error %v: %w", site, err)
	}

	var resp *http.Response

	ch := make(chan struct{})
	go func() {
		resp, err = h.client.Do(request)
		close(ch)
	}()

	select {
	case <-ch:
	// pass on
	case <-time.After(time.Duration(h.timeout) * time.Second):
		return result, fmt.Errorf("timeout on %#q, %w", site, timeoutError)
	}

	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to connect %#q, status code %d, %w", site, resp.StatusCode, httpStatusError)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	result = md5.Sum(body)
	return result, nil
}
