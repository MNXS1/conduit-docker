package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Psiphon-Labs/psiphon-tunnel-core/psiphon"
	"github.com/Psiphon-Labs/psiphon-tunnel-core/psiphon/common/buildinfo"
)

func main() {
	var configFile string
	var dataRoot string
	var maxClients int
	var limitUpstream int
	var limitDownstream int
	var version bool

	flag.StringVar(&configFile, "config", "", "Path to psiphon config JSON file")
	flag.StringVar(&dataRoot, "dataRootDirectory", "/data", "Directory for persistent data")
	flag.IntVar(&maxClients, "maxClients", 10, "Maximum number of concurrent clients")
	flag.IntVar(&limitUpstream, "limitUpstream", 0, "Upstream bandwidth limit (bytes/sec, 0=unlimited)")
	flag.IntVar(&limitDownstream, "limitDownstream", 0, "Downstream bandwidth limit (bytes/sec, 0=unlimited)")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		b := buildinfo.GetBuildInfo()
		fmt.Printf("Inproxy Node\n  Build Date: %s\n  Built With: %s\n  Repository: %s\n  Revision: %s\n",
			b.BuildDate, b.GoVersion, b.BuildRepo, b.BuildRev)
		os.Exit(0)
	}

	psiphon.SetNoticeWriter(os.Stderr)
	defer psiphon.ResetNoticeWriter()

	psiphon.SetEmitDiagnosticNotices(true, false)

	var config *psiphon.Config
	var err error

	if configFile != "" {
		configFileContents, err := ioutil.ReadFile(configFile)
		if err != nil {
			psiphon.NoticeError("error loading configuration file: %s", err)
			os.Exit(1)
		}
		config, err = psiphon.LoadConfig(configFileContents)
		if err != nil {
			psiphon.NoticeError("error processing configuration file: %s", err)
			os.Exit(1)
		}
	} else {
		configJSON := fmt.Sprintf(`{
			"DataRootDirectory": "%s",
			"DisableTunnels": true,
			"PropagationChannelId": "0000000000000000",
			"SponsorId": "0000000000000000",
			"InproxyMaxClients": %d,
			"InproxyLimitUpstreamBytesPerSecond": %d,
			"InproxyLimitDownstreamBytesPerSecond": %d
		}`, dataRoot, maxClients, limitUpstream, limitDownstream)

		config, err = psiphon.LoadConfig([]byte(configJSON))
		if err != nil {
			psiphon.NoticeError("error creating configuration: %s", err)
			os.Exit(1)
		}
	}

	if dataRoot != "" {
		config.DataRootDirectory = dataRoot
	}

	if maxClients > 0 {
		config.InproxyMaxClients = maxClients
	}
	if limitUpstream > 0 {
		config.InproxyLimitUpstreamBytesPerSecond = limitUpstream
	}
	if limitDownstream > 0 {
		config.InproxyLimitDownstreamBytesPerSecond = limitDownstream
	}

	config.DisableTunnels = true

	err = config.Commit(true)
	if err != nil {
		psiphon.NoticeError("error committing configuration: %s", err)
		os.Exit(1)
	}

	psiphon.NoticeBuildInfo()

	err = psiphon.OpenDataStore(config)
	if err != nil {
		psiphon.NoticeError("error initializing datastore: %s", err)
		os.Exit(1)
	}
	defer psiphon.CloseDataStore()

	controller, err := psiphon.NewController(config)
	if err != nil {
		psiphon.NoticeError("error creating controller: %s", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		controller.Run(ctx)
	}()

	select {
	case <-sigChan:
		psiphon.NoticeInfo("shutdown by system signal")
		cancel()
	case <-ctx.Done():
		psiphon.NoticeInfo("shutdown by context")
	}

	wg.Wait()

	psiphon.NoticeInfo("inproxy node stopped")
}
