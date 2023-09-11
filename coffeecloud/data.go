//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package coffeecloud

import (
	"fmt"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"reflect"
	"strconv"
	"template/apiserver"
	"time"
)

type ExampleDevice struct {
	ID   string `eliona:"id" subtype:"info"`
	Name string `eliona:"name,filterable" subtype:"info"`
}

type CoffeeGroup struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	SerialNumbers []string `json:"serialNumbers"`
}

type Meta[T any] struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Result []T `json:"result"`
}

type CoffeeMachine struct {
	ID          string `json:"id"`
	MachineName string `json:"machineName"`
	Origin      struct {
		SerialNumber string `json:"sn"`
		Firmware     int    `json:"fw"`
	} `json:"origin"`
	NumberOfCups int `json:"numberOfCups"`
	Relay        struct {
		Location []float64 `json:"location"`
	} `json:"relay"`
	HoursSinceClean int `json:"hoursSinceClean"`
}

type MachineError struct {
	ErrorCode  int    `json:"errorCode"`
	Error      string `json:"error"`
	ErrorShort string `json:"errorShort"`
	Origin     struct {
		SerialNumber string `json:"sn"`
	} `json:"origin"`
}

type HealthMeta struct {
	MachineKPIDetails []HealthStatus `json:"machineKPIDetails"`
}

type HealthStatus struct {
	Id     string `json:"id"`
	Origin struct {
		SerialNumber string `json:"sn"`
	} `json:"origin"`
	Reason       string `json:"reason"`
	Cause        string `json:"cause"`
	HealthStatus string `json:"healthStatus"`
}

type Body struct {
	Count    bool           `json:"count"`
	Criteria struct{}       `json:"criteria"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
	Sort     map[string]any `json:"sort"`
}

func GetGroups(url string, apiKey string, token string, timeout time.Duration) ([]CoffeeGroup, error) {
	request, err := http.NewRequestWithHeaders(url+"/rest/groups", map[string]string{
		"Authorization": "Bearer " + token,
		"API-Key":       apiKey,
	})
	if err != nil {
		return nil, err
	}
	groups, err := http.Read[[]CoffeeGroup](request, timeout, true)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// todo: implement paging (max 100 per page)
func GetMachines(url string, apiKey string, token string, groupId uint, timeout time.Duration) (map[string]CoffeeMachine, error) {
	request, err := http.NewPostRequestWithHeaders(url+"/rest/overview/data?groupid="+strconv.Itoa(int(groupId)),
		Body{
			Count:  false,
			Limit:  100,
			Offset: 0,
		},
		map[string]string{
			"Authorization": "Bearer " + token,
			"API-Key":       apiKey,
		},
	)
	if err != nil {
		return nil, err
	}
	meta, err := http.Read[Meta[CoffeeMachine]](request, timeout, true)
	if err != nil {
		return nil, err
	}
	machines := make(map[string]CoffeeMachine)
	for _, machine := range meta.Result {
		serialNumber := machine.Origin.SerialNumber
		if _, exists := machines[serialNumber]; !exists {
			machines[serialNumber] = machine
		}
	}
	return machines, nil
}

// todo: implement paging (max 100 per page)
func GetMachineErrors(url string, apiKey string, token string, groupId uint, timeout time.Duration) (map[string]MachineError, error) {
	request, err := http.NewPostRequestWithHeaders(url+"/rest/dashboard/error/search?groupid="+strconv.Itoa(int(groupId)),
		Body{
			Count:  false,
			Limit:  100,
			Offset: 0,
			Sort: map[string]any{
				"timestamp.milliseconds": "desc",
			},
		},
		map[string]string{
			"Authorization": "Bearer " + token,
			"API-Key":       apiKey,
		},
	)
	if err != nil {
		return nil, err
	}
	meta, err := http.Read[Meta[MachineError]](request, timeout, true)
	if err != nil {
		return nil, err
	}
	machineErrors := make(map[string]MachineError)
	for _, machineError := range meta.Result {
		serialNumber := machineError.Origin.SerialNumber
		if _, exists := machineErrors[serialNumber]; !exists {
			machineErrors[serialNumber] = machineError
		}
	}
	return machineErrors, nil
}

func GetHealthStatuses(url string, apiKey string, token string, groupId uint, timeout time.Duration) (map[string]HealthStatus, error) {
	request, err := http.NewPostRequestWithHeaders(url+"/rest/dashboard/healthkpi?groupid="+strconv.Itoa(int(groupId)),
		Body{
			Count:  false,
			Limit:  100,
			Offset: 0,
		},
		map[string]string{
			"Authorization": "Bearer " + token,
			"API-Key":       apiKey,
		},
	)
	if err != nil {
		return nil, err
	}
	meta, err := http.Read[HealthMeta](request, timeout, true)
	if err != nil {
		return nil, err
	}
	healthStatuses := make(map[string]HealthStatus)
	for _, healthStatus := range meta.MachineKPIDetails {
		serialNumber := healthStatus.Origin.SerialNumber
		if _, exists := healthStatuses[serialNumber]; !exists {
			healthStatuses[serialNumber] = healthStatus
		}
	}
	return healthStatuses, nil
}

func GetTags(config apiserver.Configuration) ([]ExampleDevice, error) {
	return nil, nil
}

func (tag *ExampleDevice) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	f := apiFilterToCommonFilter(filter)
	fp, err := structToMap(tag)
	if err != nil {
		return false, fmt.Errorf("converting strict to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
}

func structToMap(input interface{}) (map[string]string, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
		inputType = inputType.Elem()
	}

	if inputValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	output := make(map[string]string)
	for i := 0; i < inputValue.NumField(); i++ {
		fieldValue := inputValue.Field(i)
		fieldType := inputType.Field(i)
		output[fieldType.Name] = fieldValue.String()
	}

	return output, nil
}

func apiFilterToCommonFilter(input [][]apiserver.FilterRule) [][]common.FilterRule {
	result := make([][]common.FilterRule, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = make([]common.FilterRule, len(input[i]))
		for j := 0; j < len(input[i]); j++ {
			result[i][j] = common.FilterRule{
				Parameter: input[i][j].Parameter,
				Regex:     input[i][j].Regex,
			}
		}
	}
	return result
}
