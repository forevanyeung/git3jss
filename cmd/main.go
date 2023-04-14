package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var auth_bearer = ""
var auth_exp = time.Time{}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing subcommand, [config, download, sync]\n")
		os.Exit(1)
	}

	var configServer, configUsername, configPassword, configArg, typeArg string
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	configCmd.StringVar(&configServer, "server", "def", "Jamf server URL")
	configCmd.StringVar(&configUsername, "username", "", "")
	configCmd.StringVar(&configPassword, "password", "", "")
	configCmd.StringVar(&configArg, "config", "", "path to configuration file")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadCmd.StringVar(&typeArg, "type", "scripts,extensions", "comma-separated")
	downloadCmd.StringVar(&configArg, "config", "", "path to configuration file")

	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncCmd.StringVar(&typeArg, "type", "scripts,extensions", "comma-separated")
	syncCmd.StringVar(&configArg, "config", "", "path to configuration file")

	switch os.Args[1] {
	case "config":
		configCmd.Parse(os.Args[2:])
		fmt.Printf("Server: %s\nUsername:%s\n", configServer, configUsername)
	case "download":
		downloadCmd.Parse(os.Args[2:])
		items := strings.Split(typeArg, ",")
		start(configArg)

		for _, item := range items {
			switch item {
			case "scripts":
				// get scripts from Jamf
				fmt.Printf("Getting scripts from Jamf server\n")

				createBaseDirIfNotExists("scripts")

				for _, script := range getJssScripts() {
					saveJssScript(script, "scripts")
				}

			case "extensions":
				// get extension attributes from Jamf
				fmt.Printf("Getting extension attributes from Jamf server\n")
				// sendJssExtensionAttributeRequest(-1)
				// sendJssExtensionAttributeRequest(144)

			default:
				fmt.Printf("Unknown download argument: %s\n", item)
			}
		}
	case "review":
		start(configArg)

		// get list of all script
		scripts := getJssScripts()
		var names []string

		fmt.Println(longest(names))

		break

		// get list of all policies
		policies := getJssPolicies()
		for _, policy := range policies.Policies {
			// loop through policies and get each policy details (General and Scripts)
			policyDetails := getJssPolicyScript(policy)
			for i := range policyDetails.Policy.Scripts {
				// find matching script and append policy to script
				for ii := range scripts {
					if scripts[ii].Id == policyDetails.Policy.Scripts[i].Id {
						scripts[ii].JssPolicy = append(scripts[ii].JssPolicy, policyDetails.Policy)
					}
				}
			}
		}

		// fmt.Println(scriptIds)
		// fmt.Println(usedScriptsIds)

		// print scripts
		// for i := range scripts {
		// 	fmt.Printf("%-4s", scripts[i].Id)
		// 	fmt.Printf("")
		// }

	case "sync":
		//
	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}

}

func start(path string) {
	// default path
	if path == "" {
		path = "config.json"
	}

	// open config
	config, err := LoadConfiguration(path)
	if err != nil {
		log.Fatal(err)
	}

	// if authentication is expired, run authentication
	if auth_exp.Before(time.Now().UTC()) {
		fmt.Printf("Authentication has expired: %s\n", auth_exp)

		auth_bearer, auth_exp = auth(config)

		fmt.Printf("Authentication expiration is set to: %s\n", auth_exp)
	}
}

func createBaseDirIfNotExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
