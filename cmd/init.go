package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/st2projects/ssh-sentinel-client/config"
	"os"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the client",
	Run: func(cmd *cobra.Command, args []string) {

		configBytes, err := json.MarshalIndent(&config.ConfigType{}, "", "    ")

		if err != nil {
			panic(err)
		}

		userHome, err := os.UserHomeDir()

		if err != nil {
			panic(err)
		}

		configPath := filepath.Join(userHome, ".ssh-sentinel.json")
		err = os.WriteFile(configPath, configBytes, os.FileMode(0600))
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
