package main

import (
	"context"
	"fmt"
	"time"
)

func (j *Job) test1() {
	select {
	case <-j.ctx.Done():
		fmt.Println("test1, exit")
		return
	}
}

func (j *Job) test2() {
	select {
	case <-j.ctx.Done():
		fmt.Println("test2, exit")
		return
	}
}

func (j *Job) test3() {
	select {
	case <-j.ctx.Done():
		fmt.Println("test3 exit")
		return
	}
}

type Job struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

var Gjob *Job

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	job := &Job{
		ctx:        ctx,
		cancelFunc: cancel,
	}
	go job.test1()
	go job.test2()
	go job.test3()
	time.Sleep(2 * time.Second)
	job.cancelFunc()
	time.Sleep(1 * time.Second)

}
