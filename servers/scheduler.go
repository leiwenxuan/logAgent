package servers

import "github.com/sirupsen/logrus"

type Scheduler struct {
	logEntry          []*LogEntry
	jobExecutingTable map[string]*TailJobMgr // 正在执行的任务表
	jobEventChan      chan *JobEvent
}

var (
	GScheduler *Scheduler
)

func init() {
	GScheduler = &Scheduler{
		logEntry:          make([]*LogEntry, 100),
		jobExecutingTable: make(map[string]*TailJobMgr),
		jobEventChan:      make(chan *JobEvent, 10000),
	}
	go GScheduler.ScheduleLoop()
}

// 添加任务
func (s *Scheduler) PushJobEvent(jobEvent *JobEvent) {
	logrus.Debug("加入任务处理队列: ", " LogPath: ", jobEvent.Job.LogPath, "  Topic: ",
		jobEvent.Job.Topic, "    EventType: ", jobEvent.EventType)
	s.jobEventChan <- jobEvent
}

// 调度协程
func (s *Scheduler) ScheduleLoop() {
	for {
		select {
		case jobEvent := <-s.jobEventChan:
			// 维护任务队列
			s.HandleJobEvent(jobEvent)
		}
	}

}

func (s *Scheduler) HandleJobEvent(jobEvent *JobEvent) {
	switch jobEvent.EventType {
	case JOB_EVENT_SAVE:
		// 构造执行对象
		jobPlan, ok := s.jobExecutingTable[jobEvent.Job.logName]
		if ok {
			jobPlan.cancelFunc()
			logrus.Error("监听任务变化, 重新打开")
		}

		jobPlan = BuildJobScheduleExecutor(jobEvent.Job)

		s.jobExecutingTable[jobEvent.Job.logName] = jobPlan
		logrus.Debug("开始执行logAgent")
	case JOB_EVENT_DELETE:
		jobPlan, ok := s.jobExecutingTable[jobEvent.Job.logName]
		if ok {
			jobPlan.cancelFunc()
			delete(s.jobExecutingTable, jobEvent.Job.logName)
			logrus.Debug("删除logAgent: ", jobEvent.Job.LogPath)
		}
	}
}
