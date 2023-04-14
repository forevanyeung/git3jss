package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type JssScriptResponse struct {
	TotalCount int
	Scripts    []JssScript `json:"results"`
}

type JssScript struct {
	Id             json.Number
	Name           string
	Info           string
	Notes          string
	Priority       string
	CategoryId     string
	CategoryName   string
	Parameter4     string
	Parameter5     string
	Parameter6     string
	Parameter7     string
	Parameter8     string
	Parameter9     string
	Parameter10    string
	Parameter11    string
	OSRequirements string
	ScriptContents string
	JssPolicy      []JssPolicy
}

func sendJssScriptReq(page int, size int) JssScriptResponse {
	url := url.URL{
		Scheme: "https",
		Host:   "clear.jamfcloud.com",
		Path:   "api/v1/scripts",
	}
	query := url.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("page-size", strconv.Itoa(size))
	url.RawQuery = query.Encode()
	fmt.Printf("Url: %s\n", url.String())

	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", auth_bearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response JssScriptResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	return response
}

func getJssScripts() []JssScript {
	page := 0
	count := 0
	sum := 1

	collect := []JssScript{}
	for count < sum {
		response := sendJssScriptReq(page, 100)

		count = count + len(response.Scripts)
		sum = response.TotalCount // set this every time in case of changes mid run

		collect = append(collect, response.Scripts...)

		// for _, s := range response.Scripts {
		// 	fmt.Printf("Id: %s\n", s.Id)
		// 	fmt.Printf("Name: %s\n", s.Name)
		// }

		page++
	}

	fmt.Printf("Retrieved %d/%d scripts\n", count, sum)

	return collect
}

func saveJssScript(script JssScript, parent string) {
	// sanitize file name
	invalidFileChars := regexp.MustCompile(`[,:*?"<>|%/\\]`)
	// whitespaceNotNearDash := regexp.MustCompile(`(?<!-)\s(?!-)`)
	whitespaceChars := regexp.MustCompile(`\s`)
	name := invalidFileChars.ReplaceAllString(script.Name, "_")
	// name = whitespaceNotNearDash.ReplaceAllString(name, "-")
	name = whitespaceChars.ReplaceAllString(name, "-")

	// file extension
	shebangMap := map[string]string{
		"#!/bin/sh":             ".sh",
		"#!/bin/bash":           ".sh",
		"#!/bin/zsh":            ".sh",
		"#!/usr/bin/env sh":     ".sh",
		"#!/usr/bin/env bash":   ".sh",
		"#!/usr/bin/env zsh":    ".sh",
		"#!/usr/bin/python":     ".py",
		"#!/usr/bin/env python": ".py",
		"#!/usr/bin/perl":       ".pl",
		"#!/usr/bin/ruby":       ".rb",
	}
	interpreter, body, _ := strings.Cut(script.ScriptContents, "\n")
	extension, ok := shebangMap[interpreter]
	if !ok {
		extension = ".sh"
		fmt.Printf("Unknown interpreter string (%s), defaulting to .sh\n", interpreter)

		// rejoin first line to body
		body = interpreter + "\n" + body
		interpreter = ""
	}
	scriptName := name + extension

	// contents
	metadata := WriteMetadata(script)
	contents := interpreter + metadata + body

	err := os.WriteFile(filepath.Join(parent, scriptName), []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
