package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const (
	LogLevel = "log.level"

	DragonflyHost = "dragonfly.host"
	DragonflyPort = "dragonfly.port"
	DragonflyAuth = "dragonfly.auth"

	CockroachURL = "cockroach.url"

	MetricsEnabled = "metrics.enabled"
	MetricsPort    = "metrics.port"

	TracingEnabled = "tracing.enabled"
	TracingRatio   = "tracing.ratio"
	ServiceName    = "service.name"
	ServiceVersion = "service.version"

	JWESecret = "jwe.secret"

	AESKey = "aes.key"
)

func Get(prefix string) (*koanf.Koanf, error) {
	k := koanf.New(".")

	envVarPrefix := fmt.Sprintf("%s_", strings.ToUpper(prefix))

	err := k.Load(env.Provider(envVarPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envVarPrefix)), "_", ".", -1)
	}), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return k, nil
}
