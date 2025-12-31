package htp

import (
	"io"
	"net/http"
	"net/url"
)

// HTTP 的编程接口

func Get(url string) (*http.Response, error)
func Post(url, contentType string, body io.Reader) (*http.Response, error)
func PostForm(url string, data url.Values) (*http.Response, error)

func NewRequest(method, url string, body io.Reader) (*http.Request, error)

type Client struct{}

func (c *Client) Do(req *http.Request) (*http.Response, error)

func NewServerMux() *http.ServeMux

type ServeMux struct{}

func (mux *ServeMux) Handle(pattern string, handler http.Handler)
func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

func ListenAndServe(addr string, handler http.Handler) error
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error
