package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	version = "1.1.0"
	author  = "sunyifei83@gmail.com"
	project = "https://github.com/sunyifei83/devops-toolkit"
)

var (
	inputPath      string
	outputFile     string
	recursive      bool
	includePattern string
	excludePattern string
	pageSize       string
	orientation    string
	marginTop      string
	marginBottom   string
	marginLeft     string
	marginRight    string
	headerTitle    string
	footerText     string
	enableTOC      bool
	tocDepth       int
	cssFile        string
	verbose        bool
	showVersion    bool
	showHelp       bool
	tempDir        string
	cleanTemp      bool
	maxDepth       int
	timeout        int
	userAgent      string
	downloadImages bool
)

func init() {
	// 基础参数
	flag.StringVar(&inputPath, "input", "", "输入路径(文件、目录或URL)")
	flag.StringVar(&inputPath, "i", "", "输入路径(文件、目录或URL)的简写")
	flag.StringVar(&outputFile, "output", "output.pdf", "输出PDF文件名")
	flag.StringVar(&outputFile, "o", "output.pdf", "输出PDF文件名的简写")
	flag.BoolVar(&recursive, "recursive", true, "递归处理子目录")
	flag.BoolVar(&recursive, "r", true, "递归处理子目录的简写")

	// 文件过滤
	flag.StringVar(&includePattern, "include", "*.html,*.md,*.markdown", "包含的文件模式(逗号分隔)")
	flag.StringVar(&excludePattern, "exclude", "", "排除的文件模式(逗号分隔)")

	// 页面设置
	flag.StringVar(&pageSize, "page-size", "A4", "页面大小(A4, A3, Letter等)")
	flag.StringVar(&orientation, "orientation", "Portrait", "页面方向(Portrait或Landscape)")
	flag.StringVar(&marginTop, "margin-top", "15mm", "顶部边距")
	flag.StringVar(&marginBottom, "margin-bottom", "15mm", "底部边距")
	flag.StringVar(&marginLeft, "margin-left", "15mm", "左边距")
	flag.StringVar(&marginRight, "margin-right", "15mm", "右边距")

	// 页眉页脚
	flag.StringVar(&headerTitle, "header", "", "页眉标题")
	flag.StringVar(&footerText, "footer", "[page] / [topage]", "页脚文本")

	// 目录设置
	flag.BoolVar(&enableTOC, "toc", true, "生成目录")
	flag.IntVar(&tocDepth, "toc-depth", 3, "目录深度")

	// 样式设置
	flag.StringVar(&cssFile, "css", "", "自定义CSS文件路径")

	// Web抓取设置
	flag.IntVar(&maxDepth, "max-depth", 3, "URL爬取最大深度")
	flag.IntVar(&timeout, "timeout", 30, "HTTP请求超时时间(秒)")
	flag.StringVar(&userAgent, "user-agent", "Mozilla/5.0 (DocConverter/1.1.0)", "HTTP User-Agent")
	flag.BoolVar(&downloadImages, "download-images", true, "下载并转换网页中的图片")

	// 其他选项
	flag.BoolVar(&verbose, "verbose", false, "详细输出模式")
	flag.BoolVar(&verbose, "v", false, "详细输出模式的简写")
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息的简写")
	flag.BoolVar(&cleanTemp, "clean", true, "转换后清理临时文件")
}

