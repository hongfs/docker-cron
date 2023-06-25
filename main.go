package main

import (
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
)

var config = Config{}

type Config struct {
	List []CrontabItem `yaml:"list"`
}

type CrontabItem struct {
	Name    string `yaml:"name"`
	Cron    string `yaml:"cron"`
	Command string `yaml:"command"`
}

func init() {
	configStr := os.Getenv("CONFIG")

	if configStr != "" {
		if err := yaml.Unmarshal([]byte(configStr), &config); err != nil {
			panic(err)
		}

		return
	}

	// TODO 后期支持从文件读取

	if len(config.List) == 0 {
		panic("config is empty")
	}
}

func main() {
	c := cron.New(
		cron.WithSeconds(),                                          // 默认秒级
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), // 禁止并发执行同一个任务
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow,
		)))

	for _, item := range config.List {
		item := item

		_, err := c.AddFunc(item.Cron, func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[%s] recover: %v", item.Name, err)
				}
			}()

			log.Printf("[%s] %s", item.Name, item.Command)

			cmd := exec.Command("sh", "-c", item.Command)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			if err := cmd.Run(); err != nil {
				log.Printf("[%s] error: %v", item.Name, err)
			} else {
				log.Printf("[%s] success", item.Name)
			}

			return
		})

		if err != nil {
			log.Printf("[%s][AddFunc] error: %v", item.Name, err)
		}
	}

	c.Run()
}
