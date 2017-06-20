package cmd

import (
	"os"

	"github.com/ellcrys/util"
	"github.com/ellcrys/cocoon/core/common"
	"github.com/ellcrys/patchain/cockroach"
	"github.com/ellcrys/safehold/config"
	"github.com/ellcrys/safehold/servers"
	"github.com/spf13/cobra"
)

var (

	// address to bind RPC server to
	bindAddrRPC = util.Env("SH_RPC_ADDR", "localhost:9002")

	// address to bind HTTP server to
	bindAddrHTTP = util.Env("SH_HTTP_ADDR", "localhost:9001")

	// database name of the global patchain database
	databaseName = util.Env("SH_PATCHAIN_DB_NAME", "safehold_dev")

	// partitionChainConStr connection string
	partitionChainConStr = util.Env("SH_PATCHAIN_CONSTR", "postgresql://root@localhost:26257/"+databaseName+"?sslmode=disable")
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
		if envNotSet := common.HasEnv([]string{
			"AUTH_SECRET",
		}...); len(envNotSet) > 0 {
			log.Fatalf("The following environment variable must be set: %v", envNotSet)
		}

		log.Infof("Running in '%s' environment", util.Env("ENV", "development"))

		db := cockroach.NewDB()
		db.ConnectionString = partitionChainConStr
		rpcServer := servers.NewRPC(db)
		httpServer := servers.NewHTTP(db, bindAddrRPC)

		// terminate app gracefully
		util.OnTerminate(func(s os.Signal) {
			log.Info("Termination signal received. Gracefully shutting down")
			rpcServer.Stop()
			db.Close()
			os.Exit(0)
		})

		// create rpc server, pass database implementation
		if err := rpcServer.Start(bindAddrRPC, func(s *servers.RPC) {
			httpServer.Start(bindAddrHTTP, nil)
		}); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
