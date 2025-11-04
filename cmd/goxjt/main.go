package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/goxjt/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root cobra command.
// Call Init() once to initialize child commands.
// Global so it can be picked up by docs/doc-gen.go.
// nolint:gochecknoglobals
var RootCmd = &cobra.Command{
	Use:     "goxjt",
	Short:   "XML mapping tool",
	Long:    `A Go implementation for working with xml using json templates.`,
	Version: version.Version,
}

// Call Once.
func Init() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("GOXJT")
	viper.AutomaticEnv()

	// Add child commands
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(MapCmd())
}

func Execute() {
	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
