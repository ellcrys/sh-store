package cmd

import (
	"os"

	"github.com/ellcrys/cocoon/core/common"
	"github.com/ellcrys/elldb/config"
	"github.com/ellcrys/elldb/core/servers"
	"github.com/ellcrys/util"
	"github.com/spf13/cobra"
)

var (

	// address to bind RPC server to
	bindAddrRPC = util.Env("ELLDB_RPC_ADDR", "localhost:9002")

	// address to bind HTTP server to
	bindAddrHTTP = util.Env("ELLDB_HTTP_ADDR", "localhost:9001")

	// requiredEnv includes the environment variables required to start
	requiredEnv = []string{
		"AUTH_SECRET",
	}
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an ellcrys database server",
	Long:  `Start an ellcrys database server`,
	Run: func(cmd *cobra.Command, args []string) {

		var log = config.MakeLogger("rpc")
		defer log.Info("Stopped")

		// enforce required environment variables
		if envNotSet := common.HasEnv(requiredEnv...); len(envNotSet) > 0 {
			log.Fatalf("The following environment variable are required: %v", envNotSet)
		}

		log.Infof("Running in '%s' environment", util.Env("ENV", "development"))

		rpcServer := servers.NewRPC()

		// terminate app gracefully
		util.OnTerminate(func(s os.Signal) {
			log.Info("Termination signal received. Gracefully shutting down")
			rpcServer.Stop()
			os.Exit(0)
		})

		if err := rpcServer.Start(bindAddrRPC, func(s *servers.RPC) {
			httpServer := servers.NewHTTP(bindAddrRPC, rpcServer.GetDB())
			httpServer.Start(bindAddrHTTP, nil)
		}); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
