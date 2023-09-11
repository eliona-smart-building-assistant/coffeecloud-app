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

package eliona

import (
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

const CoffeecloudMachineAssetType = "coffeecloud_machine"
const CoffeecloudGroupAssetType = "coffeecloud_group"
const CoffeecloudRootAssetType = "coffeecloud_root"

type MachineGroup struct {
	GroupID   string `json:"groupId" eliona:"filterable"`
	GroupName string `json:"groupName" eliona:"filterable"`
	Machines  []Machine
}

type Machine struct {
	MachineID    string `json:"machineId" eliona:"filterable"`
	MachineName  string `json:"machineName" eliona:"filterable"`
	SerialNumber string `json:"serialNumber,omitempty" eliona:"filterable"`
	Firmware     int    `json:"firmware,omitempty" eliona:"filterable"`

	CubCount          int    `json:"cubCount,omitempty" eliona:"cub_count" subtype:"input"`
	EngineStatus      string `json:"engineStatus,omitempty" eliona:"engine_status,filterable" subtype:"status"`
	HoursSinceCleaned int    `json:"hourSinceCleaned,omitempty" eliona:"hours_since_cleaned" subtype:"status"`
	ErrorCode         int    `json:"errorCode,omitempty" eliona:"error_code,filterable" subtype:"status"`
	ErrorText         string `json:"errorText,omitempty" eliona:"error,filterable" subtype:"status"`
	ErrorDescription  string `json:"errorDescription,omitempty" eliona:"error_description" subtype:"status"`
}

func UpsertAsset(projectId string, uniqueIdentifier string, parentId *int32, assetType string, name string) (*int32, error) {
	assetId, err := asset.UpsertAsset(api.Asset{
		ProjectId:               projectId,
		GlobalAssetIdentifier:   uniqueIdentifier,
		Name:                    *api.NewNullableString(common.Ptr(name)),
		AssetType:               assetType,
		Description:             *api.NewNullableString(common.Ptr(fmt.Sprintf("%s (%v)", name, uniqueIdentifier))),
		ParentLocationalAssetId: *api.NewNullableInt32(parentId),
		DeviceIds: []string{
			uniqueIdentifier,
		},
	})
	if err != nil {
		return nil, err
	}
	return assetId, nil
}
