package daemon

import (
	"fmt"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server"

	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(ConfigCmd)
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "run the daemon",
	Run:   ConfigRun,
}

func ConfigRun(cmd *cobra.Command, args []string) {
	fmt.Println(GetConfig())
}

const CfgFormat = `Configuration:
  Database:
    Username:		%s
    Password:		%s
    Host:		%s
    Port:		%s
    Database:		%s
  API:
    Host:		%s
    Port:		%s
    Static Type:	%s`

func GetConfig() string {
	dbConfig, _ := db.LoadConfig()
	apiConfig, _ := server.LoadConfig()
	return fmt.Sprintf(CfgFormat, dbConfig.GetUsername(), dbConfig.GetPassword(), dbConfig.GetHost(), dbConfig.GetPort(), dbConfig.GetDatabase(), apiConfig.Host(), apiConfig.Port(), apiConfig.StaticType())
}
