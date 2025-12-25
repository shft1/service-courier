package cli

import (
	"context"
	"os"
	"service-courier/internal/config/appcfg"

	"github.com/urfave/cli/v3"
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
				os.Setenv("COURIER_LOCALPORT", port)
				env.AppPort = port
			}
			return nil
		},
	}
}
