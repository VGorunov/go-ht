package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	url        string
	countOfReq int
	timeOut    int
)

const (
	urlDefault     = "https://www.google.com/"
	countDefault   = 1
	timeOutDefault = 500000000
)

type OutputResult struct {
	totalTimeReq    int
	averageTimeReq  int
	maxResponseTime int
	minResponseTime int
	countMissResp   int
}

func sendReq() (out OutputResult, err error) {
	var mutex sync.Mutex
	var myWg sync.WaitGroup
	var allTimeReq []int
	var countMissResp int

	for i := 0; i < countOfReq; i++ {
		myWg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			start := time.Now()
			client := http.Client{
				Timeout: time.Duration(timeOut),
			}
			_, err := client.Get(url)

			timeReq := int(time.Since(start).Nanoseconds())

			if err, ok := err.(net.Error); ok && err.Timeout() {
				mutex.Lock()
				countMissResp++
				mutex.Unlock()
			} else if err != nil {
				panic(err)
			} else {
				mutex.Lock()
				allTimeReq = append(allTimeReq, timeReq)
				mutex.Unlock()
			}
		}(&myWg)
	}
	myWg.Wait()

	if len(allTimeReq) == 0 {
		return out, errors.New("all requests failed")
	}

	minTime := allTimeReq[0]
	maxTime := allTimeReq[0]
	calculateMaxMinTimeResp(allTimeReq, &maxTime, &minTime)

	totalReqTime := calculateTotalReq(allTimeReq)
	averageTimeReq := totalReqTime / len(allTimeReq)
	out = OutputResult{totalReqTime, averageTimeReq,
		maxTime, minTime, countMissResp}
	return
}

func calculateTotalReq(allTimeResp []int) (totalReq int) {
	for _, el := range allTimeResp {
		totalReq += el
	}
	return
}

func init() {
	url = *flag.String("url", urlDefault, "url")
	countOfReqFlag := flag.String("count", string(countDefault), "count of request")
	timeOutFlag := flag.String("timeOut", string(timeOutDefault), "timeOut")
	flag.Parse()

	var err error
	countOfReq, err = strconv.Atoi(*countOfReqFlag)
	if err != nil {
		fmt.Printf("parameter countOfReq contain an error, using default value = %d\n", countDefault)
		countOfReq = countDefault
	}

	timeOut, err = strconv.Atoi(*timeOutFlag)
	if err != nil {
		fmt.Printf("parameter timeOut contain an error, using default value = %d\n", timeOutDefault)
		timeOut = timeOutDefault
	}
}

func calculateMaxMinTimeResp(allTimeResp []int, max, min *int) {
	for i := 0; i < len(allTimeResp); i++ {
		if allTimeResp[i] < *min {
			*min = allTimeResp[i]
		} else if allTimeResp[i] > *max {
			*max = allTimeResp[i]
		}
	}
}

func (response OutputResult) String() string {
	return fmt.Sprintf("Time during which all requests worked: %d\n"+
		"Average request time: %d\n"+
		"Maximum response time: %d\n"+
		"Minimum response time: %d\n"+
		"Number of missed responses: %d",
		response.totalTimeReq, response.averageTimeReq,
		response.maxResponseTime, response.minResponseTime, response.countMissResp)
}

func main() {
	response, err := sendReq()
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
