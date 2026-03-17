package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FluxerAPIURL     *url.URL
	FluxerGatewayURL *url.URL
	ListenAddr       string
	LogLevel         slog.Level
}

func mustParseURL(rawURL string) *url.URL {
	result, err := url.Parse(rawURL)
	if err != nil {
		panic(fmt.Errorf("failed to parse URL constant: %w", err))
	}

	return result
}

var defaults = Config{
	FluxerAPIURL:     mustParseURL("https://api.fluxer.app/"),
	FluxerGatewayURL: mustParseURL("wss://gateway.fluxer.app"),
	ListenAddr:       ":8080",
	LogLevel:         slog.LevelInfo,
}

func Load() (Config, error) {
	dotenv, err := godotenv.Read()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("failed to load dotenv: %w", err)
	}

	lookup := func(key string) string {
		if fromDotenv, _ := dotenv[key]; fromDotenv != "" {
			return fromDotenv
		} else {
			return os.Getenv(key)
		}
	}

	result := defaults

	if v := lookup("FLUXER_API_URL"); v != "" {
		parsed, err := url.Parse(v)
		if err != nil {
			return Config{}, fmt.Errorf("failed to parse FLUXER_API_URL: %w", err)
		}

		result.FluxerAPIURL = parsed
	}

	if v := lookup("FLUXER_GATEWAY_URL"); v != "" {
		parsed, err := url.Parse(v)
		if err != nil {
			return Config{}, fmt.Errorf("failed to parse FLUXER_GATEWAY_URL: %w", err)
		}

		result.FluxerGatewayURL = parsed
	}

	if v := lookup("FLINE_LISTEN_ADDR"); v != "" {
		result.ListenAddr = v
	}

	if v := lookup("FLINE_LOG_LEVEL"); v != "" {
		err := result.LogLevel.UnmarshalText([]byte(v))
		if err != nil {
			return Config{}, fmt.Errorf("failed to parse FLINE_LOG_LEVEL: %w", err)
		}
	}

	return result, nil
}
