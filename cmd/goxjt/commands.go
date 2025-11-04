// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/goxjt/pkg/goxjt"
	"github.com/beyondcivic/goxjt/pkg/version"
	"github.com/spf13/cobra"
)

// Version Command.
// Displays tool version and build information.
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the goxjt tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			stamp := version.RetrieveStamp()
			fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
			fmt.Printf("  Git ref: %s\n", stamp.VCSRevision)
			fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", stamp.InfoGoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		},
	}
}

func MapCmd() *cobra.Command {
	var mapCmd = &cobra.Command{
		Use:   "map [xmlPath] [schemaPath]",
		Short: "Map XML to JSON based on schema",
		Long: `Maps an XML document to a JSON object based on a user-defined schema.
		
Takes two arguments:
  xmlPath    - Path to the source XML file
  schemaPath - Path to the JSON schema file that defines the mapping`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			xmlPath := args[0]
			schemaPath := args[1]
			flagOutputPath, _ := cmd.Flags().GetString("output")

			// Validate input files
			if !fileExists(xmlPath) {
				fmt.Printf("Error: XML file '%s' does not exist.\n", xmlPath)
				os.Exit(1)
			}

			if !fileExists(schemaPath) {
				fmt.Printf("Error: Schema file '%s' does not exist.\n", schemaPath)
				os.Exit(1)
			}

			// Read XML file
			xmlData, err := os.ReadFile(xmlPath)
			if err != nil {
				fmt.Printf("Error reading XML file '%s': %v\n", xmlPath, err)
				os.Exit(1)
			}

			// Read schema file
			schemaData, err := os.ReadFile(schemaPath)
			if err != nil {
				fmt.Printf("Error reading schema file '%s': %v\n", schemaPath, err)
				os.Exit(1)
			}

			// Map XML to JSON
			fmt.Printf("Mapping XML '%s' to JSON using schema '%s'...\n", xmlPath, schemaPath)
			jsonData, err := goxjt.MapXMLToJSON(xmlData, schemaData)
			if err != nil {
				fmt.Printf("Error mapping XML to JSON: %v\n", err)
				os.Exit(1)
			}

			// Output result
			if flagOutputPath != "" {
				err := os.WriteFile(flagOutputPath, jsonData, 0644)
				if err != nil {
					fmt.Printf("Error writing output file '%s': %v\n", flagOutputPath, err)
					os.Exit(1)
				}
				fmt.Printf("✓ JSON output saved to: %s\n", flagOutputPath)
			} else {
				fmt.Printf("✓ JSON output:\n%s\n", string(jsonData))
			}
		},
	}

	mapCmd.Flags().StringP("output", "o", "", "Output path for the JSON file (if not specified, prints to stdout)")

	return mapCmd
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
