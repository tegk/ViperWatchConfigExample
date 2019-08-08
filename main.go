package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Host struct {
	Address string
	Port    int
}

type Config struct {
	Host Host
}

type Env struct {
	config *Host
}


//printMessage is a placeholder function that reads the config value periodicity
//this could be in real life for example a web server
func (env *Env) printMessage() {
	for {
		fmt.Println("Port is:", env.config.Port)
		time.Sleep(time.Second)
	}
}

//readConfig reads and is unmarshalling values that will be assigned to the env used for dependency injection
func readConfig(env *Env) {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Unable to unmarshal config")
	}
	//the value of env.config.Port get updated at runtime
	env.config.Port = config.Host.Port
}

//watchConfig get notified when the config file get changed and calls a function to reload the values
func watchConfig(env *Env) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		readConfig(env)
	})
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in

	//init values in env to avoid panic
	//the values should be checked in business logic and waited for real values to come
	env := &Env{config: &Host{Address: "init", Port: 0}}
	viper.WatchConfig()

	readConfig(env)
	watchConfig(env)

	//starting function that reads the injected values and prints them out
	//this could be for example a web server
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go env.printMessage()
	wg.Wait()
}
