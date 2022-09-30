package main

import (
	"crypto/tls"
	"net"
	"os"
	"strconv"

	"github.com/romnn/flags4urfavecli/flags"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"layeh.com/gumble/gumble"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// Version is incremented using bump2version
const Version = "0.0.2"

func main() {
	app := &cli.App{
		Version:   Version,
		Usage:     "Create channel hierarchies on mumble server",
		Name:      "mumble_createchannels",
		ArgsUsage: "server_address",
		Flags: []cli.Flag{
			&flags.LogLevelFlag,
			&cli.PathFlag{
				Name:  "cert",
				Usage: "Path to certificate `FILE` (PEM)",
				//Required: true,
			},
			&cli.PathFlag{
				Name:  "key",
				Usage: "Path to key `FILE` (PEM)",
				//Required: true,
			},
			&cli.StringFlag{
				Name:  "pass",
				Usage: "Server password",
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "Server `PORT`",
				Value: gumble.DefaultPort,
			},
			&cli.StringFlag{
				Name:  "user",
				Usage: "Username to use",
				Value: "gumble_channel_bot",
			},
			&cli.BoolFlag{
				Name:  "insecure",
				Usage: "Do not verify server certificate",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			if level, err := log.ParseLevel(ctx.String("log")); err == nil {
				log.SetLevel(level)
			}
			if ctx.Args().Len() < 1 {
				log.Error("server address missing")
				cli.ShowAppHelpAndExit(ctx, 1)
			}

			var tlsConfig tls.Config
			if ctx.Bool("insecure") {
				tlsConfig.InsecureSkipVerify = true
			}
			cert := ctx.Path("cert")
			if cert != "" {
				key := ctx.Path("key")
				if key == "" {
					// In case cert and key are in same file
					key = cert
				}
				certificate, err := tls.LoadX509KeyPair(cert, key)
				if err != nil {
					log.Error(err)
					return cli.Exit("Cannot load certificate", 1)
				}
				tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
			}

			config := gumble.NewConfig()
			config.Username = ctx.String("user")
			serverpass := ctx.String("pass")
			if serverpass != "" {
				config.Password = serverpass //pragma: allowlist secret
			}

			address := net.JoinHostPort(ctx.Args().First(), strconv.Itoa(ctx.Int("port")))

			log.Info("Dialing ", address)
			client, err := gumble.DialWithDialer(new(net.Dialer), address, config, &tlsConfig)
			if err != nil {
				log.Error(err)
				return cli.Exit("Could not connect", 1)
			}

			log.Info("Disconnecting")
			err = client.Disconnect()
			if err != nil {
				log.Error(err)
				return cli.Exit("Could not disconnect", 1)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
