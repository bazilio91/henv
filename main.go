package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Services map[string]map[string]map[string][]string
}

var config Config
var hosts string
var yamlSource string

const hostsFile = "/etc/hosts"

func main() {
	app := cli.NewApp()
	config = loadConfig()
	hosts = loadHosts()

	app.Action = func(c *cli.Context) error {
		service := c.Args().Get(0)
		environment := c.Args().Get(1)

		services := map[string]bool{}

		if len(service) == 0 {
			fmt.Printf(yamlSource)
		}

		if service == "all" {
			for k := range config.Services {
				services[k] = true
			}
		}

		for serviceName := range services {
			hosts = applyService(serviceName, environment, hosts)
		}

		ioutil.WriteFile(hostsFile, []byte(hosts), 0644)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
func applyService(service string, environment string, hostsContent string) string {
	out := ""

	if environment != "undo" {
		for serviceName, serviceEnvs := range config.Services {
			if serviceName != service {
				continue
			}

			if env, ok := serviceEnvs[environment]; ok {
				fmt.Printf("Applying env %s for %s\n", environment, serviceName)
				for ip, hosts := range env {
					for _, host := range hosts {
						out += fmt.Sprintf("%s %v\n", ip, host)
					}
				}
			}
		}
	}

	startString := fmt.Sprintf("# -- henv service %s start", service)
	endString := fmt.Sprintf("# -- henv service %s end", service)

	startIndex := strings.Index(hostsContent, startString)
	endIndex := strings.Index(hostsContent, endString)

	if len(out) != 0 {
		out = startString + "\n" + out + endString
	}

	hostsContent = strings.TrimRight(hostsContent, "\n")

	if startIndex >= 0 && endIndex >= 0 {
		fmt.Printf("Replacing %v config\n", service)

		var re = regexp.MustCompile(startString + `(\s*.*\s*)+` + endString)
		hostsContent = re.ReplaceAllString(hostsContent, out)
	} else {
		hostsContent += "\n\n" + out
	}

	return hostsContent
}

func loadConfig() Config {
	usr, _ := user.Current()
	dir := usr.HomeDir

	filename := filepath.Join(dir, ".henv.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	yamlSource = string(yamlFile)

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func loadHosts() string {
	file, err := ioutil.ReadFile(hostsFile)

	if err != nil {
		panic(err)
	}

	str := string(file)

	return str
}
