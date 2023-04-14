package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// FIXME: merge ExtensionAttributes and ExtensionAttributes into one
type JssExtensionAttributeResponse struct {
	TotalCount          int                     `json:"size"`
	ExtensionAttributes []JssExtensionAttribute `json:"computer_extension_attributes"`
	ExtensionAttribute  JssExtensionAttribute   `json:"computer_extension_attribute"`
}

// FIXME: InputType to accept multiple different types
type JssExtensionAttribute struct {
	Id               int
	Name             string
	Enabled          bool
	Description      string
	DataType         string                          `json:"data_type"`
	InputType        JssExtentionAttributeScriptType `json:"input_type"`
	InventoryDisplay string                          `json:"inventory_display"`
	ReconDisplay     string                          `json:"recon_display"`
}

type JssExtensionAttributeTextType struct {
	Type string
}

type JssExtensionAttributeMenuType struct {
	Type    string
	Choices []string
}

type JssExtentionAttributeScriptType struct {
	Type     string
	Platform string
	Script   string
}

func sendJssExtensionAttributeRequest(id int) JssExtensionAttributeResponse {
	// uses Jamf Classic API

	url := url.URL{
		Scheme: "https",
		Host:   "clear.jamfcloud.com",
		Path:   "JSSResource/computerextensionattributes",
	}
	// if id == -1, return all extension attributes
	if id >= 0 {
		url.Path = url.Path + "/id/" + strconv.Itoa(id)
	}
	fmt.Printf("Url: %s\n", url.String())

	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", auth_bearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response JssExtensionAttributeResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	// per Jamf spec, this should be returned from server
	// here we set it manually
	response.TotalCount = len(response.ExtensionAttributes)

	fmt.Println(response)

	return response
}
