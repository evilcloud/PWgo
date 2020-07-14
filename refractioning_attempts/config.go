package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Level    string `yaml:"Level"`
	LastName string `yaml:"LastName"`
	LastPass string `yaml:"LastPass"`
}

func saveConfig(c Configuration, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func loadConfig(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}

func createInitialConfig() Configuration {
	return Configuration{
		Level:    "SFW",
		LastName: "",
		LastPass: "",
	}
}

type Worksafe struct {
	SFW    bool
	NSFW   bool
	Sailor bool
}

func swearState(state Worksafe) string {
	if state.NSFW {
		if state.Sailor {
			return "Sailor"
		} else {
			return "NSFW"
		}
	}
	return "SFW"
}

func main() {
	s := States{true, false, false}
	s.NSFW = true
	fmt.Println(s.NSFW)
	configName := "config.yaml"
	_, err := os.Stat(configName)
	if os.IsNotExist(err) {
		err := saveConfig(createInitialConfig(), configName)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("No config file found. New one created")
	}
	c, err := loadConfig(configName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(c)
}
