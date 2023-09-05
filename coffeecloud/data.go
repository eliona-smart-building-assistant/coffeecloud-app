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

type CoffeeMachineMeta struct {
	GroupID uint            `json:"groupId"`
	Result  []CoffeeMachine `json:"result"`
}

type CoffeeMachine struct {
	ID          string       `json:"id"`
	MachineName string       `json:"machineName"`
	Origin      CoffeeOrigin `json:"origin"`
}

type CoffeeOrigin struct {
	Serial string `json:"sn"`
}

type Body struct {
	Count    bool     `json:"count"`
	Criteria struct{} `json:"criteria"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Sort     struct{} `json:"sort"`
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

func GetMachines(url string, apiKey string, token string, groupId uint, timeout time.Duration) ([]CoffeeMachine, error) {
	request, err := http.NewPostRequestWithHeaders(url+"/rest/overview/data?groupid="+strconv.Itoa(int(groupId)),
		Body{
			Count:  false,
			Limit:  10000,
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
	meta, err := http.Read[CoffeeMachineMeta](request, timeout, true)
	if err != nil {
		return nil, err
	}
	return meta.Result, nil
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
