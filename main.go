package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/austinarmbruster/warc"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type webData struct {
	Timestamp time.Time         `json:"@timestamp"`
	URI       string            `json:"uri"`
	IP        string            `json:"ip"`
	Headers   map[string]string `json:"headers"`
	Info      string            `json:"info"`
	Content   string            `json:"content"`
}

type esbulk struct {
	Index  string
	Action string
}

func (e *esbulk) MarshalJSON() ([]byte, error) {
	index := make(map[string]interface{})
	index["_index"] = e.Index

	internalAction := "index"
	switch e.Action {
	case "index":
		internalAction = e.Action
	case "create":
		internalAction = e.Action
	case "delete":
		internalAction = e.Action
	case "update":
		internalAction = e.Action
	}

	action := make(map[string]interface{})
	action[internalAction] = index
	return json.Marshal(action)
}

func run() error {

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		err := handleWARCFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleWARCFile(file fs.FileInfo) error {
	if file.IsDir() {
		return nil
	}

	if !strings.HasSuffix(file.Name(), "warc.gz") {
		return nil
	}

	fin, err := os.Open(file.Name())
	if fin != nil {
		defer fin.Close()
	}
	if err != nil {
		return err
	}

	zr, err := gzip.NewReader(fin)
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s.ndjson", file.Name())
	fout, err := os.Create(outputFile)
	if fout != nil {
		defer fout.Close()
	}
	if err != nil {
		return err
	}

	err = warcToNDJSON(fout, zr)
	if err != nil {
		return err
	}
	return nil
}

func warcToNDJSON(w io.Writer, r io.Reader) error {
	numRecords := math.MaxInt32
	index := "sample-data"

	wr, err := warc.NewReader(io.NopCloser(r))
	if err != nil {
		return err
	}

	var rec warc.Record
	for i := 0; i < numRecords; i++ {
		rec, err = wr.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if rec.Type == warc.RecordTypeResponse {
			wd, err := convert(rec)
			if err != nil {
				return err
			}

			idx := &esbulk{
				Index: index,
			}
			jBytes, err := json.Marshal(&idx)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%s\n",string(jBytes))

			jBytes, err = json.Marshal(&wd)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%s\n",string(jBytes))
		}
	}

	return nil
}

func convert(rec warc.Record) (*webData, error) {
	var err error
	var wd webData
	wd.Headers = make(map[string]string)
	for k, v := range rec.Headers {
		switch k {
		case "WARC-Target-URI":
			wd.URI = v
		case "WARC-IP-Address":
			wd.IP = v
		case "Content-Type":
			wd.Headers[k] = v
		case "Content-Length":
			wd.Headers[k] = v
		case "WARC-Date":
			wd.Timestamp, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
		}
	}

	split := bytes.SplitN(rec.Content.Bytes(), []byte("\n\n"), 2)
	if len(split) != 2 {
		split = bytes.SplitN(rec.Content.Bytes(), []byte("\r\n\r\n"), 2)
	}

	switch len(split) {
	case 1:
		wd.Content = string(split[0])
	case 2:
		wd.Info = string(split[0])
		wd.Content = string(split[1])
	}

	return &wd, nil
}
