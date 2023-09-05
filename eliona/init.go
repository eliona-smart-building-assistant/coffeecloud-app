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
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

// InitEliona initialize the app in aliona
func InitEliona(connection db.Connection) error {
	if err := asset.InitAssetTypeFile("eliona/coffeecloud_root.json")(connection); err != nil {
		return fmt.Errorf("init root asset type: %v", err)
	}
	if err := asset.InitAssetTypeFile("eliona/coffeecloud_group.json")(connection); err != nil {
		return fmt.Errorf("init group asset type: %v", err)
	}
	if err := asset.InitAssetTypeFile("eliona/coffeecloud_brewer.json")(connection); err != nil {
		return fmt.Errorf("init brewer asset type: %v", err)
	}
	return nil
}
