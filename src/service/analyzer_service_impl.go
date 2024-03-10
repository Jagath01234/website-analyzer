package service

import (
	colly "github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"net/url"
	"strings"
	"website-analyzer/src/entity"
)

const logPrefixAnalyzerService = "web-analyzer/src/service/analyzer_service_impl"

type AnalyzerService struct {
	collector *colly.Collector
}

func NewAnalyzerService(collector *colly.Collector) *AnalyzerService {
	return &AnalyzerService{
		collector: collector,
	}
}

func (a *AnalyzerService) AnalyzeWebsiteContent(req entity.Analysis) (entity.Analysis, error) {

	setupCollector(a.collector, &req)
	err := a.collector.Visit(req.TargetUrl)
	if err != nil {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "URL analyze error:", err.Error(), req.TargetUrl)
		req.JobStatus = entity.StatusFail
		return req, err
	}
	req.JobStatus = entity.StatusSuccess

	return req, nil
}

func setupCollector(collector *colly.Collector, req *entity.Analysis) {
	setDoctypeHandler(collector, req)
	setTitleHandler(collector, req)
	setHeadingsHandler(collector, req)
	setLoginFormHandler(collector, req)
	setLinksHandler(collector, req)
}

func setDoctypeHandler(collector *colly.Collector, req *entity.Analysis) {
	collector.OnResponse(func(r *colly.Response) {
		docType := r.Ctx.Get("doctype")
		if docType == "" {
			docType = extractDoctype(r.Body)
		}
		if docType == "" || docType == "html" || strings.HasPrefix(docType, "html ") {
			req.HtmlVersion = "HTML5"
		}
		if strings.Contains(docType, "HTML 4.1") {
			req.HtmlVersion = "HTML 4.1"
		}
		if strings.Contains(docType, "XHTML") {
			req.HtmlVersion = "XHTML"
		}
	})
}

func setTitleHandler(collector *colly.Collector, req *entity.Analysis) {
	collector.OnHTML("title", func(e *colly.HTMLElement) {
		req.Title = e.Text
	})
}

func setHeadingsHandler(collector *colly.Collector, req *entity.Analysis) {
	headingCounts := make(map[string]int)
	collector.OnHTML("h1, h2, h3, h4, h5, h6", func(e *colly.HTMLElement) {
		level := e.Name
		headingCounts[level]++
	})

	collector.OnScraped(func(r *colly.Response) {
		var headings []entity.HeadingInfo
		for level, count := range headingCounts {
			heading := entity.HeadingInfo{
				Level: level,
				Count: count,
			}
			headings = append(headings, heading)
		}
		req.Headings = headings
	})
}

func setLoginFormHandler(collector *colly.Collector, req *entity.Analysis) {
	collector.OnHTML("form", func(e *colly.HTMLElement) {
		formAction := e.Attr("action")
		formMethod := e.Attr("method")
		formName := e.Attr("name")
		if strings.Contains(formName, "login") {
			req.IsLogin = true
		} else {
			if (formAction != "" || formMethod != "") && (formAction == "/login" || formAction == "/signin") {
				req.IsLogin = true
			} else {
				req.IsLogin = false
			}
		}
	})
}

func setLinksHandler(collector *colly.Collector, req *entity.Analysis) {
	var internalLinks, externalLinks, inaccessibleLinks int

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if link == "" || !isValidURL(absoluteURL) {
			return
		}

		if isInternalLink(absoluteURL, e.Request.URL.String()) {
			internalLinks++
		} else {
			externalLinks++
		}

		if !isAccessibleURL(absoluteURL) {
			inaccessibleLinks++
		}
	})

	collector.OnScraped(func(r *colly.Response) {
		req.Links.InternalLinks = internalLinks
		req.Links.ExternalLinks = externalLinks
		req.Links.InaccessibleLinks = inaccessibleLinks
	})
}

func extractDoctype(body []byte) string {
	docType := string(body)
	startIndex := strings.Index(docType, "<!DOCTYPE")
	if startIndex == -1 {
		return ""
	}
	endIndex := strings.Index(docType[startIndex:], ">")
	if endIndex == -1 {
		return ""
	}

	doctype := docType[startIndex : startIndex+endIndex+1]
	doctype = strings.ReplaceAll(doctype, "<!DOCTYPE", "")
	doctype = strings.ReplaceAll(doctype, ">", "")
	return strings.TrimSpace(doctype)
}

func isValidURL(targetUrl string) bool {
	parsed, err := url.Parse(targetUrl)
	if err != nil {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "URL validation error:", err.Error(), targetUrl)
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		log.Printf("%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "URL validation error:", targetUrl)
		return false
	}
	return true
}

func isInternalLink(link, baseURL string) bool {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "baseUrl parse error:", err.Error(), baseURL)
		return false
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "link parse error:", err.Error(), link)
		return false
	}
	return parsedLink.Host == parsedBase.Host
}

func isAccessibleURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		log.Printf("%v,%v,%v,%v,%v", "ERROR", logPrefixAnalyzerService, "head assigning error:", err.Error(), url)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
