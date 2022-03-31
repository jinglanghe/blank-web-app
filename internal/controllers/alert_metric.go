package controllers

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AlertMetricDao = &dao.AlertMetric{}
)

func registerMetric(rg *gin.RouterGroup) {
	ctrl := &metricController{}

	g := rg.Group("/metrics")
	g.GET("", ctrl.metrics)
}

type metricController struct {
	BaseController
}

func (m *metricController) metrics(c *gin.Context) {
	cnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aom_alert_total",
			Help: "Number of alerts.",
		},
		[]string{"alert_type"},
	)

	registry := prometheus.NewRegistry()
	gatherer := registry
	registry.MustRegister(cnt)

	alertTypes, _ := AlertMetricDao.ListTypes()
	for _, alertType := range alertTypes {
		cnt.WithLabelValues(alertType.Type).Add(float64(alertType.Count))
	}

	promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}
