package cli

import (
	"context"

	steveauth "github.com/rancher/steve/pkg/auth"
	authcli "github.com/rancher/steve/pkg/auth/cli"
	"github.com/rancher/steve/pkg/server"
	"github.com/rancher/steve/pkg/ui"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/urfave/cli"
)

type Config struct {
	KubeConfig      string
	Context         string
	HTTPSListenPort int
	HTTPListenPort  int
	UIPath          string
	Offline         string

	WebhookConfig authcli.WebhookConfig
}

func (c *Config) MustServer(ctx context.Context) *server.Server {
	cc, err := c.ToServer(ctx)
	if err != nil {
		panic(err)
	}
	return cc
}

func (c *Config) ToServer(ctx context.Context) (*server.Server, error) {
	var (
		auth steveauth.Middleware
	)

	restConfig, err := kubeconfig.GetNonInteractiveClientConfigWithContext(c.KubeConfig, c.Context).ClientConfig()
	if err != nil {
		return nil, err
	}
	restConfig.RateLimiter = ratelimit.None

	if c.WebhookConfig.WebhookAuthentication {
		auth, err = c.WebhookConfig.WebhookMiddleware()
		if err != nil {
			return nil, err
		}
	}

	return server.New(ctx, restConfig, &server.Options{
		AuthMiddleware: auth,
		Next:           ui.New(c.UIPath, c.Offline),
	})
}

func Flags(config *Config) []cli.Flag {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			EnvVar:      "KUBECONFIG",
			Destination: &config.KubeConfig,
		},
		cli.StringFlag{
			Name:        "context",
			EnvVar:      "CONTEXT",
			Destination: &config.Context,
		},
		cli.StringFlag{
			Name:        "ui-path",
			Destination: &config.UIPath,
		},
		cli.StringFlag{
			Name:        "offline",
			Value:       "dynamic",
			Usage:       "Determine whether or not to run the UI offline, accepts values true/false/dynamic",
			Destination: &config.Offline,
		},
		cli.IntFlag{
			Name:        "https-listen-port",
			Value:       9443,
			Destination: &config.HTTPSListenPort,
		},
		cli.IntFlag{
			Name:        "http-listen-port",
			Value:       9080,
			Destination: &config.HTTPListenPort,
		},
	}

	return append(flags, authcli.Flags(&config.WebhookConfig)...)
}
