package httpclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var Log = log.New(os.Stderr, "httpclient", log.LstdFlags)

func MustHttpPostForm(url string, v url.Values) string {
	return MustHttpPostFormWithTimeout(url, v, 0)
}

func MustHttpPostFormWithTimeout(url string, v url.Values, timeout time.Duration) string {
	r, err := HttpPostFormWithTimeout(url, v, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func HttpPostForm(url string, v url.Values) (string, error) {
	return HttpPostFormWithTimeout(url, v, 0)
}

func HttpPostFormWithTimeout(url string, v url.Values, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.PostForm(url, v)
	if err != nil {
		Log.Printf("Error %s posting to %s:%s\n", err, url, v)
		return "", err
	}
	Log.Printf("StatusCode %d posting %s:%s\n", resp.StatusCode, url, v)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
