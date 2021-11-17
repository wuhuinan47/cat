package runner

import (
	"context"
	"log"
	"sync"
	"time"

	"gtlb.zhongzefun.com/base/go/pgx"

	"github.com/codingeasygo/util/xsql"
)

const (
	RunnerPrefixPromoteSync = "promote_sync_"
)

// RunnerContext -
var RunnerContext = ContextMap{data: map[string]ContextItem{}}

// ContextItem -
type ContextItem struct {
	context.Context
	context.CancelFunc
}

// ContextMap -
type ContextMap struct {
	sync.Mutex
	data map[string]ContextItem
}

// Exec -
func (m *ContextMap) Exec(fn func(*ContextMap)) {
	m.Lock()
	fn(m)
	m.Unlock()
}

//NamedRunner will run call by delay
func NamedRunner(name string, delay time.Duration, running *bool, call func() error) {
	log.Printf("%v is starting", name)
	var finishCount = 0
	for *running {
		err := call()
		log.Println("err:", err)
		time.Sleep(delay)
		if err == nil {
			finishCount++
			log.Println("finishCount:", finishCount)
			continue
		}
		if err != pgx.ErrNoRows {
			log.Printf("%v is fail with %v", name, err)
		} else if finishCount > 0 {
			log.Printf("%v is having %v finished", name, finishCount)
		}
		finishCount = 0
		time.Sleep(delay)
	}
	log.Printf("%v is stopped", name)
}

// NamedRunnerWithCtx -
func NamedRunnerWithCtx(ctx context.Context, name string, delay time.Duration, call func() error) {
	log.Printf("%v is starting", name)
	for {
		select {
		case <-time.After(delay):
			err := call()
			if err != nil {
				log.Printf("%v is fail with %v", name, err)
			} else {
				log.Printf("%v is finished", name)
			}
		case <-ctx.Done():
			log.Printf("%v is stopped", name)
			return
		}
	}
}

// NamedSchedule will run call at schedule time
func NamedSchedule(ctx context.Context, name string, schedule int64, call func() error) {
	var getNext = func() time.Time {
		base := time.Now()
		if schedule < xsql.Time(base).Timestamp()-xsql.TimeStartOfToday().Timestamp() {
			base = base.Add(24 * time.Hour)
		}
		return time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location()).Add(time.Duration(schedule) * time.Millisecond)
	}
	log.Printf("%v is starting", name)
	after := getNext().Sub(time.Now())
	for {
		log.Printf("next %v will run at %v", name, getNext().Format("2006-01-02 15:04:05"))
		select {
		case <-time.After(after):
			err := call()
			if err != nil {
				log.Printf("%v is fail with %v", name, err)
			}
			after = getNext().Sub(time.Now())
		case <-ctx.Done():
			log.Printf("%v is stopped", name)
			return
		}
	}
}

//NamedRunner will run call by delay
func NamedRunnerWithHMS(name string, hour, minute, second int64, running *bool, call func() error) {
	log.Printf("%v is starting", name)
	var finishCount = 0
	time.Sleep(NextDiff(hour, minute, second))
	for *running {
		err := call()
		log.Println("err:", err)
		nextDiff := NextDiff(hour, minute, second)
		time.Sleep(nextDiff)
		if err == nil {
			finishCount++
			log.Println("finishCount:", finishCount)
			continue
		}
		if err != pgx.ErrNoRows {
			log.Printf("%v is fail with %v", name, err)
		} else if finishCount > 0 {
			log.Printf("%v is having %v finished", name, finishCount)
		}
		finishCount = 0
		log.Println("nextDiff:", nextDiff)
		time.Sleep(nextDiff)
	}
	log.Printf("%v is stopped", name)
}

//NamedRunner will run call by delay
func NamedRunnerWithSeconds(name string, seconds int, running *bool, call func() error) {
	log.Printf("%v is starting \n", name)
	firstDiff := time.Duration(RunPerfectTime(seconds)) * time.Second
	log.Println("first runner interval is :", firstDiff)
	time.Sleep(firstDiff)
	for *running {
		err := call()
		if err != nil {
			log.Printf("%v is fail with %v \n", name, err)
		}
		nextDiff := time.Duration(RunPerfectTime(seconds)) * time.Second
		// nextDiff := time.Duration(seconds) * time.Second
		log.Println("next runner interval is :", nextDiff)
		time.Sleep(nextDiff)
	}
	log.Printf("%v is stopped \n", name)
}

//NamedRunner will run call by delay
func NamedRunnerWithSecondsOnly(name string, seconds, extra int, running *bool, call func() error) {
	log.Printf("%v is starting \n", name)
	firstDiff := time.Duration(RunPerfectTime(seconds)) * time.Second
	log.Println("first runner interval is :", firstDiff)
	time.Sleep(firstDiff)
	for *running {
		err := call()
		if err != nil {
			log.Printf("%v is fail with %v \n", name, err)
		}
		nextDiff := time.Duration(seconds+extra) * time.Second
		// nextDiff := time.Duration(seconds) * time.Second
		log.Println("next runner interval is :", nextDiff)
		time.Sleep(nextDiff)
	}
	log.Printf("%v is stopped \n", name)
}

// 指定时间 x时x分x秒
func NextDiff(hour, minute, second int64) (seconds time.Duration) {
	log.Println("hour:", hour)
	log.Println("minute:", minute)
	log.Println("second:", second)

	sec := hour*3600 + minute*60 + second
	now := time.Now()
	daySec := now.Hour()*3600 + now.Minute()*60 + now.Second()

	var secTemp int64
	if sec > int64(daySec) {
		secTemp = sec - int64(daySec)
	} else {
		secTemp = sec + 86400 - int64(daySec)
	}

	log.Println("setTemp:", secTemp)
	seconds = time.Duration(secTemp) * time.Second
	log.Println("seconds:", seconds)

	return
}

// 指定完整时间 比如 second = 3600 表示每隔整一个小时执行一次
func RunPerfectTime(second int) (restTime int) {
	now := time.Now()
	sec := now.Minute()*60 + now.Second()
	past := sec % second
	restTime = second - past
	return
}
