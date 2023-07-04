package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachBasicAuthHeader(t *testing.T) {
	tests := []struct {
		testName               string
		disableSameOriginCheck bool
		resourcePath           string
		repositoryURL          string
		requestURL             string
		redirectStatusCode     int
		redirectURL            string
		expectedPass           bool
	}{
		{
			"Download index.yaml from repository with disableSameOriginCheck=false",
			false,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download index.yaml from repository with disableSameOriginCheck=true",
			true,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download index.yaml from repository with redirect to sub domain with disableSameOriginCheck=false",
			false,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://storage.charts.ranger.io/repository",
			true,
		},
		{
			"Download index.yaml from repository with redirect to sub domain with disableSameOriginCheck=true",
			true,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			308,
			"https://storage.charts.ranger.io/repository",
			true,
		},
		{
			"Download index.yaml from different origin url redirect with redirect with disableSameOriginCheck=false",
			false,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.io/repository",
			false,
		},
		{
			"Download index.yaml from different origin url redirect with redirect with disableSameOriginCheck=true",
			true,
			"/index.yaml",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			308,
			"https://blobstorage.io/repository",
			true,
		},
		{
			"Download chart from repository with disableSameOriginCheck=false",
			false,
			"/assets/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download chart from repository with disableSameOriginCheck=true",
			true,
			"/assets/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download chart from repository with redirect to sub domain with disableSameOriginCheck=false",
			false,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.charts.ranger.io/repository",
			true,
		},
		{
			"Download chart from repository with redirect to sub domain with disableSameOriginCheck=true",
			true,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.charts.ranger.io/repository",
			true,
		},
		{
			"Download chart from different origin url redirect with disableSameOriginCheck=false",
			false,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.io/repository",
			false,
		},
		{
			"Download chart from different origin url redirect with disableSameOriginCheck=true",
			true,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.io/repository",
			true,
		},
		{
			"Download chart from different origin url with disableSameOriginCheck=false",
			false,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://localhost.charts.io/repository",
			200,
			"",
			false,
		},
		{
			"Download chart from different origin url with disableSameOriginCheck=true",
			true,
			"/_blob/nginx-sample-1.1.0.tgz",
			"https://charts.ranger.io",
			"https://localhost.charts.io/repository",
			200,
			"",
			true,
		},
		{
			"Download icon from repository with disableSameOriginCheck=false",
			false,
			"/assets/logos/fleet.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download icon from repository with disableSameOriginCheck=true",
			true,
			"/assets/logos/fleet.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			200,
			"",
			true,
		},
		{
			"Download icon from repository with redirect to sub domain with disableSameOriginCheck=false",
			false,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.charts.ranger.io/repository",
			true,
		},
		{
			"Download icon from repository with redirect to sub domain with disableSameOriginCheck=true",
			true,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"https://blobstorage.charts.ranger.io/repository",
			true,
		},
		{
			"Download icon from different origin url redirect with disableSameOriginCheck=false",
			false,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"http://blobstorage.io/repository",
			false,
		},
		{
			"Download icon from different origin url redirect with disableSameOriginCheck=true",
			true,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			307,
			"http://blobstorage.io/repository",
			true,
		},
		{
			"Download icon from different origin url with disableSameOriginCheck=false",
			false,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://cattle.charts.io",
			200,
			"",
			false,
		},
		{
			"Download icon from different origin url with disableSameOriginCheck=true",
			true,
			"/assets/logos/istio.svg",
			"https://charts.ranger.io",
			"https://charts.cattle.io",
			200,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			parsedRequestURL, err := url.Parse(tt.requestURL + tt.resourcePath)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}
			repositoryRequest := &http.Request{
				URL:      parsedRequestURL,
				Response: nil,
			}
			if tt.redirectURL != "" {
				resp := &http.Response{
					StatusCode: tt.redirectStatusCode,
				}
				redirectRequestURL, err := url.Parse(tt.redirectURL + tt.resourcePath)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					t.FailNow()
				}
				repositoryRequest.URL = redirectRequestURL
				repositoryRequest.Response = resp
			}
			attachHeader, _ := shouldAttachBasicAuthHeader(tt.repositoryURL, tt.disableSameOriginCheck, repositoryRequest)
			assert.Equal(t, tt.expectedPass, attachHeader)
			t.Logf("Expected %v, got %v", tt.expectedPass, attachHeader)
		})
	}
}

func TestIsDomainOrSubdomain(t *testing.T) {
	tests := []struct {
		testName     string
		repoURL      string
		requestURL   string
		expectedPass bool
	}{
		{
			"exactly matching urls",
			"https://charts.ranger.io",
			"https://charts.ranger.io",
			true,
		},
		{
			"exactly matching urls with matching paths",
			"https://charts.ranger.io/path/here",
			"https://charts.ranger.io/path/here",
			true,
		},
		{
			"matching domains, but mismatch schema",
			"https://charts.ranger.io",
			"http://charts.ranger.io",
			false,
		},
		{
			"matching domains",
			"https://123.123.12.1:8443",
			"https://123.123.12.1:8443/path/here",
			true,
		},
		{
			"sub domain does not match, but domains match",
			"https://ranger.io",
			"https://assets.ranger.io",
			true,
		},
		{
			"mismatch domains",
			"https://charts.ranger.io",
			"https://other.ranger.io",
			false,
		},
		{
			"no matching urls",
			"https://ranger.com",
			"https://assets.ranger.io",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repoURL, err := url.Parse(tt.repoURL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}
			reqURL, err := url.Parse(tt.requestURL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}
			isSubDomainOrDomain := isDomainOrSubdomain(reqURL, repoURL)
			assert.Equal(t, tt.expectedPass, isSubDomainOrDomain)
			t.Logf("Expected %v, got %v", tt.expectedPass, isSubDomainOrDomain)
		})
	}
}