func main() {
	flag.Parse()

	if showVersion {
		printVersion()
		return
	}

	if showHelp || inputPath == "" {
		printHelp()
		return
	}

	// 检查wkhtmltopdf是否安装
	if err := checkWkhtmltopdf(); err != nil {
		log.Fatal("错误: wkhtmltopdf未安装或不可用\n请安装: sudo apt-get install wkhtmltopdf 或 sudo yum install wkhtmltopdf")
	}

	// 创建临时目录
	var err error
	tempDir, err = ioutil.TempDir("", "docconverter_*")
	if err != nil {
		log.Fatal("创建临时目录失败:", err)
	}
	if cleanTemp {
		defer os.RemoveAll(tempDir)
	}

	var files []string

	// 判断输入是URL还是本地路径
	if isURL(inputPath) {
		fmt.Printf("检测到URL输入: %s\n", inputPath)
		// URL模式：可以直接使用wkhtmltopdf转换或先爬取内容
		if strings.HasSuffix(inputPath, "/") || !hasFileExtension(inputPath) {
			// 爬取网站内容
			fmt.Println("开始爬取网站内容...")
			files, err = crawlWebsite(inputPath)
			if err != nil {
				log.Fatal("爬取网站失败:", err)
			}
		} else {
			// 单个URL页面
			files = []string{inputPath}
		}
	} else {
		// 本地文件模式
		files, err = collectFiles(inputPath)
		if err != nil {
			log.Fatal("收集文件失败:", err)
		}
	}

	if len(files) == 0 {
		log.Fatal("未找到要转换的内容")
	}

	if verbose {
		fmt.Printf("找到 %d 个文件/URL待转换\n", len(files))
	}

	// 处理文件
	htmlFiles, err := processFiles(files)
	if err != nil {
		log.Fatal("处理文件失败:", err)
	}

	// 生成PDF
	if err := generatePDF(htmlFiles); err != nil {
		log.Fatal("生成PDF失败:", err)
	}

	fmt.Printf("✅ PDF生成成功: %s\n", outputFile)

	// 显示文件信息
	if info, err := os.Stat(outputFile); err == nil {
		fmt.Printf("文件大小: %.2f MB\n", float64(info.Size())/(1024*1024))
		fmt.Printf("页面数量: %d 个源文件\n", len(files))
	}
}

func printVersion() {
	fmt.Printf("DocConverter - 文档转PDF工具\n")
	fmt.Printf("版本: %s\n", version)
	fmt.Printf("作者: %s\n", author)
	fmt.Printf("项目: %s\n", project)
}

func printHelp() {
	fmt.Println("DocConverter - 文档转PDF工具")
	fmt.Println("\n用法:")
	fmt.Println("  docconverter -i <输入路径> [选项]")
	fmt.Println("\n示例:")
	fmt.Println("  # 转换单个文件")
	fmt.Println("  docconverter -i README.md -o readme.pdf")
	fmt.Println("  ")
	fmt.Println("  # 转换整个目录")
	fmt.Println("  docconverter -i ./docs -o documentation.pdf")
	fmt.Println("  ")
	fmt.Println("  # 转换网页URL")
	fmt.Println("  docconverter -i https://example.com -o website.pdf")
	fmt.Println("  ")
	fmt.Println("  # 爬取并转换整个网站")
	fmt.Println("  docconverter -i https://example.com/docs/ -o docs.pdf --max-depth 3")
	fmt.Println("  ")
	fmt.Println("  # 自定义页面设置")
	fmt.Println("  docconverter -i ./docs -o docs.pdf --page-size A3 --orientation Landscape")
	fmt.Println("  ")
	fmt.Println("  # 添加页眉页脚")
	fmt.Println("  docconverter -i ./docs --header \"项目文档\" --footer \"[page] / [topage]\"")
	fmt.Println("\n选项:")
	flag.PrintDefaults()
}

