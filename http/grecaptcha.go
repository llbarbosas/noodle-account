package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type siteVerifyResponse struct {
	Success     bool     `json:"success"`
	Score       float32  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS uint32   `json:"challenge-ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

const (
	gReCAPTCHAURL = "https://www.google.com/recaptcha/api/siteverify"
)

var (
	errVerifyUnsuccess = errors.New("ReCAPTCHA inválido")
)

func GReCAPTCHAVerifier(token string, userip ...string) error {
	gReCAPTCHASecret := os.Getenv("RECAPTCHA_SECRET")

	if token == "" {
		return errVerifyUnsuccess
	}

	postBody := url.Values{}
	postBody.Set("secret", gReCAPTCHASecret)
	postBody.Set("response", token)
	postBody.Set("remoteip", userip[0])

	responseBody := strings.NewReader(postBody.Encode())
	resp, err := http.Post(gReCAPTCHAURL, "application/x-www-form-urlencoded", responseBody)

	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	var verifyResp siteVerifyResponse
	if err := json.Unmarshal(b, &verifyResp); err != nil {
		return err
	}

	if verifyResp.Success == false {
		return errVerifyUnsuccess
	}

	return nil
}
