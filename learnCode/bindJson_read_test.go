package bufio

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

/**
为什么BindJSON不能调用第二次，对应的测试用例Uri: /bingJson1
*/

func TestJSonDeCode(t *testing.T) {
	// 示例JSON字符串
	jsonStr := `{"name":"张三","age":30}`
	// 将字符串转换为io.Reader
	ioR := strings.NewReader(jsonStr)
	newDecoder := json.NewDecoder(ioR)
	obj := new(map[string]interface{})
	if err := newDecoder.Decode(obj); err != nil {
		t.Error(err)
	} else {
		t.Log(obj)
		//第二次的调用Decode方法,Json内部使用的是dec.readValue(),碰到具体的
		if err = newDecoder.Decode(obj); err != nil {
			t.Error(err) //json_read_test.go:22: EOF 读到了文件的末尾
		} else {
			t.Log(obj)
		}
	}
}

func TestIoReadAll(t *testing.T) {
	// 示例JSON字符串
	jsonStr := `{"name":"张三","age":30}`
	// 将字符串转换为io.Reader
	ioR := strings.NewReader(jsonStr)
	//这里可以对比一下 io.ReadAll 和 readValue
	body, err := io.ReadAll(ioR)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(body))
		body, err = io.ReadAll(ioR)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(string(body))
		}
	}
}
