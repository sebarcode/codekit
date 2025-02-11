package codekit

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type HttpCallOpts struct {
	Cookie         *cookiejar.Jar
	FormValues     M
	ExpectedStatus int
}

func HttpCall(url string, method string, payload []byte, headers map[string]string, callOpts *HttpCallOpts) (*http.Response, error) {
	var (
		req        *http.Request
		byteReader *bytes.Buffer
		err        error
	)
	if len(payload) > 0 {
		byteReader = bytes.NewBuffer(payload)
		req, err = http.NewRequest(method, url, byteReader)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	var opts HttpCallOpts
	if callOpts != nil {
		opts = *callOpts
	}
	return httpDo(req, opts)
}

func httpDo(req *http.Request, opts HttpCallOpts) (*http.Response, error) {
	var client *http.Client

	//-- handling cookie
	if opts.Cookie != nil {
		client = new(http.Client)
	} else {
		//-- preparing cookie jar and http client
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize cookie jar: %s", err.Error())
		}

		tjar := opts.Cookie
		if tjar != nil {
			jar = tjar
		}

		client = &http.Client{
			Jar: jar,
		}
	}

	if req.URL.Scheme == "https" {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client.Transport = tr
	}

	var resp *http.Response
	var errCall error

	if len(opts.FormValues) > 0 {
		fvs := opts.FormValues
		vs := url.Values{}
		for k, v := range fvs {
			vs.Set(k, v.(string))
		}
		resp, errCall = client.PostForm(req.URL.String(), vs)
	} else {
		resp, errCall = client.Do(req)
	}

	if errCall == nil {
		if opts.ExpectedStatus != 0 {
			if expectedStatus := opts.ExpectedStatus; expectedStatus != 0 && resp.StatusCode != expectedStatus {
				defer resp.Body.Close()
				bs, _ := io.ReadAll(resp.Body)
				return nil, fmt.Errorf("code error: %s : %v", resp.Status, string(bs))
			}
		}
	}
	return resp, errCall
}

func HttpContent(r *http.Response) []byte {
	defer r.Body.Close()
	bytes, _ := io.ReadAll(r.Body)
	return bytes
}

func HttpContentM(r *http.Response) M {
	bytes := HttpContent(r)
	obj := M{}
	_ = json.Unmarshal(bytes, &obj)
	return obj
}

func HttpContentString(r *http.Response) string {
	bytes := HttpContent(r)
	return string(bytes)
}
