package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (

	// MF usedforward metrics
	MF MFServices

	// App metadata needed for forwarding metrics
	App APPInfo
)

// MFServices is the struct use to parse VCAP_SERVICES Variable and forward metrics
type MFServices struct {
	MFCreds []struct {
		Creds struct {
			AccessKey string `json:"access_key"`
			Endpoint  string `json:"endpoint"`
		} `json:"credentials"`
	} `json:"metrics-forwarder"`
}

// Endpoint returns the first endpoint found
func (mf *MFServices) Endpoint() string {
	for i := range mf.MFCreds {
		return mf.MFCreds[i].Creds.Endpoint
	}
	return ""
}

// AccessKey returns the first access key
func (mf *MFServices) AccessKey() string {
	for i := range mf.MFCreds {
		return mf.MFCreds[i].Creds.AccessKey
	}
	return ""
}

// APPInfo is used to parse the application info from the environemnt
type APPInfo struct {
	AppID         string `json:"application_id"`
	AppInstanceID string `json:"application_instance_id`
	AppIndex      string `json:"appIndex"`
}

// Metric is the struct emitted in the post request
type Metric struct {
	Applications []Application `json:"applications"`
}
type Application struct {
	ID        string     `json:"id"`
	Instances []Instance `json:"instances"`
}
type Instance struct {
	ID      string        `json:"id"`
	Index   string        `json:"index"`
	Metrics []MetricEntry `json:"metrics"`
}
type MetricEntry struct {
	Name  string            `json:"name"`
	Type  string            `json:"type"`
	Value int               `json:"value"`
	Unit  string            `json:"unit"`
	Tags  map[string]string `json:"tags"`
}

func parseServices() {
	vcap := os.Getenv("VCAP_SERVICES")
	app := os.Getenv("VCAP_APPLICATION")
	MF = MFServices{}
	err := json.Unmarshal([]byte(vcap), &MF)
	if err != nil {
		panic(err)
	}

	type VCAPAPP struct {
		ApplicationID string `json:"application_id"`
	}

	vAPP := VCAPAPP{}
	err = json.Unmarshal([]byte(app), &vAPP)
	if err != nil {
		panic(err)
	}
	App = APPInfo{vAPP.ApplicationID, os.Getenv("CF_INSTANCE_GUID"), os.Getenv("CF_INSTANCE_INDEX")}
}

func startMetrics() {
	count := 0
	for {
		count++
		time.Sleep(30 * time.Second)
		fmt.Println("Forwarding metrics")
		tags := make(map[string]string)
		tags["severity"] = "1"
		metric := Metric{
			Applications: []Application{
				Application{
					ID: App.AppID,
					Instances: []Instance{
						Instance{
							ID:    App.AppInstanceID,
							Index: App.AppIndex,
							Metrics: []MetricEntry{
								MetricEntry{
									Name:  "mf-sample-app",
									Type:  "counter",
									Value: count,
									Unit:  "number",
									Tags:  tags,
								},
							},
						},
					},
				},
			},
		}
		b, err := json.Marshal(metric)
		fmt.Printf("%s\n", b)
		if err != nil {
			fmt.Printf("failed to marshal metric: %s\n", err)
			continue
		}
		req, err := http.NewRequest("POST", MF.Endpoint(), bytes.NewBuffer(b))
		if err != nil {
			fmt.Printf("building metric request failed: %s", err)
			continue
		}

		req.Header.Add("Authorization", MF.AccessKey())
		req.Header.Add("Content-Type", "application/json")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request failed: %s", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Reading response body failed: %s\n", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			fmt.Printf("Response: %d\nBody: %s", resp.StatusCode, body)
			continue
		}
		fmt.Println("Successfully Forwarded metrics")
	}
}

func main() {
	parseServices()
	startMetrics()
}
