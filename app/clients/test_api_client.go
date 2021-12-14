package clients

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"gopkg.in/resty.v1"
	"playground/app/models"
	"strconv"
)

const (
	LimitApiClientHystrix = "LimitApiClientHystrix"
	LimitApiUrl           = "http://127.0.0.1:8082/api/fm/collect-limit/check"
)

type LimitApiClientStruct struct {
	httpClient *resty.Client
}
type LimitApiClient interface {
	Check()
}

func NewTestApiClient() *LimitApiClientStruct {

	hystrix.ConfigureCommand(LimitApiClientHystrix, hystrix.CommandConfig{
		Timeout:               100000,
		RequestVolumeThreshold: 2,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 50,//50
		SleepWindow: 10000,
	})

	return &LimitApiClientStruct{
		httpClient: resty.New(),
	}
}

func (client *LimitApiClientStruct) Check(id int) {
	output := make(chan models.LimitCheckResponse, 1)
	checkCallFunction := client.CheckCall(id,output)

	errors := hystrix.Go("LimitApiClientHystrix",
		checkCallFunction,
	// fallback function
	func(err error) error {
		fmt.Println("an error occured", err)
		return nil
	})
	fmt.Println("asd")
	select {
	case out := <-output:
		fmt.Print("success ",out)
	case err := <-errors:
		fmt.Println("an error occured", err)
	}
}

func (client *LimitApiClientStruct) CheckCall(id int,output chan models.LimitCheckResponse) func() error{
	return func() error {
		response := new(models.LimitCheckResponse)

		req := client.httpClient.R().
			SetQueryParam("xDockId", strconv.Itoa(id)).
			SetResult(response)
		_, err := req.Get(LimitApiUrl)

		if err != nil {
			return err
		}

		output <- *response

		return err
	}
}