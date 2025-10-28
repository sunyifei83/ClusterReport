package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/devops-toolkit/clusterreport/pkg/collector"
	"github.com/devops-toolkit/clusterreport/pkg/generator"
	"github.com/gorilla/mux"
)

// Server HTTP服务器
type Server struct {
	router    *mux.Router
	collector *collector.MultiCollector
	analyzer  *analyzer.Chain
	generator *generator.MultiFormat
}

// NewServer 创建新的服务器
func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
		collector: collector.NewMultiCollector(
			collector.NewNodeProbeCollector(),
			collector.NewPerfSnapCollector(),
		),
		analyzer: analyzer.NewChain(
			analyzer.NewConfigAnalyzer(),
			analyzer.NewPerfAnalyzer(),
			analyzer.NewAnomalyDetector(),
		),
		generator: generator.NewMultiFormat(
			generator.NewHTMLGenerator(),
			generator.NewPDFGenerator(),
			generator.NewMarkdownGenerator(),
		),
	}

	s.setupRoutes()
	return s
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// API路由
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/health", s.handleHealth).Methods("GET")
	api.HandleFunc("/collect", s.handleCollect).Methods("POST")
	api.HandleFunc("/analyze", s.handleAnalyze).Methods("POST")
	api.HandleFunc("/generate", s.handleGenerate).Methods("POST")
	api.HandleFunc("/report", s.handleReport).Methods("POST")
	api.HandleFunc("/job/{id}", s.handleJobStatus).Methods("GET")

	// 静态文件服务
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

// handleHealth 健康检查
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// CollectRequest 采集请求
type CollectRequest struct {
	Nodes []string `json:"nodes"`
}

// handleCollect 处理采集请求
func (s *Server) handleCollect(w http.ResponseWriter, r *http.Request) {
	var req CollectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 转换节点列表
	nodes := make([]collector.Node, len(req.Nodes))
	for i, host := range req.Nodes {
		nodes[i] = collector.Node{
			Name:     fmt.Sprintf("node-%d", i+1),
			Host:     host,
			Port:     22,
			Username: "root",
		}
	}

	// 执行采集
	ctx := context.Background()
	results := s.collector.CollectAll(ctx, nodes)

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// AnalyzeRequest 分析请求
type AnalyzeRequest struct {
	Data []collector.Data `json:"data"`
}

// handleAnalyze 处理分析请求
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 执行分析
	ctx := context.Background()
	report, err := s.analyzer.Process(ctx, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	Report  *analyzer.Report `json:"report"`
	Formats []string         `json:"formats"`
}

// handleGenerate 处理生成请求
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 生成报告
	ctx := context.Background()
	results, err := s.generator.GenerateAll(ctx, req.Report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"formats": req.Formats,
		"results": results,
	})
}

// ReportRequest 报告请求
type ReportRequest struct {
	Nodes   []string `json:"nodes"`
	Formats []string `json:"formats"`
	Title   string   `json:"title"`
}

// handleReport 一键生成报告
func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// 1. 采集数据
	nodes := make([]collector.Node, len(req.Nodes))
	for i, host := range req.Nodes {
		nodes[i] = collector.Node{
			Name:     fmt.Sprintf("node-%d", i+1),
			Host:     host,
			Port:     22,
			Username: "root",
		}
	}

	collectionResults := s.collector.CollectAll(ctx, nodes)

	// 提取数据
	var data []collector.Data
	for _, result := range collectionResults {
		if result.Success {
			for _, d := range result.Data {
				data = append(data, *d)
			}
		}
	}

	// 2. 分析数据
	report, err := s.analyzer.Process(ctx, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. 生成报告
	results, err := s.generator.GenerateAll(ctx, report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"title":   req.Title,
		"formats": req.Formats,
		"results": results,
		"report":  report,
	})
}

// handleJobStatus 获取任务状态
func (s *Server) handleJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	// TODO: 实现任务状态查询
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     jobID,
		"status": "completed",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

func main() {
	server := NewServer()

	// 启动服务器
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
