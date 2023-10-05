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
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/http"
)

type Login struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type AuthToken struct {
	IdToken string `json:"id_token"`
}

func GetAuthToken(url string, username string, password string, timeout time.Duration) (*string, error) {
	request, err := http.NewPostRequest(url+"/rest/login", Login{
		Username:   username,
		Password:   password,
		RememberMe: false,
	})
	if err != nil {
		return nil, err
	}
	authToken, err := http.Read[AuthToken](request, timeout, true)
	if err != nil {
		return nil, err
	}
	return &authToken.IdToken, nil
}
