// anomaly_analyzer.go
package analyzers

import (
	"context"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
)

type AnomalyAnalyzer struct {
	threshold float64
}

func (a *AnomalyAnalyzer) Analyze(ctx context.Context, data []analyzer.Data) (*analyzer.Analysis, error) {
	analysis := &analyzer.Analysis{
		Type:      "anomaly",
		Timestamp: time.Now(),
		Results:   make(map[string]interface{}),
	}

	// 异常检测逻辑
	anomalies := []string{}
	for _, d := range data {
		if isAnomaly(d, a.threshold) {
			anomalies = append(anomalies, d.Node)
		}
	}

	analysis.Results["anomalies"] = anomalies
	analysis.Results["count"] = len(anomalies)

	return analysis, nil
}
