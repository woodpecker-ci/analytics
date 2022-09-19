package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	client := influxdb2.NewClient("http://localhost:8086", "my-token")
	defer client.Close()

	writeAPI := client.WriteAPIBlocking("woodpecker", "analytics")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/api/v1/meta", func(w http.ResponseWriter, r *http.Request) {
		p := influxdb2.NewPointWithMeasurement("stat").
			AddTag("type", "meta").
			AddField("version", 123).
			AddField("repos", 200).
			AddField("users", 3).
			AddField("pipelines", 123456).
			AddField("pipelines_time", 4321).
			AddField("forge", "gitea").
			AddField("server_os", "linux/amd64").
			AddField("agents", 1).
			AddField("agent_backend", "123").
			SetTime(time.Now())
		writeAPI.WritePoint(context.Background(), p)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
