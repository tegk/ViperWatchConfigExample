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

func (env *Env) printMessage() {
	for {
		fmt.Println("Port is:", env.config.Port)
		time.Sleep(time.Second)
	}
}

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
	env.config.Port = config.Host.Port
	fmt.Println(config.Host.Port)
}

func watchConfig(env *Env) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		readConfig(env)
	})
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in

	env := &Env{config: &Host{Address: "init", Port: 0}}
	viper.WatchConfig()

	readConfig(env)
	watchConfig(env)

	var wg = sync.WaitGroup{}
	wg.Add(1)
	go env.printMessage()
	wg.Wait()
}

