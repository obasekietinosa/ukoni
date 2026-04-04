package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"ukoni/internal/config"
)

type AgentHandler struct {
	Proxy *httputil.ReverseProxy
}

func NewAgentHandler(cfg *config.Config) (*AgentHandler, error) {
	targetURL, err := url.Parse(cfg.AgentServiceURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	return &AgentHandler{Proxy: proxy}, nil
}

func (h *AgentHandler) HandleProxy(w http.ResponseWriter, r *http.Request) {
	h.Proxy.ServeHTTP(w, r)
}
