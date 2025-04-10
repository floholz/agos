package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/floholz/agos/cmd"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:           "agos",
		Version:        "0.1.1",
		HelpName:       "AGOS - adb go screenshare",
		Description:    "An ADB + SCRCPY wrapper, with automated port discovery",
		Usage:          "Select a device for screensharing or connect a new one",
		DefaultCommand: "screenshare",
		Commands: []*cli.Command{
			{
				Name:    "screenshare",
				Aliases: []string{"s", "screen", "share"},
				Usage:   "List devices for screensharing",
				Action: func(c *cli.Context) error {
					app := cmd.NewAgosApp()
					err := app.Adb.StartServer() //.StartAdb()
					if err != nil {
						return err
					}

					devices, err := app.Adb.ListDevices() //.AdbListDevices()
					if err != nil {
						return err
					}

					if len(devices) == 0 {
						fmt.Println("No devices found.")
						return nil
					}

					prompt := promptui.Select{
						Label:    "Select a Device",
						Items:    devices,
						HideHelp: true,
					}

					index, result, err := prompt.Run()
					if err != nil {
						return fmt.Errorf("prompt failed: %v", err)
					}

					fmt.Printf("Selected: %s\n", result)
					fmt.Printf("Starting action on device %d: %s\n", index+1, result)

					err = app.ScrCpy.Run(devices[index].Address)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "connect",
				Aliases: []string{"c"},
				Usage:   "Connect new device",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "min-port",
						Value:   32000,
						Usage:   "Minimum port range",
						EnvVars: []string{"AGOS_MIN_PORT"},
					},
					&cli.IntFlag{
						Name:    "max-port",
						Value:   48000,
						Usage:   "Maximum port range",
						EnvVars: []string{"AGOS_MAX_PORT"},
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("IP address is required as an argument")
					}

					ip := c.Args().Get(0)
					if !strings.Contains(ip, ":") {
						minPort := c.Int("min-port")
						maxPort := c.Int("max-port")

						port, err := cmd.DiscoverAdbPort(ip, minPort, maxPort)
						if err != nil {
							return err
						}
						ip = fmt.Sprintf("%s:%d", ip, port)
					}

					app := cmd.NewAgosApp()
					err := app.Adb.StartServer() //.StartAdb()
					if err != nil {
						return err
					}
					err = app.Adb.Connect(ip) //.AdbConnect(ip)
					if err != nil {
						return err
					}
					err = app.ScrCpy.Run(ip)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "pair",
				Aliases: []string{"p"},
				Usage:   "Pair new device",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "min-port",
						Value:   32000,
						Usage:   "Minimum port range",
						EnvVars: []string{"AGOS_MIN_PORT"},
					},
					&cli.IntFlag{
						Name:    "max-port",
						Value:   48000,
						Usage:   "Maximum port range",
						EnvVars: []string{"AGOS_MAX_PORT"},
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return fmt.Errorf("IP address and pairing code are required as arguments")
					}

					ip := c.Args().Get(0)
					if !strings.Contains(ip, ":") {
						minPort := c.Int("min-port")
						maxPort := c.Int("max-port")

						port, err := cmd.DiscoverAdbPort(ip, minPort, maxPort)
						if err != nil {
							return err
						}
						ip = fmt.Sprintf("%s:%d", ip, port)
					}

					code := c.Args().Get(1)
					if code == "" {
						return fmt.Errorf("pairing code cant be empty")
					}

					app := cmd.NewAgosApp()
					err := app.Adb.StartServer() //.StartAdb()
					if err != nil {
						return err
					}
					err = app.Adb.Pair(ip, code) //.AdbPair(ip, code)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "pair-qr",
				Aliases: []string{"qr"},
				Usage:   "Pair new device with QR code",
				Action: func(c *cli.Context) error {
					app := cmd.NewAgosApp()
					err := app.Adb.StartServer() //.StartAdb()
					if err != nil {
						return err
					}
					err = app.Adb.PairQR() //.AdbPairQr()
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
