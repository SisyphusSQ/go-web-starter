package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-web-starter/internal/scaf_fold"
)

var (
	moduleNameFlag string
	binaryNameFlag string
	dbFlag         string
)

var newCmd = &cobra.Command{
	Use:   "new <output-dir>",
	Short: "Generate a new web project from template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := args[0]
		data, err := buildTemplateData(
			outputDir,
			moduleNameFlag,
			binaryNameFlag,
			dbFlag,
		)
		if err != nil {
			return err
		}

		if err := scaf_fold.Generate(outputDir, data); err != nil {
			return fmt.Errorf("generate project: %w", err)
		}

		fmt.Printf("Project generated at %s\n\n", outputDir)
		printNextSteps(outputDir, true)
		return nil
	},
}

func initNew() {
	newCmd.Flags().StringVarP(
		&moduleNameFlag,
		"module",
		"m",
		"",
		"Go module path (default: example.com/<directory-name>)",
	)
	newCmd.Flags().StringVarP(
		&binaryNameFlag,
		"binary",
		"b",
		"",
		"Binary name (default: inferred from directory name)",
	)
	newCmd.Flags().StringVar(
		&dbFlag,
		"db",
		"mysql,mongodb",
		"Database engines: mysql, mongodb, or mysql,mongodb",
	)

	rootCmd.AddCommand(newCmd)
}

func buildTemplateData(
	outputDir string,
	moduleValue string,
	binaryValue string,
	dbValue string,
) (scaf_fold.TemplateData, error) {
	projectName, err := inferProjectName(outputDir)
	if err != nil {
		return scaf_fold.TemplateData{}, err
	}

	moduleName := strings.TrimSpace(moduleValue)
	if moduleName == "" {
		moduleName = defaultModuleName(projectName)
	}

	binaryName := strings.TrimSpace(binaryValue)
	if binaryName == "" {
		binaryName = projectName
	}

	mysql, mongodb, err := scaf_fold.ParseDBFlag(dbValue)
	if err != nil {
		return scaf_fold.TemplateData{}, err
	}

	return scaf_fold.TemplateData{
		ModuleName:  moduleName,
		BinaryName:  binaryName,
		ProjectName: projectName,
		MySQL:       mysql,
		MongoDB:     mongodb,
	}, nil
}

func inferProjectName(outputDir string) (string, error) {
	cleanedOutputDir := filepath.Clean(strings.TrimSpace(outputDir))
	if cleanedOutputDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve current directory: %w", err)
		}
		return sanitizeDerivedProjectName(filepath.Base(wd), wd)
	}

	return sanitizeDerivedProjectName(
		filepath.Base(cleanedOutputDir),
		cleanedOutputDir,
	)
}

func sanitizeDerivedProjectName(projectName, source string) (string, error) {
	name := strings.TrimSpace(projectName)
	if name == "" || name == "." || name == ".." {
		return "", fmt.Errorf(
			"cannot infer project name from %q: use a non-root/non-dot output directory",
			source,
		)
	}
	if strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf(
			"cannot infer project name from %q: derived name %q is invalid",
			source,
			name,
		)
	}

	return name, nil
}

func defaultModuleName(projectName string) string {
	return fmt.Sprintf("example.com/%s", projectName)
}

func printNextSteps(outputDir string, includeCD bool) {
	fmt.Println("Next steps:")
	if includeCD {
		fmt.Printf("  cd %s\n", outputDir)
	}
	fmt.Println("  go mod tidy")
	fmt.Println("  # edit config/config.yml")
	fmt.Println("  go run ./app/main.go http")
}
