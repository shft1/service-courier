package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func CliHandler(ctx context.Context, env *Env) {
	cmd := &cli.Command{
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
	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
	}
}
