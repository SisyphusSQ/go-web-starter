package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-web-starter/v2/internal/scaf_fold"
)

var (
	initModuleNameFlag    string
	initBinaryNameFlag    string
	initDBFlag            string
	initWithFlag          string
	initExamplesFlag      bool
	initIssueProviderFlag string
	initIssuePrefixFlag   string
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

		data.Examples = initExamplesFlag
		data.IssueProvider = initIssueProviderFlag
		data.IssuePrefix = initIssuePrefixFlag
		if err := data.SetFeatures(initWithFlag); err != nil {
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
		"none",
		"Database engines: none, mysql, mongodb, or mysql,mongodb",
	)

	initCmd.Flags().StringVar(&initWithFlag, "with", "", "Optional components: redis,cron,lark,prometheus-query,jwt")
	initCmd.Flags().BoolVar(&initExamplesFlag, "examples", false, "Include complete User CRUD example (requires a database)")
	initCmd.Flags().StringVar(&initIssueProviderFlag, "issue-provider", "linear", "Issue provider: linear,github,gitlab,repo,other")
	initCmd.Flags().StringVar(&initIssuePrefixFlag, "issue-prefix", "", "Issue identifier prefix")
	rootCmd.AddCommand(initCmd)
}
