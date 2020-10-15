package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/starrybarry/schedule/pkg/scheduler"
)

func main() {

	b, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&scheduler.Task{
		ExecTime: time.Now(),
		Name:     "create",
	})

	r, _ := http.Post("http://localhost:8080/tasks", "application/json", bytes.NewBuffer(b))

	var k interface{}
	jsoniter.NewDecoder(r.Body).Decode(&k)

	fmt.Println(k)
}
