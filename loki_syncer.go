package miilog

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/grafana/loki/pkg/logproto"
)

// LokiSyncer is a struct for log forwarder to Grafana Loki
type LokiSyncer struct {
	URL       string
	TenantID  string
	Labels    string
	WaitGroup sync.WaitGroup
	Client    *http.Client
}

// Write required by zapcore.WriteSyncer
func (ls *LokiSyncer) Write(p []byte) (n int, err error) {
	ls.WaitGroup.Add(1)
	go func(p []byte) {
		defer ls.WaitGroup.Done()
		err = ls.Push(p)
		if err != nil {
			println(err)
		}
		return
	}(p)
	return len(p), nil
}

// Sync waits processing push jobs when zap.S().Sync() called. Required by zapcore.WriteSyncer.
func (ls *LokiSyncer) Sync() error {
	ls.WaitGroup.Wait()
	return nil
}

// Push sends messages to Grafana Loki using protocol buffers.
func (ls *LokiSyncer) Push(message []byte) error {
	reqData, err := ls.GenPushRequest(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", ls.URL, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	if ls.TenantID != "" {
		req.Header.Set("X-Scope-OrgID", ls.TenantID)
	}
	res, err := ls.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	if err != nil {
		return err
	}

	return nil
}

// GenPushRequest convert message to protocol buffers and compress to generates push request body.
func (ls *LokiSyncer) GenPushRequest(msg []byte) ([]byte, error) {
	var entries []logproto.Entry
	e := logproto.Entry{Timestamp: time.Now(), Line: string(msg)}
	entries = append(entries, e)
	var streams []logproto.Stream
	stream := &logproto.Stream{
		Labels:  ls.Labels,
		Entries: entries,
	}
	streams = append(streams, *stream)
	pushReq := &logproto.PushRequest{
		Streams: streams,
	}
	buf, err := proto.Marshal(pushReq)
	if err != nil {
		return nil, err
	}
	buf = snappy.Encode(nil, buf)
	return buf, nil
}
