package req

// ClientMetricParam 客户端指标参数
type ClientMetricParam struct {
	Metrics []struct {
		Fields    map[string]interface{} `json:"fields"`
		Name      string                 `json:"name"`
		Tags      map[string]string      `json:"tags"`
		Timestamp int64                  `json:"timestamp"`
	} `json:"metrics"`
}
