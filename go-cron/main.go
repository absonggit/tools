package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Name          string `yaml:"name"`
	Cron          string `yaml:"cron"`
	Command       string `yaml:"command"`
	Retries       int    `yaml:"retries"`
	RetryInterval int    `yaml:"retry_interval"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
}

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Tasks    []Task         `yaml:"tasks"`
}

func sendTelegramMessage(botToken string, chatID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("创建 Telegram Bot 失败: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("发送 Telegram 消息失败: %v", err)
	}

	return nil
}

func shellTask(name, command string, retries, retryInterval int, botToken string, chatID int64) {
	go func() {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		logger.Printf("[%s] 任务开始: %s", name, command)

		for i := 0; i <= retries; i++ {
			cmd := exec.Command("bash", "-c", command)
			var out, stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err == nil {
				logger.Printf("[%s] 任务完成: %s", name, out.String())
				return
			}

			logger.Printf("[%s] 任务失败: %s", name, stderr.String())
			if i < retries {
				logger.Printf("[%s] 等待 %d 秒后重试", name, retryInterval)
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				logger.Printf("[%s] 已达最大重试次数, 任务失败", name)
				err := sendTelegramMessage(botToken, chatID, fmt.Sprintf("[%s] 任务失败: %s", name, stderr.String()))
				if err != nil {
					logger.Printf("[%s] 发送 Telegram 消息失败: %v", name, err)
				}
			}
		}
	}()
}

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件出错: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件出错: %v", err)
	}

	scheduler := cron.New(cron.WithSeconds())

	for _, task := range config.Tasks {
		_, err := scheduler.AddFunc(task.Cron, func(name, command string, retries, retryInterval int, botToken string, chatID int64) func() {
			return func() { shellTask(name, command, retries, retryInterval, botToken, chatID) }
		}(task.Name, task.Command, task.Retries, task.RetryInterval, config.Telegram.BotToken, config.Telegram.ChatID))
		if err != nil {
			log.Fatalf("无法调度任务 %s: %v", task.Name, err)
		}
	}

	scheduler.Start()
	log.Println("调度器已启动")

	select {}
}
