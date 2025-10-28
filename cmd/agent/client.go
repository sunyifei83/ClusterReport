package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client Agent客户端
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient 创建新的客户端
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CollectionJob 采集任务
type CollectionJob struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Nodes     []string  `json:"nodes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Report 报告
type Report struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Format    string    `json:"format"`
	Content   []byte    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ReportParams 报告参数
type ReportParams struct {
	JobID   string   `json:"job_id"`
	Formats []string `json:"formats"`
	Title   string   `json:"title"`
}

// CollectData 收集数据
func (c *Client) CollectData(nodes []string) (*CollectionJob, error) {
	payload := map[string]interface{}{
		"nodes": nodes,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/collect", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var job CollectionJob
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// GenerateReport 生成报告
func (c *Client) GenerateReport(params ReportParams) (*Report, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/report", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var report Report
	if err := json.Unmarshal(body, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

// GetJobStatus 获取任务状态
func (c *Client) GetJobStatus(jobID string) (*CollectionJob, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/job/"+jobID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var job CollectionJob
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func main() {
	// Agent客户端示例
	client := NewClient("http://localhost:8080", "your-token")

	// 收集数据
	job, err := client.CollectData([]string{"node1", "node2"})
	if err != nil {
		fmt.Printf("Failed to collect data: %v\n", err)
		return
	}

	fmt.Printf("Collection job started: %s\n", job.ID)

	// 等待完成
	for {
		status, err := client.GetJobStatus(job.ID)
		if err != nil {
			fmt.Printf("Failed to get job status: %v\n", err)
			return
		}

		if status.Status == "completed" {
			fmt.Println("Collection completed")
			break
		}

		time.Sleep(5 * time.Second)
	}

	// 生成报告
	report, err := client.GenerateReport(ReportParams{
		JobID:   job.ID,
		Formats: []string{"html", "pdf"},
		Title:   "Cluster Report",
	})

	if err != nil {
		fmt.Printf("Failed to generate report: %v\n", err)
		return
	}

	fmt.Printf("Report generated: %s\n", report.ID)
}
