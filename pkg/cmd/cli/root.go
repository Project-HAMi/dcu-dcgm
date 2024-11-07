package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

var dcgmInitialized bool // 追踪 DCGM 是否成功初始化

var rootCmd = &cobra.Command{
	Use:   "dcgm",
	Short: "DCGM CLI tool",
	Long:  "Command-line interface for managing and interacting with DCGM. Use dcgm-cli [command] --help for more information on a command.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 在执行任何命令之前运行初始化
		if err := dcgm.Init(); err != nil {
			return fmt.Errorf("initialization failed: %v", err)
		}
		dcgmInitialized = true // 表示初始化成功
		return nil
	},
}

// Execute 执行 root 命令
func Execute() {
	defer func() {
		// 仅当 DCGM 成功初始化时才调用 ShutDown
		if dcgmInitialized {
			if err := dcgm.ShutDown(); err != nil {
				fmt.Println("Failed to shut down properly:", err)
			}
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
