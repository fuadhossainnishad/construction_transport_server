package utils

type Metrics interface {
	IncDBRetry()
	IncDBFailure()
	IncDBSuccess()
}

type NoopMetrics struct{}

func (m *NoopMetrics) IncDBRetry()   {}
func (m *NoopMetrics) IncDBFailure() {}
func (m *NoopMetrics) IncDBSuccess() {}
