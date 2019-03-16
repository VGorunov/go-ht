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

type OutputResult struct {
	totalTimeReq    int
	averageTimeReq  int
	maxResponseTime int
	minResponseTime int
	countMissResp   int
}

type InputParameters struct {
	url        string
	countOfReq int
	timeOut    int
}

func sendReq() (out OutputResult, err error) {
	parameters := initParameters()
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var allTimeReq []int
	var countMissResp int

	for i := 0; i < parameters.countOfReq; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			client := http.Client{
				Timeout: time.Duration(parameters.timeOut),
			}
			_, err := client.Get(parameters.url)

			timeReq := int(time.Since(start).Nanoseconds())

			if err, ok := err.(net.Error); ok && err.Timeout() {
				mutex.Lock()
				countMissResp++
				mutex.Unlock()
			} else if err != nil {
				checkErr(err)
			} else {
				mutex.Lock()
				allTimeReq = append(allTimeReq, timeReq)
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()

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

func initParameters() InputParameters {
	urlFlag := flag.String("url", "https://www.google.com/", "url")
	countOfReqFlag := flag.String("count", "1", "count of request")
	timeOutFlag := flag.String("timeOut", "500000000", "timeOut")
	flag.Parse()

	url := *urlFlag
	countOfReq, err := strconv.Atoi(*countOfReqFlag)
	checkErr(err)
	timeOut, err := strconv.Atoi(*timeOutFlag)
	checkErr(err)

	return InputParameters{url: url, countOfReq: countOfReq, timeOut: timeOut}
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

func checkErr(err error) {
	if err != nil {
		panic(err)
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
	fmt.Println(response)
	checkErr(err)
}
