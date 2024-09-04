package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"g.sugon.com/das/dcgm-dcu/pkg/dcgm"
)

var vDeviceInfoCmd = &cobra.Command{
	Use:   "vdevice-info [device-index]",
	Short: "Get virtual device information",
	Long:  `Retrieve detailed information about a virtual device using its device index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dvInd, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid device index:", err)
			os.Exit(1)
		}

		info, err := dcgm.VDeviceSingleInfo(dvInd)
		if err != nil {
			fmt.Println("Error fetching virtual device info:", err)
			os.Exit(1)
		}

		fmt.Printf("Virtual Device Info: %+v\n", info)
	},
}

var destroyVDeviceCmd = &cobra.Command{
	Use:   "destroy-vdevice<dvInd>",
	Short: "Destroy a single virtual device",
	Long:  `This command destroys a single virtual device by its index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vDvInd, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid virtual device index:", err)
			os.Exit(1)
		}

		err = dcgm.DestroySingleVDevice(vDvInd)
		if err != nil {
			fmt.Println("Error destroying virtual device:", err)
			os.Exit(1)
		}

		fmt.Printf("Virtual device %d destroyed successfully.\n", vDvInd)
	},
}
var allDeviceInfosCmd = &cobra.Command{
	Use:   "all-device-infos",
	Short: "Get information for all physical devices",
	Long:  `Retrieve detailed information about all physical devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		infos, err := dcgm.AllDeviceInfos()
		if err != nil {
			fmt.Println("Error fetching all device infos:", err)
			os.Exit(1)
		}
		fmt.Println("==========allDevices==========")
		fmt.Printf(dataToJson(infos))

	},
}

func init() {
	rootCmd.AddCommand(vDeviceInfoCmd)
	rootCmd.AddCommand(destroyVDeviceCmd)
	rootCmd.AddCommand(allDeviceInfosCmd)
}
