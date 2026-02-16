package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitAndNewDBBehaviorConsistency(t *testing.T) {
	baseDir := t.TempDir()
	newOut := filepath.Join(baseDir, "from-new")
	initOut := filepath.Join(baseDir, "from-init")

	if err := os.MkdirAll(initOut, 0o755); err != nil {
		t.Fatalf("mkdir initOut: %v", err)
	}
	if err := os.Mkdir(filepath.Join(initOut, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir initOut .git: %v", err)
	}

	origModule := moduleNameFlag
	origBinary := binaryNameFlag
	origDB := dbFlag
	origInitModule := initModuleNameFlag
	origInitBinary := initBinaryNameFlag
	origInitDB := initDBFlag
	t.Cleanup(func() {
		moduleNameFlag = origModule
		binaryNameFlag = origBinary
		dbFlag = origDB
		initModuleNameFlag = origInitModule
		initBinaryNameFlag = origInitBinary
		initDBFlag = origInitDB
	})

	moduleNameFlag = "github.com/acme/from-new"
	binaryNameFlag = "from-new"
	dbFlag = "mongodb"
	if err := newCmd.RunE(newCmd, []string{newOut}); err != nil {
		t.Fatalf("newCmd.RunE() error = %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})
	if err := os.Chdir(initOut); err != nil {
		t.Fatalf("chdir initOut: %v", err)
	}

	initModuleNameFlag = "github.com/acme/from-init"
	initBinaryNameFlag = "from-init"
	initDBFlag = "mongodb"
	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("initCmd.RunE() error = %v", err)
	}

	assertMongoOnlyLayout := func(dir string) {
		t.Helper()
		if _, err := os.Stat(filepath.Join(dir, "internal", "lib", "mongodb", "mongodb.go")); err != nil {
			t.Fatalf("expected mongodb file in %s: %v", dir, err)
		}
		if _, err := os.Stat(filepath.Join(dir, "internal", "lib", "gorm", "gorm.go")); !os.IsNotExist(err) {
			t.Fatalf("expected mysql gorm file absent in %s, got err=%v", dir, err)
		}

		goMod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err != nil {
			t.Fatalf("read go.mod in %s: %v", dir, err)
		}
		modText := string(goMod)
		if !strings.Contains(modText, "github.com/qiniu/qmgo") {
			t.Fatalf("go.mod in %s should include qmgo dependency", dir)
		}
		if strings.Contains(modText, "gorm.io/gorm") {
			t.Fatalf("go.mod in %s should not include gorm dependency in mongodb mode", dir)
		}
	}

	assertMongoOnlyLayout(newOut)
	assertMongoOnlyLayout(initOut)
}
