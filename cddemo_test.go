package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gorest/types"
)

func TestDemo(t *testing.T) {
	server := newServer()
	go server.Run("0.0.0.0:1234")
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Post("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters", "application/json", bytes.NewBufferString(fmt.Sprintf("{\"name\":\"bar\"}")))
	ut.Equal(t, err, nil)

	var cluster Cluster
	err = parseBody(resp, &cluster)
	ut.Equal(t, err, nil)
	ut.Equal(t, cluster.Name, "bar")

	resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters/%s", cluster.GetID()))
	ut.Equal(t, err, nil)

	var getcluster Cluster
	err = parseBody(resp, &getcluster)
	ut.Equal(t, err, nil)
	ut.Equal(t, getcluster.GetID(), cluster.GetID())
	ut.Equal(t, getcluster.Name, "bar")

	resp, err = http.Get("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters")
	ut.Equal(t, err, nil)
	var clusterCollection types.Collection
	err = parseBody(resp, &clusterCollection)
	ut.Equal(t, err, nil)
	ut.Equal(t, reflect.TypeOf(clusterCollection.Data).Kind(), reflect.Slice)
	ut.Equal(t, reflect.ValueOf(clusterCollection.Data).Len(), 1)
	data := reflect.ValueOf(clusterCollection.Data).Index(0).Elem()
	name := data.MapIndex(reflect.ValueOf("name"))
	ut.Equal(t, fmt.Sprintf("%s", name), "bar")

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters/%s", cluster.GetID()), bytes.NewBufferString(fmt.Sprintf("{\"name\":\"dar\"}")))
	ut.Equal(t, err, nil)

	client := &http.Client{}
	resp, err = client.Do(req)
	ut.Equal(t, err, nil)

	var putcluster Cluster
	err = parseBody(resp, &putcluster)
	ut.Equal(t, err, nil)
	ut.Equal(t, putcluster.Name, "dar")

	req, err = http.NewRequest("DELETE", fmt.Sprintf("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters/%s", cluster.GetID()), nil)
	ut.Equal(t, err, nil)

	resp, err = client.Do(req)
	ut.Equal(t, err, nil)

	ut.Equal(t, resp.StatusCode, 204)

	resp, err = http.Get("http://127.0.0.1:1234/apis/zcloud.example/v1/clusters")
	ut.Equal(t, err, nil)
	var emptyclusterCollection types.Collection
	err = parseBody(resp, &emptyclusterCollection)
	ut.Equal(t, err, nil)
	ut.Equal(t, reflect.TypeOf(emptyclusterCollection.Data).Kind(), reflect.Slice)
	ut.Equal(t, reflect.ValueOf(emptyclusterCollection.Data).Len(), 0)
}

func parseBody(resp *http.Response, result interface{}) error {
	reqBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(reqBody, result)
	if err != nil {
		return err
	}

	return nil
}
