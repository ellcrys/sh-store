package cmd

import (
	"os"

	"github.com/ellcrys/util"
	"github.com/ncodes/bifrost/config"
	"github.com/ncodes/bifrost/servers"
	"github.com/ncodes/cocoon/core/common"
	"github.com/spf13/cobra"
)

var (

	// address to bind RPC server to
	bindAddrRPC = util.Env("SH_RPC_ADDR", "localhost:10008")

	// address to bind HTTP server to
	bindAddrHTTP = util.Env("SH_HTTP_ADDR", "localhost:10001")
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts an agent",
	Long:  `Starts an agent`,
	Run: func(cmd *cobra.Command, args []string) {

		var log = config.MakeLogger("rpc")
		defer log.Info("Stopped")

		// enforce required environment variables
		if envNotSet := common.HasEnv([]string{}...); len(envNotSet) > 0 {
			log.Fatalf("The following environment variable must be set: %v", envNotSet)
		}

		log.Infof("Running in '%s' environment", util.Env("ENV", "development"))

		var rpcServer = servers.NewRPC()

		// terminate app gracefully
		util.OnTerminate(func(s os.Signal) {
			log.Info("Termination signal received. Gracefully shutting down")
			rpcServer.Stop()
			os.Exit(0)
		})

		// create rpc server, pass database implementation
		if err := rpcServer.Start(bindAddrRPC, func(s *servers.RPC) {
		}); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
