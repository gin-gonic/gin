package test_url

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const Host = "http://127.0.0.1:8080/"

func TestUrl(t *testing.T) {

	json, err := PostJsonFormServer(context.TODO(), "bingJson")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("对应的响应是:", json)

	url := []string{"ping", "someJSON", "someJSON2"}
	for _, v := range url {
		json, err := GetJsonFormServer(context.TODO(), v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(v, "对应的响应是:", json)
	}
}

func PostJsonFormServer(ctx context.Context, v string) (string, error) {
	httpUrl := Host + v
	// 发送POST请求
	var reqBody io.Reader
	resp, err := http.Post(httpUrl, "application/json", reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // 确保关闭响应体
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetJsonFormServer
func GetJsonFormServer(todo context.Context, url string) (string, error) {
	httpUrl := Host + url
	// 发送GET请求
	resp, err := http.Get(httpUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // 确保关闭响应体
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
