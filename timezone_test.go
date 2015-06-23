// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// More information about Google Directions API is available on
// https://developers.google.com/maps/documentation/directions/

package maps // import "google.golang.org/maps"

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestTimezoneNevada(t *testing.T) {

	response := `{
   "dstOffset" : 0,
   "rawOffset" : -28800,
   "status" : "OK",
   "timeZoneId" : "America/Los_Angeles",
   "timeZoneName" : "Pacific Standard Time"
}`

	server := mockServer(200, response)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &TimezoneRequest{
		Location: &LatLng{
			Lat: 39.6034810,
			Lng: -119.6822510,
		},
		Timestamp: time.Unix(1331161200, 0),
	}

	resp, err := r.Get(ctx)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := TimezoneResult{
		DstOffset:    0,
		RawOffset:    -28800,
		TimeZoneID:   "America/Los_Angeles",
		TimeZoneName: "Pacific Standard Time",
	}

	if !reflect.DeepEqual(resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
		pretty.Println(resp)
		pretty.Println(correctResponse)
	}
}

func TestTimezoneLocationMissing(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &TimezoneRequest{
		Timestamp: time.Unix(1331161200, 0),
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Missing Location should return error")
	}
}

func TestTimezoneFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &TimezoneRequest{
		Location: &LatLng{Lat: 36.578581, Lng: -118.291994},
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Failing server should return error")
	}
}