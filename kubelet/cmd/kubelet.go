package main

import (
	"context"

	kubelet "github.com/kevindrosendahl/kubernetes-honeypot/kubelet/pkg"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	cli "github.com/virtual-kubelet/node-cli"
	logruscli "github.com/virtual-kubelet/node-cli/logrus"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	logruslogger "github.com/virtual-kubelet/virtual-kubelet/log/logrus"
)

func main() {
	ctx := cli.ContextWithCancelOnSignal(context.Background())
	logger := logrus.StandardLogger()

	log.L = logruslogger.FromLogrus(logrus.NewEntry(logger))
	logConfig := &logruscli.Config{LogLevel: "info"}

	node, err := cli.New(ctx,
		cli.WithProvider("honeypot", func(cfg provider.InitConfig) (provider.Provider, error) {
			var conf kubelet.HoneypotConfig
			if _, err := toml.Decode(cfg.ConfigPath, &conf); err != nil {
				return nil, err
			}

			return kubelet.NewHoneypotProviderFromConfig(&conf)
		}),
		// Adds flags and parsing for using logrus as the configured logger
		cli.WithPersistentFlags(logConfig.FlagSet()),
		cli.WithPersistentPreRunCallback(func() error {
			return logruscli.Configure(logConfig, logger)
		}),
	)

	if err != nil {
		panic(err)
	}

	// Notice that the context is not passed in here, this is due to limitations
	// of the underlying CLI library (cobra).
	// contexts get passed through from `cli.New()`
	//
	// Args can be specified here, or os.Args[1:] will be used.
	if err := node.Run(); err != nil {
		panic(err)
	}
}
