package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/devops-toolkit/clusterreport/pkg/collector"
	"github.com/gorilla/mux"
)

// Server HTTP服务器
type Server struct {
	router       *mux.Router
	sysCollector *collector.BaseCollector
	analyzer     *analyzer.SystemAnalyzer
}

// NewServer 创建新的服务器
func NewServer() *Server {
	config := collector.Config{
		Interval: 60,
		Enabled:  true,
	}

	analyzerConfig := analyzer.DefaultAnalyzerConfig()

	s := &Server{
		router:       mux.NewRouter(),
		sysCollector: collector.NewBaseCollector("system", config),
		analyzer:     analyzer.NewSystemAnalyzer(analyzerConfig),
	}

	s.setupRoutes()
	return s
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// API路由
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/health", s.handleHealth).Methods("GET")
	api.HandleFunc("/metrics", s.handleMetrics).Methods("GET")
	api.HandleFunc("/analyze", s.handleAnalyze).Methods("GET")
	api.HandleFunc("/report", s.handleReport).Methods("GET")

	// 静态文件服务
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

// handleHealth 健康检查
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleMetrics 获取系统指标
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	data, err := s.sysCollector.Collect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleAnalyze 分析系统状态
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	// 采集数据
	data, err := s.sysCollector.Collect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 分析数据
	result, err := s.analyzer.Analyze(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleReport 生成完整报告
func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.sysCollector.GenerateReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果报告包含数据，进行分析
	if report.Data != nil {
		analysisResult, err := s.analyzer.Analyze(report.Data)
		if err == nil {
			// 将分析结果添加到报告中
			reportWithAnalysis := map[string]interface{}{
				"collection": report,
				"analysis":   analysisResult,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reportWithAnalysis)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	log.Printf("启动服务器: %s", addr)
	log.Printf("健康检查: http://%s/api/v1/health", addr)
	log.Printf("系统指标: http://%s/api/v1/metrics", addr)
	log.Printf("系统分析: http://%s/api/v1/analyze", addr)
	log.Printf("完整报告: http://%s/api/v1/report", addr)
	return http.ListenAndServe(addr, s.router)
}

func main() {
	server := NewServer()

	// 启动服务器
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
