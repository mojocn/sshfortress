package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sshfortress/model"
	"time"
)

//LoginGithub github OAuth 登陆
func LoginGithub(c *gin.Context) {
	code := c.Query("code")

	gu, err := fetchGithubUser(code)
	if err != nil {
		jsonError(c, err)
		return
	}

	user := model.User{}
	data, err := user.LoginGithub(*gu.Email, gu.Login, gu.Name, gu.Bio, gu.AvatarURL, gu.Token)
	if handleError(c, err) {
		return
	}
	llogM := model.SigninLog{}
	llogM.CreatedAt = time.Now()
	llogM.ClientIp = c.ClientIP()
	llogM.UserName = data.User.Name
	llogM.UserId = data.User.Id
	llogM.Email = data.User.Email
	llogM.LoginType = "github"
	llogM.UserAgent = c.GetHeader("User-Agent")
	err = llogM.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}

//fetchGithubUser 获取github 用户信息
func fetchGithubUser(code string) (*githubUser, error) {
	client := http.Client{}
	params := fmt.Sprintf(`{"client_id":"%s","client_secret":"%s","code":"%s"}`, model.GithubClientId, model.GithubClientSecret, code)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBufferString(params))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	gt := githubToken{}
	err = json.Unmarshal(bs, &gt)
	if err != nil {
		return nil, err
	}

	//开始获取用户信息
	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", "Bearer "+gt.AccessToken)

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("using github token to fetch User Info failed with not 200 error")
	}
	defer res.Body.Close()
	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	gu := &githubUser{}
	err = json.Unmarshal(bs, gu)
	if err != nil {
		return nil, err
	}
	if gu.Email == nil {
		tEmail := fmt.Sprintf("%d@github.com", gu.ID)
		gu.Email = &tEmail
	}
	gu.Token = gt.AccessToken
	return gu, nil
}

type githubToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
type githubUser struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             *string   `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Token             string    `json:"-"`
}
