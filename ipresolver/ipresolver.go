package ipresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	apiKeyHeader       = "x-api-key"
	defaultHTTPTimeout = 10 * time.Second
)

type response struct {
	IP        string `json:"ip"`
	IPAddress string `json:"ip_address"`
}

func ResolveIP(ctx context.Context, resolverURL string, resolverAPIKey string) (string, error) {
	log.Printf("ipresolver: building request")

	client := &http.Client{Timeout: defaultHTTPTimeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, resolverURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	req.Header.Set(apiKeyHeader, resolverAPIKey)
	log.Printf("ipresolver: executing request to %q", resolverURL)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	log.Printf("ipresolver: response status=%s body_bytes=%d", resp.Status, len(body))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		bodyPreview := strings.TrimSpace(string(body))
		if len(bodyPreview) > 256 {
			bodyPreview = bodyPreview[:256]
		}
		return "", fmt.Errorf("resolver returned status=%d body=%q", resp.StatusCode, bodyPreview)
	}

	resolved := &response{}
	if err := json.Unmarshal(body, resolved); err != nil {
		return "", fmt.Errorf("unmarshal resolver response: %w", err)
	}

	ip := strings.TrimSpace(resolved.IPAddress)
	if ip == "" {
		ip = strings.TrimSpace(resolved.IP)
	}
	if ip == "" {
		return "", fmt.Errorf("resolver response missing ip/ip_address field")
	}

	if parsed := net.ParseIP(ip); parsed == nil {
		return "", fmt.Errorf("resolver returned invalid IP address %q", ip)
	}

	log.Printf("ipresolver: successfully parsed IP=%s", ip)
	return ip, nil
}
