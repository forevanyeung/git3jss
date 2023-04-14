package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type JssPolicyResponse struct {
	Policies []JssPolicy
	Policy   JssPolicy
}

type JssPolicy struct {
	Id      int
	Name    string
	Enabled bool
	Scripts []JssScript
}

func getJssPolicies() JssPolicyResponse {
	url := url.URL{
		Scheme: "https",
		Host:   "clear.jamfcloud.com",
		Path:   "JSSResource/policies",
	}

	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", auth_bearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response JssPolicyResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	return response
}

func getJssPolicyScript(policy JssPolicy) JssPolicyResponse {
	url := url.URL{
		Scheme: "https",
		Host:   "clear.jamfcloud.com",
		// Path:   "JSSResource/policies/id/{id}/subset/Scripts",
	}
	url.Path = fmt.Sprintf("JSSResource/policies/id/%d/subset/General&Scripts", policy.Id)
	fmt.Printf("Url: %s\n", url.String())

	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", auth_bearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))

	var response JssPolicyResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	return response
}
