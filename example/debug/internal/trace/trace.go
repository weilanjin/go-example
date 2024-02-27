package trace

import (
	"context"
	"sync"
	"time"
)

type TraceLog struct {
	TraceID    string     `json:"trace_id" bson:"trace_id"`
	Span       string     `json:"span" bson:"span"`
	Tip        string     `json:"tip" bson:"tip"`
	Err        string     `json:"err" bson:"err"`
	LogData    any        `json:"log" bson:"log"`
	CreateTime *time.Time `json:"create_time" bson:"create_time"` // 创建时间
}

func (l *TraceLog) reset() {
	l.TraceID = ""
	l.Span = ""
	l.Tip = ""
	l.Err = ""
	l.LogData = nil
	l.CreateTime = nil
}

var (
	once             sync.Once
	mu               sync.Mutex
	defaultConnector Connector
)

func Use(ctx context.Context, c Connector) {
	once.Do(func() {
		pool = sync.Pool{
			New: func() any {
				return &TraceLog{}
			},
		}
	})

	if err := defaultConnector.Init(ctx); err != nil {
		defaultConnector.Logger(err)
		return
	}
	mu.Lock()
	defaultConnector = c
	mu.Unlock()
}

func Debug(traceID, span, tip string, fn func(o *TraceLog)) {
	if defaultConnector == nil || !defaultConnector.Enable() {
		return
	}
	now := time.Now()
	go func(tid, span, tip string, now *time.Time, fn func(o *TraceLog)) {
		o := Get()
		o.TraceID = tid
		o.Span = span
		o.Tip = tip
		o.CreateTime = now
		fn(o)
		if err := defaultConnector.Push(context.Background(), o); err != nil {
			defaultConnector.Logger(err)
		}
		Put(o)
	}(traceID, span, tip, &now, fn)
}
