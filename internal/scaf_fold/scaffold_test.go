package scaf_fold

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const commandTimeout = 5 * time.Minute

func TestGenerate(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := filepath.Join(baseDir, "my-new-web")
	data := TemplateData{
		ModuleName:  "github.com/test/my-new-web",
		BinaryName:  "my-new-web",
		ProjectName: "my-new-web",
		MySQL:       true,
		MongoDB:     true,
	}

	if err := Generate(outputDir, data); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	goModPath := filepath.Join(outputDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}
	if !strings.Contains(string(goModContent), "module github.com/test/my-new-web") {
		t.Fatalf("go.mod does not contain expected module path: %s", goModContent)
	}
	if !strings.Contains(string(goModContent), "\ngo "+defaultGoVersion()+"\n") {
		t.Fatalf("go.mod does not contain expected go version directive: %s", goModContent)
	}

	for _, relPath := range []string{
		filepath.Join("app", "main.go"),
		filepath.Join("app", "cmd", "root.go"),
		filepath.Join("app", "cmd", "http.go"),
		filepath.Join("internal", "http", "server.go"),
		filepath.Join("internal", "lib", "gorm", "gorm.go"),
		filepath.Join("internal", "lib", "mongodb", "mongodb.go"),
		filepath.Join("vars", "vars.go"),
		filepath.Join("README.md"),
	} {
		assertFileExists(t, filepath.Join(outputDir, relPath))
	}
}

func TestGenerateRejectsNonEmptyOutputDir(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := filepath.Join(baseDir, "existing")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "placeholder.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write placeholder file: %v", err)
	}

	err := Generate(outputDir, TemplateData{
		ModuleName:  "github.com/test/sample",
		BinaryName:  "sample",
		ProjectName: "sample",
		MySQL:       true,
		MongoDB:     true,
	})
	if err == nil {
		t.Fatal("Generate() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("Generate() unexpected error: %v", err)
	}
}

func TestGenerateAllowsHiddenEntriesOnly(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := filepath.Join(baseDir, "existing")

	if err := os.MkdirAll(filepath.Join(outputDir, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git dir: %v", err)
	}

	err := Generate(outputDir, TemplateData{
		ModuleName:  "github.com/test/sample",
		BinaryName:  "sample",
		ProjectName: "sample",
		MySQL:       true,
		MongoDB:     false,
	})
	if err != nil {
		t.Fatalf("Generate() expected success with hidden-only entries, got error: %v", err)
	}
}

func TestGenerateRejectsHiddenEntriesExceptGit(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := filepath.Join(baseDir, "existing")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, ".gitignore"), []byte(""), 0o644); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}

	err := Generate(outputDir, TemplateData{
		ModuleName:  "github.com/test/sample",
		BinaryName:  "sample",
		ProjectName: "sample",
		MySQL:       true,
		MongoDB:     true,
	})
	if err == nil {
		t.Fatal("Generate() expected error for hidden non-.git entry, got nil")
	}
	if !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("Generate() unexpected error: %v", err)
	}
}

func TestGenerateRejectsInvalidTemplateData(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "out")
	err := Generate(outputDir, TemplateData{
		ModuleName:  "bad module",
		BinaryName:  "my-web",
		ProjectName: "my-web",
		MySQL:       true,
		MongoDB:     true,
	})
	if err == nil {
		t.Fatal("Generate() expected error for invalid module, got nil")
	}
	if !strings.Contains(err.Error(), "invalid template data") {
		t.Fatalf("Generate() unexpected error: %v", err)
	}
}

func TestGenerateE2EDBCombos(t *testing.T) {
	testGenerateE2EDBCombos(t, false)
}

