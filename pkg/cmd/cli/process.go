package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

var pidListCmd = &cobra.Command{
	Use:   "pid-list",
	Short: "Get a list of PIDs",
	Long:  `Retrieve a list of process IDs (PIDs) managed by the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		pidList, err := dcgm.PidList()
		if err != nil {
			fmt.Println("Error fetching PID list:", err)
			os.Exit(1)
		}

		fmt.Println("PID List:")
		for _, pid := range pidList {
			fmt.Println(pid)
		}
	},
}

var showPidsCmd = &cobra.Command{
	Use:   "show-pids",
	Short: "Show running KFD process information",
	Long:  `Retrieve and display detailed information about KFD processes currently running on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 调用 ShowPids 函数
		err := dcgm.ShowPids()
		if err != nil {
			fmt.Println("Error displaying KFD process information:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pidListCmd)
	rootCmd.AddCommand(showPidsCmd)
}