func checkWkhtmltopdf() error {
	cmd := exec.Command("wkhtmltopdf", "--version")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func hasFileExtension(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	path := u.Path
	ext := filepath.Ext(path)
	return ext != ""
}

// downloadImage 下载图片并保存到指定目录
func downloadImage(client *http.Client, imgURL string, imgDir string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	// 从URL中提取文件名
	u, err := url.Parse(imgURL)
	if err != nil {
		return "", err
	}

	// 生成文件名
	fileName := filepath.Base(u.Path)
	if fileName == "" || fileName == "/" || fileName == "." {
		// 如果没有合适的文件名，使用时间戳
		fileName = fmt.Sprintf("img_%d.jpg", time.Now().UnixNano())
	}

	// 确保文件名是唯一的
	imgPath := filepath.Join(imgDir, fileName)
	baseImgPath := imgPath
	counter := 1
	for {
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(baseImgPath)
		nameWithoutExt := strings.TrimSuffix(baseImgPath, ext)
		imgPath = fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		counter++
	}

	// 读取图片数据
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 保存图片文件
	if err := ioutil.WriteFile(imgPath, imgData, 0644); err != nil {
		return "", err
	}

	return imgPath, nil
}

func crawlWebsite(startURL string) ([]string, error) {
	visited := make(map[string]bool)
	var pages []string

	baseURL, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// 创建图片目录
	imgDir := filepath.Join(tempDir, "images")
	if downloadImages {
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			return nil, fmt.Errorf("创建图片目录失败: %v", err)
		}
	}

	// 爬取函数 - 下载页面并处理图片
	var crawl func(urlStr string, depth int) error
	crawl = func(urlStr string, depth int) error {
		if depth > maxDepth {
			return nil
		}

		// 检查是否已访问
		if visited[urlStr] {
			return nil
		}
		visited[urlStr] = true

		if verbose {
			fmt.Printf("爬取页面 [深度:%d]: %s\n", depth, urlStr)
		}

		// 获取页面
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
		}

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// 如果需要下载图片，保存HTML并处理图片
		if downloadImages {
			// 解析HTML
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			if err != nil {
				return err
			}

			// 下载并替换图片
			doc.Find("img").Each(func(i int, s *goquery.Selection) {
				src, exists := s.Attr("src")
				if !exists || src == "" {
					return
				}

				// 解析图片URL
				imgURL, err := url.Parse(src)
				if err != nil {
					return
				}

				// 转换为绝对URL
				absoluteImgURL := baseURL.ResolveReference(imgURL)

				// 下载图片
				imgPath, err := downloadImage(client, absoluteImgURL.String(), imgDir)
				if err != nil {
					if verbose {
						fmt.Printf("⚠️ 下载图片失败: %s - %v\n", src, err)
					}
					return
				}

				// 替换图片路径为本地路径
				s.SetAttr("src", imgPath)
				if verbose {
					fmt.Printf("✓ 下载图片: %s\n", src)
				}
			})

			// 保存处理后的HTML到临时文件
			processedHTML, err := doc.Html()
			if err != nil {
				return err
			}

			// 创建临时HTML文件
			tempFile := filepath.Join(tempDir, fmt.Sprintf("page_%d.html", len(pages)))
			if err := ioutil.WriteFile(tempFile, []byte(processedHTML), 0644); err != nil {
				return err
			}

			pages = append(pages, tempFile)
		} else {
			// 不下载图片，直接使用URL
			pages = append(pages, urlStr)
		}

		// 解析HTML获取链接
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			return err
		}

		// 查找所有链接
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}

			// 解析链接
			linkURL, err := url.Parse(href)
			if err != nil {
				return
			}

			// 转换为绝对URL
			absoluteURL := baseURL.ResolveReference(linkURL)

			// 只爬取同域名下的链接
			if absoluteURL.Host == baseURL.Host &&
				strings.HasPrefix(absoluteURL.Path, baseURL.Path) {
				crawl(absoluteURL.String(), depth+1)
			}
		})

		return nil
	}

	// 开始爬取
	if err := crawl(startURL, 0); err != nil {
		return nil, err
	}

	fmt.Printf("发现 %d 个页面URL\n", len(pages))
	return pages, nil
}

func collectFiles(inputPath string) ([]string, error) {
	var files []string
	includePatterns := strings.Split(includePattern, ",")
	excludePatterns := []string{}
	if excludePattern != "" {
		excludePatterns = strings.Split(excludePattern, ",")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// 单个文件
		if shouldInclude(inputPath, includePatterns, excludePatterns) {
			files = append(files, inputPath)
		}
		return files, nil
	}

	// 目录处理
	walkFunc := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !recursive && path != inputPath {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldInclude(path, includePatterns, excludePatterns) {
			files = append(files, path)
			if verbose {
				fmt.Printf("添加文件: %s\n", path)
			}
		}

		return nil
	}

	err = filepath.WalkDir(inputPath, walkFunc)
	return files, err
}

