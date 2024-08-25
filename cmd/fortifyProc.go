package cmd

import (
	"checksec/pkg/checksec"
	"checksec/pkg/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// fortifyProcCmd represents the fortifyProc command
var fortifyProcCmd = &cobra.Command{
	Use:   "fortifyProc",
	Short: "Check Fortify for running process",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Error: no process id provided")
			os.Exit(1)
		}
		proc := args[0]

		file, err := os.Readlink(filepath.Join("/proc", proc, "exe"))
		if err != nil {
			fmt.Printf("Error: Pid %s not found", proc)
			os.Exit(1)
		}

		utils.CheckElfExists(file)
		binary := utils.GetBinary(file)
		fortify := checksec.Fortify(file, binary)
		output := []interface{}{
			map[string]interface{}{
				"name": file,
				"checks": map[string]interface{}{
					"fortify_source": fortify.Output,
					"fortified":      fortify.Fortified,
					"fortifyable":    fortify.Fortifiable,
					"noFortify":      fortify.NoFortify,
					"libcSupport":    fortify.LibcSupport,
					"numLibcFunc":    fortify.NumLibcFunc,
					"numFileFunc":    fortify.NumFileFunc,
				},
			},
		}
		color := []interface{}{
			map[string]interface{}{
				"name": file,
				"checks": map[string]interface{}{
					"fortified":           fortify.Fortified,
					"fortifiedColor":      "unset",
					"noFortify":           fortify.NoFortify,
					"fortifyable":         fortify.Fortifiable,
					"fortifyableColor":    "unset",
					"fortify_source":      fortify.Output,
					"fortify_sourceColor": fortify.Color,
					"libcSupport":         fortify.LibcSupport,
					"libcSupportColor":    fortify.LibcSupportColor,
					"numLibcFunc":         fortify.NumLibcFunc,
					"numFileFunc":         fortify.NumFileFunc,
				},
			},
		}
		utils.FortifyPrinter(outputFormat, output, color)
	},
}

func init() {
	rootCmd.AddCommand(fortifyProcCmd)
}