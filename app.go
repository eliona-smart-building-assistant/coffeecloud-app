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

package main

import (
	"context"
	"net/http"
	"template/apiserver"
	"template/apiservices"
	"template/coffeecloud"
	"template/conf"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func collectData() {
	configs, err := conf.GetConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "couldn't read configs from DB: %v", err)
		return
	}
	if len(configs) == 0 {
		log.Info("conf", "no configs in DB")
		return
	}

	for _, config := range configs {

		// Skip config if disabled and set inactive
		if !conf.IsConfigEnabled(config) {
			if conf.IsConfigActive(config) {
				_, err := conf.SetConfigActiveState(context.Background(), config, false)
				if err != nil {
					log.Fatal("conf", "couldn't set config active state to DB: %v", err)
					return
				}
			}
			continue
		}

		// Signals that this config is active
		if !conf.IsConfigActive(config) {
			_, err := conf.SetConfigActiveState(context.Background(), config, true)
			if err != nil {
				log.Fatal("conf", "couldn't set config active state to DB: %v", err)
				return
			}
			log.Info("conf", "collecting initialized with Configuration %d:\n"+
				"Enable: %t\n"+
				"Refresh Interval: %d\n"+
				"Request Timeout: %d\n"+
				"Active: %t\n"+
				"Project IDs: %v\n",
				*config.Id,
				*config.Enable,
				config.RefreshInterval,
				*config.RequestTimeout,
				*config.Active,
				*config.ProjectIDs)
		}

		common.RunOnceWithParam(func(config apiserver.Configuration) {
			log.Info("main", "collecting %d started", *config.Id)

			if err := collectDevices(config); err != nil {
				return // Error is handled in the method itself.
			}

			log.Info("main", "collecting %d finished", *config.Id)

			time.Sleep(time.Second * time.Duration(config.RefreshInterval))
		}, config, *config.Id)
	}
}

func collectDevices(config apiserver.Configuration) error {

	// todo: create root asset

	token, err := coffeecloud.GetAuthToken(config.Url, config.Username, config.Password, time.Duration(*config.RequestTimeout)*time.Second)
	if err != nil {
		log.Error("coffeecloud", "getting access token: %v", err)
		return err
	}
	if token == nil {
		log.Error("coffeecloud", "no access token received", err)
		return err
	}

	groups, err := coffeecloud.GetGroups(config.Url, config.ApiKey, *token, time.Duration(*config.RequestTimeout)*time.Second)
	if err != nil {
		log.Error("coffeecloud", "getting groups: %v", err)
		return err
	}

	for _, group := range groups {

		// todo: create asset for groups
		log.Debug("coffeecloud", "found group: %s", group.Name)

		machines, err := coffeecloud.GetMachines(config.Url, config.ApiKey, *token, group.ID, time.Duration(*config.RequestTimeout)*time.Second)

		for _, machine := range machines {
			if err != nil {
				log.Error("coffeecloud", "getting machine: %v", err)
				return err
			}

			// todo: create asset for machine
			log.Debug("coffeecloud", "found machine: %s", machine.MachineName)

		}

	}

	//devices, err := kontaktio.GetDevices(config)
	//if err != nil {
	//	log.Error("kontakt-io", "getting devices info: %v", err)
	//	return err
	//}
	//if err := eliona.CreateDeviceAssetsIfNecessary(config, devices); err != nil {
	//	log.Error("eliona", "creating tag assets: %v", err)
	//	return err
	//}
	//if err := eliona.UpsertDeviceData(config, devices); err != nil {
	//	log.Error("eliona", "inserting location data into Eliona: %v", err)
	//	return err
	//}
	return nil
}

// listenApi starts the API server and listen for requests
func listenApi() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionApiController(apiservices.NewVersionApiService()),
		apiserver.NewCustomizationApiController(apiservices.NewCustomizationApiService()),
	))
	log.Fatal("main", "API server: %v", err)
}
