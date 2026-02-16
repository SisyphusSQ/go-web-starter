package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootExecuteNewRequiresOutputDir(t *testing.T) {
	err := executeRootForTest(nil, "new")
	if err == nil {
		t.Fatal("expected error for missing output-dir, got nil")
	}
	if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootExecuteNewRejectsInvalidDBFlag(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "demo")
	err := executeRootForTest(nil, "new", outDir, "--db", "postgres")
	if err == nil {
		t.Fatal("expected invalid db error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid db value") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootExecuteNewUsesDefaultDBFlag(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "demo-default-db")
	if err := executeRootForTest(nil, "new", outDir); err != nil {
		t.Fatalf("execute new with default db failed: %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}
	modText := string(goMod)
	if !strings.Contains(modText, "gorm.io/gorm") {
		t.Fatalf("default db should include mysql dependency, go.mod:\n%s", modText)
	}
	if !strings.Contains(modText, "github.com/qiniu/qmgo") {
		t.Fatalf("default db should include mongodb dependency, go.mod:\n%s", modText)
	}
}

func TestRootExecuteInitRejectsArgs(t *testing.T) {
	err := executeRootForTest(nil, "init", "unexpected")
	if err == nil {
		t.Fatal("expected error for extra init arg, got nil")
	}
	if !strings.Contains(err.Error(), "unknown command") &&
		!strings.Contains(err.Error(), "accepts 0 arg(s), received 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootExecuteNewHelpContainsDBFlag(t *testing.T) {
	var output bytes.Buffer
	if err := executeRootForTest(&output, "new", "--help"); err != nil {
		t.Fatalf("execute new --help failed: %v", err)
	}
	if !strings.Contains(output.String(), "Database engines: mysql, mongodb, or mysql,mongodb") {
		t.Fatalf("help output missing db flag description:\n%s", output.String())
	}
}

func executeRootForTest(output *bytes.Buffer, args ...string) error {
	ensureInitialized()
	resetCLIFlagStateForTest()

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	if output == nil {
		rootCmd.SetOut(bytes.NewBuffer(nil))
		rootCmd.SetErr(bytes.NewBuffer(nil))
	} else {
		output.Reset()
		rootCmd.SetOut(output)
		rootCmd.SetErr(output)
	}
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func resetCLIFlagStateForTest() {
	moduleNameFlag = ""
	binaryNameFlag = ""
	dbFlag = "mysql,mongodb"
	initModuleNameFlag = ""
	initBinaryNameFlag = ""
	initDBFlag = "mysql,mongodb"
}
