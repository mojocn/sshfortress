package stat

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// The same as above, but now as a histogram, and only for the normal
	// distribution. The buckets are targeted to the parameters of the
	// normal distribution, with 20 buckets centered on the mean, each
	// half-sigma wide.
	//GaugeVecApiDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	//	Name:    "HttpDuration",
	//	Help:    "api requset 耗时单位ms",
	//	Buckets: prometheus.LinearBuckets(0, 1000, 50),
	//})
	GaugeVecApiDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiDuration",
		Help: "api耗时单位ms",
	}, []string{"WSorAPI"})
	GaugeVecApiMethod = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiCount",
		Help: "各种网络请求次数",
	}, []string{"method"})
	GaugeVecApiError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiErrorCount",
		Help: "请求api错误的次数type: api/ws",
	}, []string{"type"})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(GaugeVecApiMethod, GaugeVecApiDuration, GaugeVecApiError)
	// Add Go module build info.
	//processCollector := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace:"Sys"})
	//goCollector := prometheus.NewGoCollector()
	//prometheus.MustRegister(processCollector)
}
