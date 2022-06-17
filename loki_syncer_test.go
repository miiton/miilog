package miilog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"testing"
	"time"
)

type testParams struct {
	lokiURL  string
	msg      []byte
	tenantID string
	labels   string
}

var params = &testParams{
	lokiURL:  "http://localhost:3100",
	msg:      []byte("{\"level\":\"info\",\"msg\": \"hogefugapiyohogefugapiyohogefuga-\"}"),
	tenantID: "MYAWESOMETENANT",
	labels:   "{host=\"localhost\", job=\"gotest\"}",
}

type lokiLog struct {
	Data struct {
		Result []struct {
			Values [][]string `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func Test_GenPushRequest(t *testing.T) {
	ls := &LokiSyncer{
		Labels: params.labels,
	}

	req, err := ls.GenPushRequest(params.msg)
	if err != nil {
		t.Error(err)
	}
	if len(req) == 0 {
		t.Error("wrong request result")
	}
}

func getLokiLog(client *http.Client, start int64) (values [][]string, err error) {
	u, err := url.Parse(params.lokiURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "loki", "api", "v1", "query_range")
	q := u.Query()
	q.Set("query", `{job="gotest"}`)
	q.Add("start", fmt.Sprint(start))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return values, err
	}
	req.Header.Add("X-Scope-OrgID", params.tenantID)
	res, err := client.Do(req)
	if err != nil {
		return values, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return values, errors.New(res.Status)
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return values, err
	}
	var l lokiLog
	err = json.Unmarshal(resBytes, &l)
	if err != nil {
		return values, err
	}
	return l.Data.Result[0].Values, nil
}

func Test_Push(t *testing.T) {
	start := time.Now().Add(-60 * time.Second).UnixNano()
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   5 * time.Second,
	}
	u, err := url.Parse(params.lokiURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "loki", "api", "v1", "push")
	ls := &LokiSyncer{
		URL:      u.String(),
		TenantID: params.tenantID,
		Labels:   params.labels,
		Client:   client,
	}
	values, err := getLokiLog(client, start)
	if err != nil {
		t.Error(err)
	}
	beforeCount := len(values)
	err = ls.Push(params.msg)
	if err != nil {
		t.Error(err)
	}

	values, err = getLokiLog(client, start)
	if err != nil {
		t.Error(err)
	}
	afterCount := len(values)
	if beforeCount+1 != afterCount {
		t.Error(errors.New("not match log numbers"))
	}
}

func BenchmarkPush(b *testing.B) {
	log.Println(time.Now().UnixNano())
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   5 * time.Second,
	}
	u, err := url.Parse(params.lokiURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "loki", "api", "v1", "push")
	ls := &LokiSyncer{
		URL:      u.String(),
		TenantID: params.tenantID,
		Labels:   params.labels,
		Client:   client,
	}
	for i := 0; i < b.N; i++ {
		err = ls.Push(params.msg)
		if err != nil {
			b.Error(err)
		}
	}
}