func shouldInclude(path string, includePatterns, excludePatterns []string) bool {
	fileName := filepath.Base(path)

	// 检查排除模式
	for _, pattern := range excludePatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return false
		}
	}

	// 检查包含模式
	for _, pattern := range includePatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}

	return false
}

func processFiles(files []string) ([]string, error) {
	var htmlFiles []string

	for i, file := range files {
		if verbose {
			fmt.Printf("[%d/%d] 处理: %s\n", i+1, len(files), file)
		}

		// 如果是URL，直接添加到列表
		if isURL(file) {
			htmlFiles = append(htmlFiles, file)
			continue
		}

		ext := strings.ToLower(filepath.Ext(file))

		switch ext {
		case ".md", ".markdown":
			// Markdown转HTML
			htmlFile, err := convertMarkdownToHTML(file)
			if err != nil {
				return nil, fmt.Errorf("转换Markdown失败 %s: %v", file, err)
			}
			htmlFiles = append(htmlFiles, htmlFile)

		case ".html", ".htm":
			// 直接使用HTML文件
			// 如果有CSS文件，需要注入CSS
			if cssFile != "" {
				htmlFile, err := injectCSS(file)
				if err != nil {
					return nil, fmt.Errorf("注入CSS失败 %s: %v", file, err)
				}
				htmlFiles = append(htmlFiles, htmlFile)
			} else {
				htmlFiles = append(htmlFiles, file)
			}

		default:
			log.Printf("跳过不支持的文件类型: %s", file)
		}
	}

	return htmlFiles, nil
}

func convertMarkdownToHTML(mdFile string) (string, error) {
	// 读取Markdown文件
	mdBytes, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return "", err
	}

	// 设置Markdown解析器
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Tables
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdBytes)

	// 设置HTML渲染器
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
		Title: filepath.Base(mdFile),
	}
	renderer := html.NewRenderer(opts)

	// 转换为HTML
	htmlBytes := markdown.Render(doc, renderer)

	// 创建完整的HTML文档
	fullHTML := createHTMLDocument(string(htmlBytes), filepath.Base(mdFile))

	// 保存到临时文件
	tempFile := filepath.Join(tempDir, fmt.Sprintf("%d_%s.html",
		time.Now().UnixNano(),
		strings.TrimSuffix(filepath.Base(mdFile), filepath.Ext(mdFile))))

	if err := ioutil.WriteFile(tempFile, []byte(fullHTML), 0644); err != nil {
		return "", err
	}

	return tempFile, nil
}

