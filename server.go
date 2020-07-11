package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)


type ActiveSeries struct {
	Series map[string]string
	Length string
}

func startServer() {
	var mutex sync.Mutex
	var totalSeries ActiveSeries
	totalSeries.Series = make(map[string]string)
	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, ts := range req.Timeseries {
			m := make(model.Metric, len(ts.Labels))
			for _, l := range ts.Labels {
				m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
			}

			mutex.Lock()
			totalSeries.Series[m.String()] = time.Now().Format("2006.01.02 15:04:05")
			mutex.Unlock()
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		totalSeries.Length = strconv.Itoa(len(totalSeries.Series))
		t := template.Must(template.New("").Parse(html))
		t.Execute(w, totalSeries)
		fmt.Fprint(w, totalSeries)
	})

	log.Fatal(http.ListenAndServe(":2233", nil))
}


var html = `
<p>{{ "Total number of Series: " }}{{ .Length}}</p>
<table>
{{range $key, $value := .Series}}  
<tr>
	<td>{{$key}}</td>
	<td>{{$value}}</td>
</tr>
{{end}}
</table>
`