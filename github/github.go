package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/kintohub/utils-go/klog"
	"github.com/valyala/fasthttp"
)

type GithubInterface interface {
	GetListRepos(page int32, installationId, userAccessToken string) (*GithubRepositories, error)
	GetUserInformation(userAccessToken string) (*GithubUserInfo, error)
	CreateGithubAppToken(installationId string) (string, error)
	GetUserAccessToken(code string) (string, error)
}

type github struct {
	appClientID     string
	appClientSecret string
}

var (
	BASE_API_URL             = "https://api.github.com"
	BASE_URL                 = "https://github.com"
	TEMP_ACCEPT_HEADER_VALUE = "application/vnd.github.machine-man-preview+json"
)

func New(appClientID, appClientSecret string) GithubInterface {
	g := &github{
		appClientID:     appClientID,
		appClientSecret: appClientSecret,
	}
	return g
}

// code is a temp string provided by github, we will send it along with the github client id and secret to get a
// user access token for that user
func (g *github) GetUserAccessToken(code string) (string, error) {
	type githubUserAccessRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}

	requestBody := githubUserAccessRequest{
		ClientID:     g.appClientID,
		ClientSecret: g.appClientSecret,
		Code:         code,
	}

	githubUrl := getUrl(BASE_URL, "/login/oauth/access_token", "")
	body, err := g.callGithub(githubUrl, "POST", "", false, requestBody)
	bodyStr := string(body)

	if err != nil {
		klog.ErrorfWithErr(err, "Github error getting access token: %v", bodyStr)
		return "", errors.New("Error getting access token from github.")
	}
	// Response is encoded as x-www-form-urlencoded so we parse it
	// as a url segment by adding `?` at the beginning
	parsedUrl, err := url.Parse("?" + bodyStr)
	parseErr := errors.New("Error parsing github response.")
	if err != nil {
		klog.Errorf("Github error parsing get token response: %v", bodyStr)
		return "", parseErr
	}

	parsedQuery, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		klog.Errorf("Github error parsing get token response: %v", bodyStr)
		return "", parseErr
	}

	accessToken := parsedQuery["access_token"]
	if len(accessToken) == 0 {
		klog.Errorf("Github error parsing response: %v", bodyStr)
		return "", parseErr
	}
	return accessToken[0], nil
}

type GithubRepository struct {
	CloneURL      string `json:"clone_url"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}

type GithubRepositories struct {
	Repositories []GithubRepository `json:"repositories"`
	TotalCount   int                `json:"total_count"`
}

func (g *github) GetListRepos(page int32, installationId, userAccessToken string) (*GithubRepositories, error) {
	query := ""
	if page != 0 {
		query = fmt.Sprintf("page=%d", page)
	}
	endpoint := fmt.Sprintf("/user/installations/%v/repositories", installationId)
	url := getUrl(BASE_API_URL, endpoint, query)

	body, err :=
		g.callGithub(url, "GET", genAuthToken(userAccessToken), true, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting repo list from github. %v", err)
	}

	githubRepos := GithubRepositories{}
	err = json.Unmarshal(body, &githubRepos)
	if err != nil {
		return nil, fmt.Errorf("Error parsing github response. %v", err)
	}
	return &githubRepos, nil
}

type GithubUserInfo struct {
	Email    string `json:"email"`
	Username string `json:"login"`
}

func (g *github) GetUserInformation(userAccessToken string) (*GithubUserInfo, error) {
	url := getUrl(BASE_API_URL, "/user", "")
	body, err :=
		g.callGithub(url, "GET", genAuthToken(userAccessToken), true, nil)
	bodyStr := string(body)
	if err != nil {
		klog.ErrorfWithErr(err, "Github error getting user info: %v", bodyStr)
		return nil, errors.New("Error getting user info from github.")
	}

	userInfo := GithubUserInfo{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		klog.ErrorfWithErr(err, "Github parsing github response for getting user info: %v", bodyStr)
		return nil, errors.New("Error parsing github response for getting user info.")
	}
	return &userInfo, nil
}

func (g *github) CreateGithubAppToken(token string) (string, error) {
	return fmt.Sprintf("x-access-token:%s", token), nil
}

func getUrl(baseUrl, endpoint, query string) string {
	return fmt.Sprintf("%s/%s?%s",
		baseUrl, strings.TrimPrefix(endpoint, "/"), strings.TrimPrefix(query, "?"))
}

// A function that calls the github api, uses github base url
// url the github url
// verb the http verb
// authHeaderValue the "Authorization" header value, the value has to be formated with "Bearer {token}" or "token {token}"
// useTempAcceptHeader some github api's require a custom "Accept" header, if this is true will set that custom header
// body when wanting to send request body (mainly used when verb is "POST")
func (g *github) callGithub(url, verb, authHeaderValue string, useTempAcceptHeader bool, body interface{}) ([]byte, error) {
	klog.Debugf("Full Github URL: %v, Auth Token: %v", url, authHeaderValue)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(verb)

	if useTempAcceptHeader {
		req.Header.Set("Accept", TEMP_ACCEPT_HEADER_VALUE)
	}

	if authHeaderValue != "" {
		req.Header.Set("Authorization", authHeaderValue)
	}
	req.Header.SetContentType("application/json")
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req.SetBody(bodyBytes)
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() > 400 {
		return nil, fmt.Errorf("github returned: %v", resp.StatusCode())
	}
	return resp.Body(), nil
}

func genAuthToken(token string) string {
	return fmt.Sprintf("token %v", token)
}

func genAuthBearer(jwtToken string) string {
	return fmt.Sprintf("Bearer %v", jwtToken)
}
