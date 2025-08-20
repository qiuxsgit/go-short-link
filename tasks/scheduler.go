package tasks

import (
	"log"
	"time"

	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	config *conf.Config
	cron   *cron.Cron
	tasks  []Task
}

// Task 定时任务接口
type Task interface {
	Name() string
	Run() error
	IsEnabled() bool
	Schedule() string
}

// NewScheduler 创建一个新的定时任务调度器
func NewScheduler(config *conf.Config) *Scheduler {
	// 创建一个支持秒级别的cron调度器
	cronScheduler := cron.New(cron.WithSeconds())

	return &Scheduler{
		config: config,
		cron:   cronScheduler,
		tasks:  make([]Task, 0),
	}
}

// RegisterTask 注册一个定时任务
func (s *Scheduler) RegisterTask(task Task) {
	s.tasks = append(s.tasks, task)
}

// Start 启动定时任务调度器
func (s *Scheduler) Start() {
	for _, task := range s.tasks {
		if !task.IsEnabled() {
			log.Printf("任务 %s 已禁用，跳过", task.Name())
			continue
		}

		// 使用闭包捕获当前任务
		t := task
		_, err := s.cron.AddFunc(task.Schedule(), func() {
			log.Printf("开始执行任务: %s", t.Name())
			startTime := time.Now()

			if err := t.Run(); err != nil {
				log.Printf("任务 %s 执行失败: %v", t.Name(), err)
			} else {
				duration := time.Since(startTime)
				log.Printf("任务 %s 执行完成，耗时: %v", t.Name(), duration)
			}
		})

		if err != nil {
			log.Printf("注册任务 %s 失败: %v", task.Name(), err)
		} else {
			log.Printf("任务 %s 已注册，调度表达式: %s", task.Name(), task.Schedule())
		}
	}

	// 启动cron调度器
	s.cron.Start()
	log.Println("定时任务调度器已启动")
}

// Stop 停止定时任务调度器
func (s *Scheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done() // 等待所有任务完成
		log.Println("定时任务调度器已停止")
	}
}
