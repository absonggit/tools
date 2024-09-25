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
	"gopkg.in/yaml.v3"
)

// 定义任务结构体
type Task struct {
	Name          string `yaml:"name"`
	Cron          string `yaml:"cron"`
	Command       string `yaml:"command"`
	Retries       int    `yaml:"retries"`       // 失败重试次数
	RetryInterval int    `yaml:"retry_interval"` // 失败重试间隔（秒）
}

// 定义配置文件结构体
type Config struct {
	Tasks []Task `yaml:"tasks"`
}

// 执行 Shell 脚本任务
func shellTask(name, command string, retries, retryInterval int) {
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
			}
		}
	}()
}

// 加载配置文件
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
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime)

	// 加载配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件出错: %v", err)
	}

	// 创建 cron 调度器
	scheduler := cron.New(cron.WithSeconds()) // 支持秒级调度

	// 遍历配置中的任务并调度
	for _, task := range config.Tasks {
		_, err := scheduler.AddFunc(task.Cron, func(name, command string, retries, retryInterval int) func() {
			return func() { shellTask(name, command, retries, retryInterval) }
		}(task.Name, task.Command, task.Retries, task.RetryInterval))
		if err != nil {
			log.Fatalf("无法调度任务 %s: %v", task.Name, err)
		}
	}

	// 启动调度器
	scheduler.Start()
	log.Println("调度器已启动")

	// 阻塞主线程
	select {}
}
