package service

import (
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"website-analyzer/src/entity"
)

func TestExtractDoctype(t *testing.T) {
	body := []byte(`<!DOCTYPE html>`)
	doctype := extractDoctype(body)
	assert.Equal(t, "html", doctype)

	body = []byte(`<!DOCTYPE html PUBLIC xyz>`)
	doctype = extractDoctype(body)
	assert.Equal(t, "html PUBLIC xyz", doctype)

	body = []byte(`<!DOCTYPE unknown>`)
	doctype = extractDoctype(body)
	assert.Equal(t, "unknown", doctype)

	body = []byte(`<!DOCTYPE>`)
	doctype = extractDoctype(body)
	assert.Equal(t, "", doctype)

	body = []byte(``)
	doctype = extractDoctype(body)
	assert.Equal(t, "", doctype)

	body = []byte(`<!DOCTYPE`)
	doctype = extractDoctype(body)
	assert.Equal(t, "", doctype)
}

func TestIsValidURL(t *testing.T) {
	valid := isValidURL("https://www.google.com")
	assert.True(t, valid)

	invalid := isValidURL("invalid-u%rl")
	assert.False(t, invalid)

	invalid = isValidURL("invalid-url")
	assert.False(t, invalid)
}

func TestIsInternalLink(t *testing.T) {
	baseURL := "https://www.google.com"
	link := "https://www.google.com/mail"
	internal := isInternalLink(link, baseURL)
	assert.True(t, internal)

	link = "https://www.somethingelse.com/page"
	internal = isInternalLink(link, baseURL)
	assert.False(t, internal)

	link = "invalid-u%rl/page"
	internal = isInternalLink(link, baseURL)
	assert.False(t, internal)

	baseURL = "invalid-u%rl"
	link = "invalid-u%rl/page"
	internal = isInternalLink(link, baseURL)
	assert.False(t, internal)

}

func TestIsAccessibleURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	accessible := isAccessibleURL(server.URL)
	assert.True(t, accessible)

	nonExistentURL := "https://www.thisurldoesnotexist.com"
	accessible = isAccessibleURL(nonExistentURL)
	assert.False(t, accessible)

	errorUrl := "vsnkjfvns;"
	accessible = isAccessibleURL(errorUrl)
	assert.False(t, accessible)
}

func TestSetupCollector(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}

	setupCollector(collector, req)
	assert.NotNil(t, collector.OnResponse)
	assert.NotNil(t, collector.OnHTML)
	assert.NotNil(t, collector.OnScraped)
}

func TestAnalyzeWebsiteContent(t *testing.T) {
	collector := colly.NewCollector()
	analyzerService := NewAnalyzerService(collector)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><head><title>Title</title></head><body><h1>Heading</h1><a href="https://www.exampleurl.com">Link</a></body></html>`))
	}))
	defer server.Close()

	req := entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}

	result, err := analyzerService.AnalyzeWebsiteContent(req)
	assert.NoError(t, err)
	assert.Equal(t, entity.StatusSuccess, result.JobStatus)
	assert.Contains(t, "HTML5", result.HtmlVersion)
	assert.Equal(t, "Title", result.Title)
	assert.Equal(t, 1, len(result.Headings))
	assert.Equal(t, "h1", result.Headings[0].Level)
	assert.Equal(t, 1, result.Headings[0].Count)
	assert.False(t, result.IsLogin)
	assert.Equal(t, 0, result.Links.InternalLinks)
	assert.Equal(t, 1, result.Links.ExternalLinks)
	assert.Equal(t, 0, result.Links.InaccessibleLinks)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE XHTML><html><head><title>Title</title></head><body><h1>Heading</h1><a href="https://www.exampleurl.com">Link</a></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.Contains(t, "XHTML", result.HtmlVersion)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE HTML 4.1><html><head><title>Title</title></head><body><h1>Heading</h1><a href="https://www.exampleurl.com">Link</a></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.Contains(t, "HTML 4.1", result.HtmlVersion)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><form name="login" action="/login" method="POST"><input type="text" name="username"><input type="password" name="password"><button type="submit">Login</button></form></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.True(t, result.IsLogin)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><form action="/login" method="POST"><input type="text" name="username"><input type="password" name="password"><button type="submit">Login</button></form></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.True(t, result.IsLogin)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><form action="/submit" method="POST"><input type="text" name="username"><input type="password" name="password"><button type="submit">Login</button></form></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.False(t, result.IsLogin)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE HTML 4.1><html><head><title>Title</title></head><body><h1>Heading</h1><a href="https://www.inaccessibleurl.com">Link</a><a href="/login">Link</a><a href="">Link</a></body></html>`))
	}))
	defer server.Close()

	req = entity.Analysis{
		Id:        1,
		TargetUrl: server.URL,
	}
	result, err = analyzerService.AnalyzeWebsiteContent(req)
	assert.Equal(t, 1, result.Links.InternalLinks)
	assert.Equal(t, 1, result.Links.ExternalLinks)
}

func TestAnalyzeWebsiteContentError(t *testing.T) {
	collector := colly.NewCollector()
	analyzerService := NewAnalyzerService(collector)
	req := entity.Analysis{
		Id:        1,
		TargetUrl: "127.0.0.a:80",
	}

	result, err := analyzerService.AnalyzeWebsiteContent(req)
	assert.Error(t, err)
	assert.Equal(t, entity.StatusFail, result.JobStatus)
}

func TestDoctypeHandler(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}
	setDoctypeHandler(collector, req)
	assert.NotNil(t, collector.OnResponse)
}

func TestSetTitleHandler(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}
	setTitleHandler(collector, req)
	assert.NotNil(t, collector.OnHTML)
}

func TestSetHeadingsHandler(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}
	setHeadingsHandler(collector, req)
	assert.NotNil(t, collector.OnHTML)
}

func TestSetLoginFormHandler(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}
	setLoginFormHandler(collector, req)
	assert.NotNil(t, collector.OnHTML)
}

func TestSetLinksHandler(t *testing.T) {
	collector := colly.NewCollector()
	req := &entity.Analysis{}
	setLinksHandler(collector, req)
	assert.NotNil(t, collector.OnHTML)
}
