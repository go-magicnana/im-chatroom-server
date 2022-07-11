package server

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

// Task 任务
type Task struct {
}

var taskMap = make(map[string]func())

func init() {
	taskMap["procTask"] = ProcTask
}

func ProcTask() {
	log.Println("hello world")
}

//// Init 初始化
//func Init() (obj *Task) {
//	obj = &Task{}
//	return
//}

func ClientCleanerTask() {
	myTimer := time.NewTimer(time.Second * 1) // 启动定时器

	for {
		select {
		case <-myTimer.C:
			fmt.Println("in timer")
			myTimer.Reset(time.Second * 1) // 每次使用完后需要人为重置下
		}
	}

	// 不再使用了，结束它
	myTimer.Stop()



}

// Execute 执行任务
func (obj *Task) Execute() {
	crontabList := make(map[string]string)
	// 每分钟执行一次
	// https://crontab.guru/ 检测 crontab 准确率
	//crontabList["procTask"] = "0 */1 * * * *"
	crontabList["procTask"] = "*/1 * * * * *"
	c := cron.New(cron.WithSeconds())
	log.Println(crontabList)
	for key, value := range crontabList {
		if _, ok := taskMap[key]; ok {
			_, err := c.AddFunc(value, taskMap[key])
			if err != nil {
				log.Println("Execute task AddFunc Failed, err=" + err.Error())
			}
		}
	}
	go c.Start()
	defer c.Stop()
}
