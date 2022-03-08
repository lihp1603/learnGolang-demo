package http

import (
	"context"
	"testing"
)

func TestHttpApiClient_Get(t *testing.T) {
	url := "http://127.0.0.1:8090/create"
	client := NewHttpApiClient()
	var body []byte
	resp, err := client.Get(context.Background(), url)
	if err != nil {
		t.Error(err)
	}
	//释放资源
	defer resp.Body.Close()
	t.Logf("HTTP status:%d", resp.StatusCode)
	t.Log(body)
}
