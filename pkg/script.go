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
	for i := range [50]struct{}{} {
		b, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&scheduler.Task{
			ExecTime: time.Now().Add(time.Duration(i) + time.Second),
			Name:     fmt.Sprintf("create task #%d", i),
		})

		http.Post("http://localhost:8080/tasks", "application/json", bytes.NewBuffer(b))
	}

}
