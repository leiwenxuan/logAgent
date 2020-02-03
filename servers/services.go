package servers

import (
	"encoding/json"
	"strings"
)

const (
	// 任务路劲
	JOB_SAVE_DIR = "/server/log/"

	// 保存任务事件
	JOB_EVENT_SAVE = 1

	// 删除任务事件
	JOB_EVENT_DELETE = 2
)

type JobEvent struct {
	EventType int // 添加修改 删除
	Job       *LogEntry
}

func BuildJobEvent(eventType int, logjob *LogEntry) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       logjob,
	}
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func UnpackJob(value []byte) (ret *LogEntry, err error) {
	var (
		log *LogEntry
	)
	log = &LogEntry{}
	if err = json.Unmarshal(value, log); err != nil {
		return
	}
	ret = log
	return
}

func BuildJobScheduleExecutor(entry *LogEntry) (tailJob *TailJobMgr) {
	tailJob = NewTailJob(entry.LogPath, entry.Topic)
	return tailJob
}
