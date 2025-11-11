package cli

import (
	"context"
	"os"
	"service-courier/internal/config"

	"github.com/urfave/cli/v3"
)

func CliHandler(env *config.Env) *cli.Command {
	return &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Usage: "Port for web-server",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if port := cmd.String("port"); port != "" {
				os.Setenv("PORT", port)
				env.Port = port
			}
			return nil
		},
	}
}
