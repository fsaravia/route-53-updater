package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsaravia/route-53-updater/ipresolver"
	"github.com/fsaravia/route-53-updater/route53"
)

const (
	envResolverURL  = "RESOLVER_URL"
	envAPIKey       = "API_KEY"
	envHostedZoneID = "HOSTED_ZONE_ID"
	envRecordSet    = "RECORD_SET"
)

type config struct {
	ResolverURL    string
	ResolverAPIKey string
	HostedZoneID   string
	RecordSet      string
}

func mustEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}

	return value, nil
}

func loadConfig() (*config, error) {
	resolverURL, err := mustEnv(envResolverURL)
	if err != nil {
		return nil, err
	}

	apiKey, err := mustEnv(envAPIKey)
	if err != nil {
		return nil, err
	}

	hostedZoneID, err := mustEnv(envHostedZoneID)
	if err != nil {
		return nil, err
	}

	recordSet, err := mustEnv(envRecordSet)
	if err != nil {
		return nil, err
	}

	return &config{
		ResolverURL:    resolverURL,
		ResolverAPIKey: apiKey,
		HostedZoneID:   hostedZoneID,
		RecordSet:      recordSet,
	}, nil
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log.Printf("configuration loaded: resolver_url=%q hosted_zone_id=%q record_set=%q", cfg.ResolverURL, cfg.HostedZoneID, cfg.RecordSet)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("resolving IP address using resolver URL %q", cfg.ResolverURL)
	ip, err := ipresolver.ResolveIP(ctx, cfg.ResolverURL, cfg.ResolverAPIKey)
	if err != nil {
		return fmt.Errorf("resolve IP: %w", err)
	}
	log.Printf("resolved public IP address: %s", ip)

	log.Printf("creating Route53 client")
	r53, err := route53.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("create Route53 client: %w", err)
	}

	log.Printf("upserting A record set=%q in hosted zone=%q with IP=%q", cfg.RecordSet, cfg.HostedZoneID, ip)
	output, err := route53.UpsertZone(ctx, r53, cfg.HostedZoneID, cfg.RecordSet, ip)
	if err != nil {
		return fmt.Errorf("upsert Route53 record: %w", err)
	}

	log.Printf("Route53 update submitted successfully: %s", output)
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)

	if err := run(); err != nil {
		log.Printf("fatal error: %v", err)
		if !errors.Is(err, context.DeadlineExceeded) {
			os.Exit(1)
		}
		os.Exit(2)
	}

	log.Printf("run completed")
}
