package cli

import (
	"fmt"
	"github.com/floholz/agos/internal/core"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func LaunchCLI(agos *core.AgosApp) {
	app := &cli.App{
		Name:           "agos",
		Version:        "0.2.0",
		HelpName:       "AGOS - adb go screenshare",
		Description:    "An ADB + SCRCPY wrapper, with automated port discovery",
		Usage:          "Select a device for screensharing or connect a new one",
		DefaultCommand: "screenshare",
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
		Before: func(c *cli.Context) error {
			// override config from arguments
			minPort := c.Int("min-port")
			if minPort > 0 {
				agos.Config.PortRange.MinPort = minPort
			}
			maxPort := c.Int("max-port")
			if maxPort > 0 {
				agos.Config.PortRange.MaxPort = maxPort
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "screenshare",
				Aliases: []string{"s", "screen", "share"},
				Usage:   "List devices for screensharing",
				Action: func(c *cli.Context) error {
					err := agos.Adb.StartServer()
					if err != nil {
						return err
					}

					devices, err := agos.Adb.ListDevices()
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

					err = agos.ScrCpy.Run(devices[index].Address)
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
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("IP address is required as an argument")
					}

					ip := c.Args().Get(0)
					if !strings.Contains(ip, ":") {
						port, err := core.DiscoverAdbPort(ip, agos.Config.PortRange)
						if err != nil {
							return err
						}
						ip = fmt.Sprintf("%s:%d", ip, port)
					}

					err := agos.Adb.StartServer()
					if err != nil {
						return err
					}
					err = agos.Adb.Connect(ip)
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
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return fmt.Errorf("IP address and pairing code are required as arguments")
					}

					address := c.Args().Get(0)
					adrSplits := strings.Split(address, ":")
					if len(adrSplits) < 2 {
						return fmt.Errorf("IP address must include the port [ip:port]")
					}
					if len(adrSplits) > 3 {
						return fmt.Errorf("IP address must be a valid ip + port [ip:port]")
					}
					if len(adrSplits[0]) < 8 {
						return fmt.Errorf("IP address must be a valid ip [X.X.X.X]")
					}
					if len(adrSplits[1]) < 4 {
						return fmt.Errorf("the address must include a valid port")
					}

					code := c.Args().Get(1)
					if code == "" {
						return fmt.Errorf("pairing code cant be empty")
					}

					err := agos.Adb.StartServer()
					if err != nil {
						return err
					}
					err = agos.Adb.Pair(address, code)
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "no-connect",
						Aliases: []string{"n"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					err := agos.Adb.StartServer()
					if err != nil {
						return err
					}

					connect := !c.Bool("no-connect")

					err = agos.Adb.PairQR(connect)
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
