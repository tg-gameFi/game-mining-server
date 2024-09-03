package api

import (
	"compress/flate"
	"compress/gzip"
	"game-mining-server/entities"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var htmlIgnoreHeaderKeys = []string{"x-frame-options", "content-encoding", "cross-origin-opener-policy", "referrer-policy"}
var reqIgnoreHeaderKeys = []string{"x-frame-options", "cross-origin-opener-policy", "referrer-policy"}
var relativePathReplacements = []string{"/_next", "/manifest.json", "/favicon.ico"}

func isIgnoreHeader(key string, ignoreKeys []string) bool {
	for _, _key := range ignoreKeys {
		if strings.ToLower(key) == _key {
			return true
		}
	}
	return false
}

func switchContentEncoding(res *http.Response) (bodyReader io.Reader, err error) {
	switch res.Header.Get("Content-Encoding") {
	case "br":
		bodyReader = brotli.NewReader(res.Body)
	case "gzip":
		bodyReader, err = gzip.NewReader(res.Body)
	case "deflate":
		bodyReader = flate.NewReader(res.Body)
	default:
		bodyReader = res.Body
	}
	return
}

// replace relative path to absolute path in html file
func replaceRelativePathInRes(realUrl string, res *http.Response) (string, error) {
	resBody, e0 := switchContentEncoding(res)
	if e0 != nil {
		return "", e0
	}

	body, e1 := io.ReadAll(resBody)
	if e1 != nil {
		return "", e1
	}
	realPath, e2 := url.Parse(realUrl)
	if e2 != nil {
		return "", e2
	}

	htmlBody := string(body)
	realPathPrefix := realPath.Scheme + "://" + realPath.Host // https://pancakeswap.com
	for _, relativePath := range relativePathReplacements {
		htmlBody = strings.ReplaceAll(htmlBody, "\""+relativePath, "\""+realPathPrefix+relativePath)
	}
	return htmlBody, nil
}

// ProxyGetHtml Request resource(html, css, js, images, assets) by proxy
func ProxyGetHtml(c *gin.Context) {
	var params entities.ProxyGetParam
	if e0 := c.ShouldBindQuery(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	proxyReq, e1 := http.NewRequest(c.Request.Method, params.Url, c.Request.Body)
	if e1 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyCreateRequestFailed, e1.Error()))
		return
	}
	proxyReq.Header = c.Request.Header

	client := &http.Client{}
	proxyRes, e2 := client.Do(proxyReq)
	if e2 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyRequestFailed, e2.Error()))
		return
	}
	defer func() {
		_ = proxyRes.Body.Close()
	}()

	resBody, e3 := replaceRelativePathInRes(params.Url, proxyRes)
	if e3 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyParseResBodyFailed, e3.Error()))
		return
	}

	for key, value := range proxyRes.Header {
		if len(value) == 1 && !isIgnoreHeader(key, htmlIgnoreHeaderKeys) {
			c.Writer.Header().Add(key, value[0])
		}
	}

	c.Status(proxyRes.StatusCode)
	_, _ = c.Writer.Write([]byte(resBody))
}

func ProxyGetNextRes(c *gin.Context) {
	realUrl := "https://pancakeswap.finance" + c.Request.URL.Path
	log.Printf("Proxy res url: %s", realUrl)
	req, e1 := http.NewRequest(c.Request.Method, realUrl, c.Request.Body)
	if e1 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyCreateRequestFailed, e1.Error()))
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, e2 := client.Do(req)
	if e2 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyRequestFailed, e2.Error()))
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, e3 := io.ReadAll(resp.Body)
	if e3 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyParseResBodyFailed, e3.Error()))
		return
	}

	for key, value := range resp.Header {
		if len(value) == 1 && !isIgnoreHeader(key, reqIgnoreHeaderKeys) {
			c.Writer.Header().Add(key, value[0])
		}
	}
	c.Status(resp.StatusCode)
	_, _ = c.Writer.Write(body)
}

func ProxyRequest(c *gin.Context) {
	var params entities.ProxyGetParam
	if e0 := c.ShouldBindQuery(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	req, e1 := http.NewRequest(c.Request.Method, params.Url, c.Request.Body)
	if e1 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyCreateRequestFailed, e1.Error()))
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, e2 := client.Do(req)
	if e2 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyRequestFailed, e2.Error()))
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, e3 := io.ReadAll(resp.Body)
	if e3 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrProxyParseResBodyFailed, e3.Error()))
		return
	}

	for key, value := range resp.Header {
		if len(value) == 1 && !isIgnoreHeader(key, reqIgnoreHeaderKeys) {
			c.Writer.Header().Add(key, value[0])
		}
	}
	c.Status(resp.StatusCode)
	_, _ = c.Writer.Write(body)
}
