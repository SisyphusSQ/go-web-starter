package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-web-starter/v2/internal/scaf_fold"
)

func initReference() {
	var target string
	var write bool
	command := &cobra.Command{Use: "reference", Short: "Compare or sync the fixed go-starter reference (preserves local edits)", RunE: func(cmd *cobra.Command, _ []string) error {
		if target == "" {
			return fmt.Errorf("--target is required")
		}
		temp, err := os.MkdirTemp("", "go-starter-reference-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(temp)
		output := filepath.Join(temp, "go-starter")
		data := scaf_fold.TemplateData{ModuleName: "github.com/SisyphusSQ/go-starter/v2", AppVersion: "v2.0.0", BinaryName: "go-starter", ProjectName: "go-starter", MySQL: true, MongoDB: true, Redis: true, Cron: true, Lark: true, Prometheus: true, JWT: true, IssueProvider: "linear"}
		if err := scaf_fold.Generate(output, data); err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
		defer cancel()
		tidy := exec.CommandContext(ctx, "go", "mod", "tidy")
		tidy.Dir = output
		tidy.Stdout = cmd.OutOrStdout()
		tidy.Stderr = cmd.ErrOrStderr()
		if err := tidy.Run(); err != nil {
			return fmt.Errorf("resolve reference dependencies: %w", err)
		}
		if err := scaf_fold.RefreshManifest(output); err != nil {
			return err
		}
		changes, err := scaf_fold.SyncReference(output, target, write)
		if err != nil {
			return err
		}
		for _, change := range changes {
			fmt.Fprintln(cmd.OutOrStdout(), change)
		}
		if !write && len(changes) > 0 {
			return fmt.Errorf("reference differs in %d files", len(changes))
		}
		fmt.Fprintf(cmd.OutOrStdout(), "reference synchronized: %t; changes: %d\n", write, len(changes))
		return nil
	}}
	command.Flags().StringVar(&target, "target", "", "absolute reference repository directory")
	command.Flags().BoolVar(&write, "write", false, "apply reviewed generated changes")
	rootCmd.AddCommand(command)
}
