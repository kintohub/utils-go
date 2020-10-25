package github

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kintohub/utils-go/klog"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type GithubInterface interface {
	GetListRepos(page int32, installationId, userAccessToken string) (*GithubRepositories, error)
	GetUserInformation(userAccessToken string) (*GithubUserInfo, error)
	CreateGithubAppToken(installationId string) (string, error)
	GetUserAccessToken(code string) (string, error)
}

type github struct {
	acceptHeader    string
	appId           string
	appClientId     string
	appClientSecret string
	appPrivateKey   []byte
}

var (
	BASE_URL       = "https://api.github.com"
	BASE_URL_PLAIN = "https://github.com"
)

func New(appId, appClientId, appClientSecret string, appPrivateKey []byte) GithubInterface {
	g := &github{
		acceptHeader:    "application/vnd.github.machine-man-preview+json",
		appId:           appId,
		appClientId:     appClientId,
		appClientSecret: appClientSecret,
		appPrivateKey:   appPrivateKey, // this is user when autheticating as the github app (when cloning)
	}
	return g
}

type githubUserAccessRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type githubUserAccess struct {
	AccessToken string `json:"access_token"`
}

func (g *github) GetUserAccessToken(code string) (string, error) {
	requestBody := githubUserAccessRequest{
		ClientId:     g.appClientId,
		ClientSecret: g.appClientSecret,
		Code:         code,
	}

	url := getUrl(BASE_URL_PLAIN, "/login/oauth/access_token", "")
	body, err := g.callGithub(url, "POST", "", requestBody)
	if err != nil {
		return "", fmt.Errorf("Error getting access token from github. %v", err)
	}

	githubAccess := githubUserAccess{}

	err = json.Unmarshal(body, &githubAccess)
	if err != nil {
		return "", fmt.Errorf("Error parsing github response. %v", err)
	}
	return githubAccess.AccessToken, nil
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
	url := getUrl(BASE_URL, endpoint, query)

	body, err :=
		g.callGithub(url, "GET", genAuthToken(userAccessToken), nil)
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
	url := getUrl(BASE_URL, "/user", "")
	body, err :=
		g.callGithub(url, "GET", genAuthToken(userAccessToken), nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user info from github. %v", err)
	}
	userInfo := GithubUserInfo{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("Error parsing github response for getting user info. %v", err)
	}
	return &userInfo, nil
}

func (g *github) CreateGithubAppToken(installationId string) (string, error) {
	// TODO check if existing token is still valid or not before generating a new one
	token, err := g.generateInstallationToken(installationId)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("x-access-token:%s", token), nil
}

func getUrl(baseUrl, endpoint, query string) string {
	return fmt.Sprintf("%s/%s?%s",
		baseUrl, strings.TrimPrefix(endpoint, "/"), strings.TrimPrefix(query, "?"))
}

// A function that calls the github api, uses github base url
// accepts an endpoint that starts with "/" (ex: /app/installation)
// accept query that will be set after "?" (ex: page=2&per_page=50)
// sets the default "Accept" header to the one used by github apps
// authHeaderValue  is the full auth token not just the value (ex: authHeaderValue="Bearer {token}")
func (g *github) callGithub(url, verb, authHeaderValue string, body interface{}) ([]byte, error) {
	klog.Debugf("Full Github URL: %v, Auth Token: %v", url, authHeaderValue)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(verb)
	req.Header.Set("Accept", g.acceptHeader)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", authHeaderValue)
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

// A function that returns installation access token (different from JWT token)
// That token is used in most api calls to identify the org the user belongs to
// Calls generateJWTToken() and calls a github endpoint to generate the token
// Note: whenever installation token is used in any endpoint, need to use "token" in
// "Authorization" header instead of "Bearer"
func (g *github) generateInstallationToken(installationId string) (string, error) {
	jwtToken, err := g.generateJWTToken()
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("/app/installations/%v/access_tokens", installationId)
	url := getUrl(BASE_URL, endpoint, "")

	body, err := g.callGithub(url, "POST", genAuthBearer(jwtToken), nil)
	if err != nil {
		return "", err
	}

	type installationTokenModel struct {
		Token string `json:"token"`
	}
	installationToken := installationTokenModel{}
	err = json.Unmarshal(body, &installationToken)
	if err != nil {
		return "", err
	}

	return installationToken.Token, nil
}

// A function that returns a JWT token that is valid for 10m
// It uses the github app private key & app id (stored in config)
func (g *github) generateJWTToken() (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    g.appId,
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(g.appPrivateKey))
	if err != nil {
		klog.Fatalf("can't parse github app private key: %v", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		klog.Fatalf("can't sign the JWT token with the github app private key: %v", err)
		return "", err
	}

	return signedToken, nil
}

func genAuthToken(token string) string {
	return fmt.Sprintf("token %v", token)
}
func genAuthBearer(jwtToken string) string {
	return fmt.Sprintf("Bearer %v", jwtToken)
}
