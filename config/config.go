package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl         string `env:"DATABASE_URL,required"`
	CorsOrigin          string `env:"CORS_ORIGIN,required"`
	Port                string `env:"PORT,default=4000"`
	SessionSecret       string `env:"SECRET,required"`
	Domain              string `env:"DOMAIN"`
	ImageKitPrivateKey  string `env:"IMAGEKITPRIVATEKEY"`
	ImageKitPublicKey   string `env:"IMAGEKITPUBLICKEY"`
	ImageKitURLEndpoint string `env:"IMAGEKITURLENDPPOINT"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}
	return
}
