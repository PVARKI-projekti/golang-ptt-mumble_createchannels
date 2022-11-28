package main

import (
	"crypto/tls"
	"net"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/romnn/flags4urfavecli/flags"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"layeh.com/gumble/gumble"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// Version is incremented using bump2version
const Version = "0.1.0"

func main() {
	app := &cli.App{
		Version:   Version,
		Usage:     "Create channel hierarchies on mumble server",
		Name:      "mumble_createchannels",
		ArgsUsage: "server_address channel_map.yaml",
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
			&cli.IntFlag{
				Name:  "wait",
				Usage: "wait for server ping for this many seconds",
				Value: 5,
			},
		},
		Action: func(ctx *cli.Context) error {
			if level, err := log.ParseLevel(ctx.String("log")); err == nil {
				log.SetLevel(level)
			}
			if ctx.Args().Len() < 1 {
				log.Fatal("Server address missing")
				cli.ShowAppHelpAndExit(ctx, 1)
			}
			if ctx.Args().Len() < 2 {
				log.Fatal("Channel map is missing")
				cli.ShowAppHelpAndExit(ctx, 1)
			}

			yfile, err := os.ReadFile(ctx.Args().Get(1))
			if err != nil {
				log.Fatal(err)
				return cli.Exit("Cannot open yaml file", 1)
			}
			yamldata := make(map[string]interface{})
			err = yaml.Unmarshal(yfile, &yamldata)
			if err != nil {
				log.Fatal(err)
				return cli.Exit("Cannot parse yaml file", 1)
			}

			//log.Info("yamldata: ", pp.Sprint(yamldata))

			channels, ok := yamldata["channels"]
			if !ok {
				log.Error("No 'channels' key at root level in YAML")
				return cli.Exit("Nothing to do", 1)
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
					log.Fatal(err)
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
			interval := time.Second * 1
			timeout := time.Second * time.Duration(ctx.Int("wait"))

			address := net.JoinHostPort(ctx.Args().First(), strconv.Itoa(ctx.Int("port")))
			log.Info("Pinging ", address)
			_, err = gumble.Ping(address, interval, timeout)
			if err != nil {
				log.Fatal(err)
				return cli.Exit("Could not ping", 1)
			}

			log.Info("Dialing ", address)
			client, err := gumble.DialWithDialer(new(net.Dialer), address, config, &tlsConfig)
			if err != nil {
				log.Fatal(err)
				return cli.Exit("Could not connect", 1)
			}

			recurseChannelMap(client, client.Channels[0], channels)

			log.Info("Disconnecting")
			err = client.Disconnect()
			if err != nil {
				log.Fatal(err)
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

// Add channel and check if it got added
func addAndCheck(client *gumble.Client, parent *gumble.Channel, name string) *gumble.Channel {
	log.Info("Trying to create channel ", name, " under ", parent.Name)
	var retval *gumble.Channel = nil

	// FIXME: listen for events with some timeout instead
	for i := range [4]int{} {
		client.Do(func() {
			parent.Add(name, false)
		})
		time.Sleep(250 * time.Millisecond)
		client.Do(func() {
			child := parent.Find(name)
			if child != nil {
				retval = child
			}
		})
		if retval != nil {
			break
		}
		log.Warning("Did not find channel, trying again in 1s, this was attempt #", i)
		time.Sleep(1000 * time.Millisecond)
	}
	if retval == nil {
		log.Error("Looks like create failed")
	}
	return retval
}

func recurseChannelMap(client *gumble.Client, parent *gumble.Channel, children interface{}) bool {
	for idx, item := range children.([]interface{}) {
		childconf := item.(map[string]interface{})
		name, found := childconf["name"]
		if !found {
			log.Error("'name' not found in item #", idx)
			continue
		}
		childname := ""
		switch name.(type) {
		case string:
			childname = name.(string)
		case int:
			childname = strconv.Itoa(name.(int))
		}
		childch := addAndCheck(client, parent, childname)
		if childch == nil {
			log.Error("Could not create child ", name)
			continue
		}
		description, found := childconf["description"]
		if found {
			client.Do(func() {
				childch.SetDescription(description.(string))
			})
		}
		grandchildren, ok := childconf["channels"]
		if ok {
			recurseChannelMap(client, childch, grandchildren)
		}
	}
	return true
}
