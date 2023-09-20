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
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

func CoffeeCloudDashboard(projectId string) (api.Dashboard, error) {
	dashboard := api.Dashboard{}
	dashboard.Name = "CoffeeCloud"
	dashboard.ProjectId = projectId
	dashboard.Widgets = []api.Widget{}

	machines, _, err := client.NewClient().AssetsAPI.
		GetAssets(client.AuthenticationContext()).
		AssetTypeName(CoffeeCloudMachineAssetType).
		ProjectId(projectId).
		Execute()
	if err != nil {
		return api.Dashboard{}, err
	}

	for _, machine := range machines {
		widget := api.Widget{
			WidgetTypeName: "coffecloud",
			AssetId:        machine.Id,
			Sequence:       nullableInt32(2),
			Details: map[string]any{
				"388": map[string]any{
					"tilesConfig": []map[string]any{
						{
							"defaultColorIndex": 3,
							"progressBar":       nil,
							"valueMapping":      [][]string{},
						},
						{
							"defaultColorIndex": 1,
							"progressBar":       nil,
							"valueMapping":      [][]string{},
						},
						{
							"defaultColorIndex": 5,
							"progressBar":       nil,
							"valueMapping": [][]string{
								{
									"0",
									"Healthy",
									"#007305",
								},
								{
									"9999999",
									"Error",
									"#9E003D",
								},
							},
						},
						{
							"defaultColorIndex": 0,
							"progressBar":       nil,
							"valueMapping":      [][]string{},
						},
					},
				},
				"size":     1,
				"timespan": 7,
			},
			Data: []api.WidgetData{
				{
					ElementSequence: nullableInt32(1),
					AssetId:         machine.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "cup_count",
						"description":         "Cups",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(1),
					AssetId:         machine.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "hours_since_cleaned",
						"description":         "Since clean",
						"key":                 "",
						"seq":                 1,
						"subtype":             "status",
					},
				},
				{
					ElementSequence: nullableInt32(1),
					AssetId:         machine.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "error_code",
						"description":         "Status",
						"key":                 "",
						"seq":                 2,
						"subtype":             "status",
					},
				},
				{
					ElementSequence: nullableInt32(1),
					AssetId:         machine.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "error",
						"description":         "Message",
						"key":                 "",
						"seq":                 3,
						"subtype":             "status",
					},
				},
			},
		}

		// add station widget to dashboard
		dashboard.Widgets = append(dashboard.Widgets, widget)
	}
	return dashboard, nil
}

func nullableInt32(val int32) api.NullableInt32 {
	return *api.NewNullableInt32(common.Ptr[int32](val))
}
