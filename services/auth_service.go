package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/david-x-chen/cartracker.api/common"
)

// Authorized function
func Authorized(w http.ResponseWriter, r *http.Request) (*common.AuthUserInfo, bool, error) {
	session, err := common.OAuthStore.Get(r, "session_cookie")
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return nil, false, err
	}

	var userInfo = "{}"

	if session.Values["userInfo"] == nil {
		url := "/authorize"
		// redirect user to authorize page
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		userInfo = session.Values["userInfo"].(string)
	}

	var user *common.AuthUserInfo
	_ = json.Unmarshal([]byte(userInfo), &user)

	return user, true, nil
}
