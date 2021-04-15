package timewheel

import (
	"errors"
	"sync"
	"time"
)

// 时间轮

type taskFn func()

const (
	Start byte = iota + 1
	Doing
	Stop
)

const MaxBuckets int = 1024 * 1024

var (
	IllegalBucketNum      = errors.New("illegal bucket num")
	IllegalTaskDelayError = errors.New("illegal delay time ")
	TaskKeyExistError     = errors.New("task key already exists")
	NotRunError           = errors.New("timeWheel not running")
)

type task struct {
	key      string
	delay    time.Duration
	pos      int
	circle   int
	fn       taskFn
	schedule bool
}

type bucket struct {
	list []*task
	mu   sync.Mutex
}

func (b *bucket) push(t *task) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.list = append(b.list, t)
}

type timeWheel struct {
	ticker       *time.Ticker
	tickDuration time.Duration
	bucketsNum   int
	buckets      []*bucket
	curPos       int
	keyPosMap    sync.Map
	status       byte
	begin        sync.Once
	end          sync.Once
	stopChan     chan struct{}
}

func normalizeTicksPerWheel(ticksPerWheel int) int {
	u := ticksPerWheel - 1
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	if u+1 > MaxBuckets {
		return MaxBuckets
	}
	return u + 1
}

// constructor
func New(tick time.Duration, bucketNum int) (*timeWheel, error) {
	if bucketNum <= 0 {
		return nil, IllegalBucketNum
	}

	num := normalizeTicksPerWheel(bucketNum)
	tw := &timeWheel{
		tickDuration: tick,
		bucketsNum:   num,
		buckets:      make([]*bucket, num),
		curPos:       0,
		status:       Start,
		stopChan:     make(chan struct{}),
	}

	for i := 0; i < num; i++ {
		tw.buckets[i] = &bucket{list: make([]*task, 0)}
	}
	return tw, nil
}

func (tw *timeWheel) addTask(key string, delay time.Duration, schedule bool, fn func()) error {
	if tw.status != Doing {
		return NotRunError
	}
	if _, exist := tw.keyPosMap.Load(key); exist {
		return TaskKeyExistError
	}
	if delay <= 0 {
		return IllegalTaskDelayError
	}

	if delay < tw.tickDuration {
		delay = tw.tickDuration
	}
	pos, circle := tw.getPositionAndCircle(delay)
	task := &task{
		delay:    delay,
		key:      key,
		pos:      pos,
		circle:   circle,
		fn:       fn,
		schedule: schedule,
	}
	tw.buckets[pos].push(task)
	tw.keyPosMap.Store(key, pos)
	return nil
}

func (tw *timeWheel) AddTask(key string, delay time.Duration, fn func()) error {
	return tw.addTask(key, delay, false, fn)
}

func (tw *timeWheel) AddScheduleTask(key string, delay time.Duration, fn func()) error {
	return tw.addTask(key, delay, true, fn)
}

func (tw *timeWheel) RemoveTask(key string) {

	pos, ok := tw.keyPosMap.Load(key)
	if !ok {
		return
	}
	bucket := tw.buckets[pos.(int)]
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	for i := 0; i < len(bucket.list); i++ {
		t := bucket.list[i]
		if key == t.key {
			bucket.list = append(bucket.list[:i], bucket.list[i+1:]...)
			tw.keyPosMap.Delete(key)
			return
		}
	}
}

func (tw *timeWheel) getPositionAndCircle(delay time.Duration) (pos int, circle int) {
	dd := int(delay / tw.tickDuration)
	circle = dd / tw.bucketsNum
	// 此处的tw.curPos表示该position中的任务都已经处理完毕
	pos = (tw.curPos + dd) & (tw.bucketsNum - 1)
	if circle > 0 && pos == tw.curPos {
		circle--
	}
	return
}

// start the timeWheel
func (tw *timeWheel) Start() {
	tw.begin.Do(func() {
		tw.ticker = time.NewTicker(tw.tickDuration)
		tw.status = Doing
		go tw.startTickerHandle()
	})
}

func (tw *timeWheel) Stop() {
	tw.end.Do(func() {
		tw.stopChan <- struct{}{}
		tw.status = Stop
	})
}

func (tw *timeWheel) startTickerHandle() {
	for {
		select {
		case <-tw.ticker.C:
			tw.handleTicker()
		case <-tw.stopChan:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *timeWheel) handleTicker() {
	curPos := (tw.curPos + 1) & (tw.bucketsNum - 1) //equals (tw.curPos + 1) % tw.bucketsNum
	bucket := tw.buckets[curPos]
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	tw.curPos = curPos
	k := 0
	for i := 0; i < len(bucket.list); i++ {
		task := bucket.list[i]
		if task.circle > 0 {
			task.circle--
			bucket.list[k] = task
			k++
			continue
		}
		go task.fn()
		if task.schedule {
			_, ok := tw.keyPosMap.Load(task.key)
			if !ok {
				continue
			}
			//reload
			pos, circle := tw.getPositionAndCircle(task.delay)
			task.pos = pos
			task.circle = circle
			if pos == curPos {
				bucket.list[k] = task
				k++
			} else {
				tw.buckets[pos].push(task)
				tw.keyPosMap.Store(task.key, pos)
			}
		} else {
			tw.keyPosMap.Delete(task.key)
		}
	}
	bucket.list = bucket.list[:k]
}
