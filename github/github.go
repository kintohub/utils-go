package github

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kintohub/utils-go/klog"
	"github.com/valyala/fasthttp"
	"time"
)

type GithubInterface interface {
	CallGithub(endpoint, verb, authToken string) ([]byte, error)
	GenerateInstallationToken(installationId string) (string, error)
}

type github struct {
	baseUrl       string
	acceptHeader  string
	appID         string
	appPrivateKey []byte
}

func New(baseUrl, acceptHeader, appID string, appPrivateKey []byte) GithubInterface {
	g := &github{
		baseUrl:       baseUrl,
		acceptHeader:  acceptHeader,
		appID:         appID,
		appPrivateKey: appPrivateKey,
	}
	return g
}

// A function that calls the github api, uses github base url
// accepts a relative url that starts with "/" (ex: /app/installation)
// sets the default "Accept" header to the one used by github apps
// authToken is the full auth token not just the value (ex: authToken="Bearer {token}")
func (g *github) CallGithub(endpoint, verb, authToken string) ([]byte, error) {
	fullUrl := g.baseUrl + endpoint
	klog.Debugf("Full Github URL: %v, Auth Token: %v", fullUrl, authToken)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fullUrl)
	req.Header.SetMethod(verb)
	req.Header.Set("Accept", g.acceptHeader)
	req.Header.Set("Authorization", authToken)

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
func (g *github) GenerateInstallationToken(installationId string) (string, error) {
	endpoint := fmt.Sprintf("/app/installations/%v/access_tokens", installationId)

	jwtToken, err := g.generateJWTToken()
	if err != nil {
		return "", err
	}

	body, err := g.CallGithub(endpoint, "POST", fmt.Sprintf("Bearer %v", jwtToken))
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
		Issuer:    g.appID,
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
