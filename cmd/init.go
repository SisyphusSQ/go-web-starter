package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-web-starter/internal/scaf_fold"
)

var (
	initModuleNameFlag string
	initBinaryNameFlag string
	initDBFlag         string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project in current directory (.git allowed)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := buildTemplateData(
			".",
			initModuleNameFlag,
			initBinaryNameFlag,
			initDBFlag,
		)
		if err != nil {
			return err
		}

		if err := scaf_fold.Generate(".", data); err != nil {
			return fmt.Errorf("initialize project: %w", err)
		}

		fmt.Println("Project initialized in current directory")
		fmt.Println()
		printNextSteps(".", false)
		return nil
	},
}

func initInit() {
	initCmd.Flags().StringVarP(
		&initModuleNameFlag,
		"module",
		"m",
		"",
		"Go module path (default: example.com/<directory-name>)",
	)
	initCmd.Flags().StringVarP(
		&initBinaryNameFlag,
		"binary",
		"b",
		"",
		"Binary name (default: inferred from directory name)",
	)
	initCmd.Flags().StringVar(
		&initDBFlag,
		"db",
		"mysql,mongodb",
		"Database engines: mysql, mongodb, or mysql,mongodb",
	)

	rootCmd.AddCommand(initCmd)
}
