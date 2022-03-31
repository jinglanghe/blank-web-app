package controllers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	config "github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/sdk/go-utils/logging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
)

func registerPrometheus(rg *gin.RouterGroup) {
	ctrl := &prometheusController{}

	g := rg.Group("/prometheus")
	g.GET("/query-range", ctrl.queryRange)
}

type prometheusController struct {
	BaseController
}

type ResData struct {
	Code   int64       `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

func (p *prometheusController) queryRange(c *gin.Context) {
	director := func(req *http.Request) {
		_ = req.ParseForm()
		_type := req.Form.Get("type")
		query := "aom_alert_total{alert_type=\"" + _type + "\"}"
		req.Form.Add("query", query)
		req.URL.RawQuery = req.Form.Encode()

		req.URL.Host = config.Config.Prometheus.ServerUrl
		req.URL.Scheme = "http"
		req.URL.Path = "/api/v1/query_range"
	}

	modifyResponse := func(resp *http.Response) error {
		defer resp.Body.Close()

		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
		defer reader.Close()

		resData := ResData{}
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&resData)
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}

		if resp.StatusCode != http.StatusOK {
			resData.Code = int64(resp.StatusCode)
			resData.Msg = resData.Error
		}

		newBody, err := json.Marshal(resData)
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}

		buf := bytes.NewBuffer(newBody)
		resp.Header["Content-Length"] = []string{fmt.Sprint(buf.Len())}
		resp.Header["Content-Encoding"] = []string{}
		resp.Body = io.NopCloser(buf)

		return nil
	}
	proxy := &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResponse}
	proxy.ServeHTTP(c.Writer, c.Request)
}
