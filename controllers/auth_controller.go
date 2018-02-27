package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/david-x-chen/cartracker.api/common"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Authorize is the handler to handle authentication
func Authorize(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	session, _ := common.OAuthStore.Get(r, "session_cookie")
	session.Values["state"] = state
	session.Save(r, w)

	url := common.OAuth2Cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	//fmt.Printf("%s\n", url)
	// redirect user to authorize page
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuth2Callback is handling the callback
func OAuth2Callback(w http.ResponseWriter, r *http.Request) {
	session, err := common.OAuthStore.Get(r, "session_cookie")
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
		return
	}

	// Get the code from the response
	code := r.FormValue("code")
	ctx := context.Background()

	tok, err := common.OAuth2Cfg.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	session.Values["accessToken"] = tok.AccessToken
	session.Save(r, w)

	//fmt.Printf("%v\n", tok)

	var url = "https://dyntech.solutions/connect/userinfo"
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)

	var user *common.AuthUserInfo
	_ = json.Unmarshal(contents, &user)

	userStr, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values["userInfo"] = string(userStr[:])
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