func testGenerateE2EDBCombos(t *testing.T, runBuildChecks bool) {
	tests := []struct {
		name      string
		mysql     bool
		mongodb   bool
		module    string
		binary    string
		project   string
		leakCheck func(*testing.T, string)
	}{
		{
			name:    "mysql-only",
			mysql:   true,
			mongodb: false,
			module:  "github.com/test/mysql-only-web",
			binary:  "mysql-only-web",
			project: "mysql-only-web",
			leakCheck: func(t *testing.T, outputDir string) {
				t.Helper()
				goMod := readFileForAssertion(t, filepath.Join(outputDir, "go.mod"))
				if strings.Contains(goMod, "github.com/qiniu/qmgo") {
					t.Fatalf("mysql-only go.mod should not contain qmgo dependency")
				}
				if strings.Contains(goMod, "go.mongodb.org/mongo-driver") {
					t.Fatalf("mysql-only go.mod should not contain mongo-driver dependency")
				}
			},
		},
		{
			name:    "mongodb-only",
			mysql:   false,
			mongodb: true,
			module:  "github.com/test/mongodb-only-web",
			binary:  "mongodb-only-web",
			project: "mongodb-only-web",
			leakCheck: func(t *testing.T, outputDir string) {
				t.Helper()
				goMod := readFileForAssertion(t, filepath.Join(outputDir, "go.mod"))
				if strings.Contains(goMod, "gorm.io/gorm") {
					t.Fatalf("mongodb-only go.mod should not contain gorm dependency")
				}
				if strings.Contains(goMod, "gorm.io/driver/mysql") {
					t.Fatalf("mongodb-only go.mod should not contain mysql gorm driver dependency")
				}
				if strings.Contains(goMod, "github.com/go-sql-driver/mysql") {
					t.Fatalf("mongodb-only go.mod should not contain mysql driver dependency")
				}
			},
		},
		{
			name:    "mysql-and-mongodb",
			mysql:   true,
			mongodb: true,
			module:  "github.com/test/full-web",
			binary:  "full-web",
			project: "full-web",
			leakCheck: func(t *testing.T, outputDir string) {
				t.Helper()
				goMod := readFileForAssertion(t, filepath.Join(outputDir, "go.mod"))
				if !strings.Contains(goMod, "gorm.io/gorm") {
					t.Fatalf("full mode go.mod should contain gorm dependency")
				}
				if !strings.Contains(goMod, "github.com/qiniu/qmgo") {
					t.Fatalf("full mode go.mod should contain qmgo dependency")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDir := t.TempDir()
			outputDir := filepath.Join(baseDir, tt.project)
			data := TemplateData{
				ModuleName:  tt.module,
				BinaryName:  tt.binary,
				ProjectName: tt.project,
				MySQL:       tt.mysql,
				MongoDB:     tt.mongodb,
			}

			if err := Generate(outputDir, data); err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			assertFileExists(t, filepath.Join(outputDir, "app", "main.go"))
			assertFileExists(t, filepath.Join(outputDir, "config", "config.go"))
			assertFileExists(t, filepath.Join(outputDir, "internal", "service", "module.go"))
			assertFileExists(t, filepath.Join(outputDir, "internal", "repository", "module.go"))
			assertFileExists(t, filepath.Join(outputDir, "internal", "controller", "module.go"))

			if tt.mysql {
				assertFileExists(t, filepath.Join(outputDir, "internal", "lib", "gorm", "gorm.go"))
				assertFileExists(t, filepath.Join(outputDir, "internal", "controller", "example_controller", "user_handler.go"))
			} else {
				assertFileNotExists(t, filepath.Join(outputDir, "internal", "lib", "gorm", "gorm.go"))
				assertFileNotExists(t, filepath.Join(outputDir, "internal", "controller", "example_controller", "user_handler.go"))
			}

			if tt.mongodb {
				assertFileExists(t, filepath.Join(outputDir, "internal", "lib", "mongodb", "mongodb.go"))
				assertFileExists(t, filepath.Join(outputDir, "internal", "controller", "example_controller", "user_mongo_handler.go"))
			} else {
				assertFileNotExists(t, filepath.Join(outputDir, "internal", "lib", "mongodb", "mongodb.go"))
				assertFileNotExists(t, filepath.Join(outputDir, "internal", "controller", "example_controller", "user_mongo_handler.go"))
			}

			if runBuildChecks {
				if stdout, stderr, err := runGoCommand(outputDir, "mod", "tidy"); err != nil {
					t.Fatalf(
						"go mod tidy failed: %v\nstdout:\n%s\nstderr:\n%s",
						err,
						stdout,
						stderr,
					)
				}
				if stdout, stderr, err := runGoCommand(outputDir, "build", "./..."); err != nil {
					t.Fatalf("go build failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
				}
				if stdout, stderr, err := runGoCommand(outputDir, "vet", "./..."); err != nil {
					t.Fatalf("go vet failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
				}
			}

			tt.leakCheck(t, outputDir)
		})
	}
}

func runGoCommand(dir string, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() != nil {
		return stdout.String(), stderr.String(), fmt.Errorf(
			"go %s: %w",
			strings.Join(args, " "),
			ctx.Err(),
		)
	}

	return stdout.String(), stderr.String(), err
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %s, err=%v", path, err)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file to be absent: %s, err=%v", path, err)
	}
}

func readFileForAssertion(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(data)
}
