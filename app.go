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
	"coffeecloud/apiserver"
	"coffeecloud/apiservices"
	"coffeecloud/coffeecloud"
	"coffeecloud/conf"
	"coffeecloud/eliona"
	"context"
	"fmt"
	"net/http"
	"strconv"
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

			groups, err := collectGroupedMachines(config)
			if err != nil {
				log.Error("coffeecloud", "error collection machines: %v", err)
				return
			} else {

				err = sendGroupedMachinesAndData(config, groups)
				if err != nil {
					log.Error("coffeecloud", "error sending assets and data: %v", err)
					return
				} else {
					log.Info("main", "collecting %d successful finished", *config.Id)
				}

			}

			time.Sleep(time.Second * time.Duration(config.RefreshInterval))
		}, config, *config.Id)
	}
}

func sendGroupedMachinesAndData(config apiserver.Configuration, groups []eliona.MachineGroup) error {

	if config.ProjectIDs == nil || len(*config.ProjectIDs) == 0 {
		log.Info("eliona", "No project id defined in configuration %d. No data is send to Eliona.")
		return nil
	}

	for _, projectId := range *config.ProjectIDs {

		rootAssetId, err := createAssetFirstTime(*config.Id, projectId, eliona.CoffeecloudRootAssetType, nil, eliona.CoffeecloudRootAssetType, "Coffeecloud Root")
		if err != nil {
			return fmt.Errorf("create root asset first time: %w", err)
		}

		for _, group := range groups {

			groupAssetId, err := createAssetFirstTime(*config.Id, projectId, eliona.CoffeecloudGroupAssetType+"_"+group.GroupID, &rootAssetId, eliona.CoffeecloudGroupAssetType, group.GroupName)
			if err != nil {
				return fmt.Errorf("create group asset first time: %w", err)
			}

			for _, machine := range group.Machines {

				machineAssetId, err := createAssetFirstTime(*config.Id, projectId, eliona.CoffeecloudMachineAssetType+"_"+machine.MachineID, &groupAssetId, eliona.CoffeecloudMachineAssetType, machine.MachineName)
				if err != nil {
					return fmt.Errorf("create machine asset first time: %w", err)
				}

				err = eliona.UpsertData(machineAssetId, eliona.CoffeecloudMachineAssetType, machine)
				if err != nil {
					return fmt.Errorf("upserting machine data: %w", err)
				}
			}

		}

	}

	return nil
}

func collectGroupedMachines(config apiserver.Configuration) ([]eliona.MachineGroup, error) {

	var eliGroups []eliona.MachineGroup

	ccToken, err := coffeecloud.GetAuthToken(config.Url, config.Username, config.Password, time.Duration(*config.RequestTimeout)*time.Second)
	if err != nil {
		return eliGroups, fmt.Errorf("getting access token: %w", err)
	}
	if ccToken == nil {
		return eliGroups, fmt.Errorf("no access token received: %w", err)
	}

	ccGroups, err := coffeecloud.GetGroups(config.Url, config.ApiKey, *ccToken, time.Duration(*config.RequestTimeout)*time.Second)
	if err != nil {
		return eliGroups, fmt.Errorf("getting groups: %w", err)
	}

	for _, ccGroup := range ccGroups {

		log.Debug("coffeecloud", "found group %s", ccGroup.Name)
		eliGroup := eliona.MachineGroup{
			GroupID:   strconv.Itoa(int(ccGroup.ID)),
			GroupName: ccGroup.Name,
		}

		shouldUse, err := eliona.AdheresToFilter(eliGroup, config.AssetFilter)
		if err != nil {
			return eliGroups, fmt.Errorf("filtering group %s: %w", eliGroup.GroupName, err)
		}
		if !shouldUse {
			continue
		}

		ccMachines, err := coffeecloud.GetMachines(config.Url, config.ApiKey, *ccToken, ccGroup.ID, time.Duration(*config.RequestTimeout)*time.Second)
		if err != nil {
			return eliGroups, fmt.Errorf("getting machines: %w", err)
		}
		ccMachineErrors, err := coffeecloud.GetMachineErrors(config.Url, config.ApiKey, *ccToken, ccGroup.ID, time.Duration(*config.RequestTimeout)*time.Second)
		if err != nil {
			return eliGroups, fmt.Errorf("getting machine errors: %w", err)
		}
		ccHealthStatuses, err := coffeecloud.GetHealthStatuses(config.Url, config.ApiKey, *ccToken, ccGroup.ID, time.Duration(*config.RequestTimeout)*time.Second)
		if err != nil {
			return eliGroups, fmt.Errorf("getting health statuses: %w", err)
		}

		for serialNumber, ccMachine := range ccMachines {
			log.Debug("coffeecloud", "found machine %s", ccMachine.MachineName)
			eliMachine := eliona.Machine{
				MachineID:         ccMachine.ID,
				MachineName:       ccMachine.MachineName,
				SerialNumber:      serialNumber,
				Firmware:          ccMachine.Origin.Firmware,
				CubCount:          ccMachine.NumberOfCups,
				HoursSinceCleaned: ccMachine.HoursSinceClean,
			}
			if ccMachineError, exists := ccMachineErrors[serialNumber]; exists {
				eliMachine.ErrorCode = ccMachineError.ErrorCode
				eliMachine.ErrorText = ccMachineError.Error
				eliMachine.ErrorDescription = ccMachineError.ErrorShort
			}
			if ccHealthyStatus, exists := ccHealthStatuses[serialNumber]; exists {
				eliMachine.EngineStatus = ccHealthyStatus.HealthStatus
			}

			shouldUse, err = eliona.AdheresToFilter(eliMachine, config.AssetFilter)
			if err != nil {
				return eliGroups, fmt.Errorf("filtering machine %s: %w", eliMachine.MachineName, err)
			}
			if !shouldUse {
				continue
			}

			eliGroup.Machines = append(eliGroup.Machines, eliMachine)
		}
		eliGroups = append(eliGroups, eliGroup)
	}
	return eliGroups, nil
}

func createAssetFirstTime(configId int64, projectId string, identifier string, parentId *int32, assetType string, name string) (int32, error) {
	uniqueIdentifier := assetType + "_" + identifier
	ctx := context.Background()

	// check if asset already exists in app
	assetId, err := conf.GetAssetId(ctx, configId, projectId, uniqueIdentifier)
	if err != nil {
		return 0, fmt.Errorf("get asset id for %s in app: %w", uniqueIdentifier, err)
	}

	// if not, create asset in Eliona also
	if assetId == nil {

		log.Debug("assets", "no asset id found for %s", uniqueIdentifier)
		assetId, err = eliona.UpsertAsset(projectId, uniqueIdentifier, parentId, assetType, name)
		if err != nil || assetId == nil {
			return 0, fmt.Errorf("upserting root asset %s in Eliona: %w", uniqueIdentifier, err)
		}

		err = conf.InsertAsset(ctx, configId, *assetId, projectId, uniqueIdentifier)
		if err != nil {
			return 0, fmt.Errorf("insert asset %s in app: %w", uniqueIdentifier, err)
		}
		log.Debug("assets", "asset created for %s with id %d", uniqueIdentifier, *assetId)

	} else {
		log.Debug("assets", "asset already created for %s with id %d", uniqueIdentifier, *assetId)
	}

	return *assetId, nil
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
