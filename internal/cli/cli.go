package cli

import (
	"context"

	"github.com/urfave/cli/v3"

	"service-courier/internal/config/appcfg"
)

// CliHandler - парсер командной строки
func CliHandler(env *appcfg.AppEnv) *cli.Command {
	return &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Usage: "Port for web-server",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if port := cmd.String("port"); port != "" {
				env.AppPort = port
			}
			return nil
		},
	}
}