func createHTMLDocument(content, title string) string {
	css := getDefaultCSS()
	if cssFile != "" {
		if customCSS, err := ioutil.ReadFile(cssFile); err == nil {
			css = string(customCSS)
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
%s
    </style>
</head>
<body>
    <div class="container">
%s
    </div>
</body>
</html>`, title, css, content)
}

func getDefaultCSS() string {
	return `
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            background: white;
        }
        
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
            page-break-after: avoid;
        }
        
        h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h3 { font-size: 1.25em; }
        h4 { font-size: 1em; }
        h5 { font-size: 0.875em; }
        h6 { font-size: 0.85em; color: #6a737d; }
        
        p {
            margin-top: 0;
            margin-bottom: 16px;
        }
        
        code {
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%;
            background-color: rgba(27,31,35,0.05);
            border-radius: 3px;
            font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
        }
        
        pre {
            padding: 16px;
            overflow: auto;
            font-size: 85%;
            line-height: 1.45;
            background-color: #f6f8fa;
            border-radius: 3px;
            page-break-inside: avoid;
        }
        
        pre code {
            padding: 0;
            background-color: transparent;
        }
        
        table {
            border-spacing: 0;
            border-collapse: collapse;
            margin-bottom: 16px;
        }
        
        table th, table td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
        }
        
        table th {
            font-weight: 600;
            background-color: #f6f8fa;
        }
        
        table tr:nth-child(2n) {
            background-color: #f6f8fa;
        }
        
        blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin: 0 0 16px 0;
        }
        
        ul, ol {
            padding-left: 2em;
            margin-bottom: 16px;
        }
        
        li + li {
            margin-top: 0.25em;
        }
        
        a {
            color: #0366d6;
            text-decoration: none;
        }
        
        a:hover {
            text-decoration: underline;
        }
        
        img {
            max-width: 100%;
            box-sizing: content-box;
        }
        
        hr {
            height: 0.25em;
            padding: 0;
            margin: 24px 0;
            background-color: #e1e4e8;
            border: 0;
        }
        `
}

func injectCSS(htmlFile string) (string, error) {
	// 读取HTML文件
	htmlBytes, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		return "", err
	}

	htmlContent := string(htmlBytes)

	// 读取CSS文件
	cssContent, err := ioutil.ReadFile(cssFile)
	if err != nil {
		return "", err
	}

	// 查找</head>标签，插入CSS
	headEnd := strings.Index(htmlContent, "</head>")
	if headEnd == -1 {
		// 如果没有head标签，创建完整的HTML文档
		htmlContent = createHTMLDocument(htmlContent, filepath.Base(htmlFile))
	} else {
		// 插入CSS
		styleTag := fmt.Sprintf("<style>\n%s\n</style>\n", string(cssContent))
		htmlContent = htmlContent[:headEnd] + styleTag + htmlContent[headEnd:]
	}

	// 保存到临时文件
	tempFile := filepath.Join(tempDir, fmt.Sprintf("%d_%s",
		time.Now().UnixNano(),
		filepath.Base(htmlFile)))

	if err := ioutil.WriteFile(tempFile, []byte(htmlContent), 0644); err != nil {
		return "", err
	}

	return tempFile, nil
}

func generatePDF(htmlFiles []string) error {
	if len(htmlFiles) == 0 {
		return fmt.Errorf("没有HTML文件可转换")
	}

	// 对于URL，每个单独处理以避免崩溃
	const urlThreshold = 5 // 如果有超过5个URL，单独处理每个
	isMainlyURLs := false
	urlCount := 0
	for _, file := range htmlFiles {
		if isURL(file) {
			urlCount++
		}
	}
	if urlCount > urlThreshold {
		isMainlyURLs = true
	}

	if isMainlyURLs && len(htmlFiles) > 1 {
		fmt.Printf("⚠️ 检测到多个URL (%d)，将单独处理每个页面以确保稳定性\n", len(htmlFiles))

		// 创建临时PDF文件列表
		var tempPDFs []string
		defer func() {
			// 清理临时PDF文件
			for _, pdf := range tempPDFs {
				os.Remove(pdf)
			}
		}()

		// 单独处理每个URL
		for i, htmlFile := range htmlFiles {
			tempPDF := filepath.Join(tempDir, fmt.Sprintf("page_%d.pdf", i))

			fmt.Printf("[%d/%d] 处理页面: %s\n", i+1, len(htmlFiles), htmlFile)

			if err := generateSinglePDF([]string{htmlFile}, tempPDF, false); err != nil {
				fmt.Printf("⚠️ 页面 %d 处理失败: %v，跳过\n", i+1, err)
				continue // 跳过失败的页面，继续处理其他页面
			}

			tempPDFs = append(tempPDFs, tempPDF)

			// 短暂延迟，避免过快请求
			if i < len(htmlFiles)-1 {
				time.Sleep(500 * time.Millisecond)
			}
		}

		if len(tempPDFs) == 0 {
			return fmt.Errorf("所有页面处理失败")
		}

		// 合并所有成功的PDF
		fmt.Printf("合并 %d 个成功的PDF文件...\n", len(tempPDFs))
		return mergePDFs(tempPDFs, outputFile)
	}

	// 文件数量不多或都是本地文件，直接处理
	return generateSinglePDF(htmlFiles, outputFile, enableTOC)
}

func generateSinglePDF(htmlFiles []string, outputPath string, withTOC bool) error {
	// 首先尝试使用Chrome headless（如果可用）
	if len(htmlFiles) == 1 && isURL(htmlFiles[0]) {
		if chromeErr := tryChromePDF(htmlFiles[0], outputPath); chromeErr == nil {
			return nil
		}
	}

	// 构建wkhtmltopdf命令参数
	args := []string{
		"--enable-local-file-access",
		"--encoding", "utf-8",
		"--load-error-handling", "ignore",
		"--load-media-error-handling", "ignore",
		"--no-stop-slow-scripts",
		"--javascript-delay", "1000",
		"--disable-smart-shrinking",
		"--disable-javascript", // 禁用JavaScript以避免崩溃
		"--print-media-type",   // 使用print媒体类型
	}

	// 如果下载了图片，启用图片
	if !downloadImages {
		args = append(args, "--no-images")
	}

	// 页面设置
	args = append(args, "-s", pageSize)
	args = append(args, "-O", orientation)
	args = append(args, "-T", marginTop)
	args = append(args, "-B", marginBottom)
	args = append(args, "-L", marginLeft)
	args = append(args, "-R", marginRight)

	// 页眉设置
	if headerTitle != "" {
		args = append(args, "--header-center", headerTitle)
		args = append(args, "--header-font-size", "10")
		args = append(args, "--header-spacing", "5")
		args = append(args, "--header-line")
	}

	// 页脚设置
	if footerText != "" {
		args = append(args, "--footer-center", footerText)
		args = append(args, "--footer-font-size", "10")
		args = append(args, "--footer-spacing", "5")
		args = append(args, "--footer-line")
	}

	// 目录设置
	if withTOC && len(htmlFiles) > 1 {
		args = append(args, "toc")
	}

	// 添加所有HTML文件或URL
	args = append(args, htmlFiles...)

	// 输出文件
	args = append(args, outputPath)

	// 执行命令
	cmd := exec.Command("wkhtmltopdf", args...)

	if verbose {
		fmt.Printf("执行命令: wkhtmltopdf %s\n", strings.Join(args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wkhtmltopdf执行失败: %v", err)
	}

	return nil
}

func mergePDFs(pdfFiles []string, outputPath string) error {
	// 使用 gs (Ghostscript) 合并PDF，如果没有安装则使用简单的wkhtmltopdf合并
	// 首先检查是否有 gs
	if _, err := exec.LookPath("gs"); err == nil {
		// 使用 Ghostscript
		args := []string{
			"-dBATCH",
			"-dNOPAUSE",
			"-q",
			"-sDEVICE=pdfwrite",
			fmt.Sprintf("-sOutputFile=%s", outputPath),
		}
		args = append(args, pdfFiles...)

		cmd := exec.Command("gs", args...)
		if verbose {
			fmt.Printf("使用 Ghostscript 合并: gs %s\n", strings.Join(args, " "))
		}
		return cmd.Run()
	}

	// 如果没有 Ghostscript，直接使用第一个PDF作为输出
	// 这不是理想的解决方案，但至少能生成部分内容
	fmt.Println("⚠️ 未找到 Ghostscript，将只保留第一批次的内容")
	fmt.Println("建议安装 Ghostscript 以获得完整功能: brew install ghostscript")

	// 复制第一个PDF作为输出
	input, err := ioutil.ReadFile(pdfFiles[0])
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, input, 0644)
}

func tryChromePDF(url string, outputPath string) error {
	// 检查Chrome路径
	chromePaths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"google-chrome",
		"chromium",
	}

	var chromePath string
	for _, path := range chromePaths {
		if _, err := exec.LookPath(path); err == nil {
			chromePath = path
			break
		} else if _, err := os.Stat(path); err == nil {
			chromePath = path
			break
		}
	}

	if chromePath == "" {
		return fmt.Errorf("Chrome not found")
	}

	if verbose {
		fmt.Printf("尝试使用Chrome Headless生成PDF...\n")
	}

	// 使用Chrome headless生成PDF
	cmd := exec.Command(chromePath,
		"--headless",
		"--disable-gpu",
		fmt.Sprintf("--print-to-pdf=%s", outputPath),
		"--no-pdf-header-footer",
		"--virtual-time-budget=10000",
		url)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Chrome PDF generation failed: %v", err)
	}

	if verbose {
		fmt.Printf("✓ Chrome Headless成功生成PDF\n")
	}

	return nil
}
