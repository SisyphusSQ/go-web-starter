package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SisyphusSQ/go-web-starter/internal/scaf_fold"
)

func TestInferProjectNameFromCurrentDirectory(t *testing.T) {
	baseDir := t.TempDir()
	projectDir := filepath.Join(baseDir, "demo-app")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir project dir: %v", err)
	}

	projectName, err := inferProjectName(".")
	if err != nil {
		t.Fatalf("inferProjectName(.) error = %v", err)
	}
	if projectName != "demo-app" {
		t.Fatalf("inferProjectName(.) = %q, want %q", projectName, "demo-app")
	}
}

func TestInferProjectNameRejectsIllegalDefaults(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
	}{
		{name: "root path", outputDir: string(filepath.Separator)},
		{name: "parent directory path", outputDir: ".."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := inferProjectName(tt.outputDir)
			if err == nil {
				t.Fatalf("inferProjectName(%q) expected error, got nil", tt.outputDir)
			}
			if !strings.Contains(err.Error(), "cannot infer project name") {
				t.Fatalf("inferProjectName(%q) error = %v", tt.outputDir, err)
			}
		})
	}
}

func TestBuildTemplateDataWithDefaults(t *testing.T) {
	baseDir := t.TempDir()
	projectDir := filepath.Join(baseDir, "web-demo")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	data, err := buildTemplateData(projectDir, "", "", "mysql,mongodb")
	if err != nil {
		t.Fatalf("buildTemplateData() error = %v", err)
	}

	if data.ProjectName != "web-demo" {
		t.Fatalf("ProjectName = %q, want %q", data.ProjectName, "web-demo")
	}
	if data.BinaryName != "web-demo" {
		t.Fatalf("BinaryName = %q, want %q", data.BinaryName, "web-demo")
	}
	if data.ModuleName != "example.com/web-demo" {
		t.Fatalf("ModuleName = %q, want %q", data.ModuleName, "example.com/web-demo")
	}
	if !data.MySQL || !data.MongoDB {
		t.Fatalf("defaults should enable both dbs, got mysql=%v mongodb=%v", data.MySQL, data.MongoDB)
	}
}

func TestBuildTemplateDataWithCustomValues(t *testing.T) {
	data, err := buildTemplateData(
		"./anything",
		"github.com/acme/web",
		"acme-web",
		"mongodb",
	)
	if err != nil {
		t.Fatalf("buildTemplateData() error = %v", err)
	}

	if data.ModuleName != "github.com/acme/web" {
		t.Fatalf("ModuleName = %q, want %q", data.ModuleName, "github.com/acme/web")
	}
	if data.BinaryName != "acme-web" {
		t.Fatalf("BinaryName = %q, want %q", data.BinaryName, "acme-web")
	}
	if data.MySQL || !data.MongoDB {
		t.Fatalf("db selection mismatch, got mysql=%v mongodb=%v", data.MySQL, data.MongoDB)
	}
}

func TestParseDBFlagBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		val       string
		wantMySQL bool
		wantMongo bool
		wantErr   bool
	}{
		{
			name:      "default both",
			val:       "mysql,mongodb",
			wantMySQL: true,
			wantMongo: true,
		},
		{
			name:      "case insensitive and spaces",
			val:       " MySQL , MONGODB ",
			wantMySQL: true,
			wantMongo: true,
		},
		{
			name:      "ignore empty token",
			val:       "mysql,,mongodb",
			wantMySQL: true,
			wantMongo: true,
		},
		{
			name:      "single mysql with empty token",
			val:       ", mysql, ",
			wantMySQL: true,
			wantMongo: false,
		},
		{
			name:    "empty after trim",
			val:     " , , ",
			wantErr: true,
		},
		{
			name:    "invalid token",
			val:     "mysql,postgres",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mysql, mongo, err := scaf_fold.ParseDBFlag(tt.val)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseDBFlag(%q) expected error, got nil", tt.val)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDBFlag(%q) error = %v", tt.val, err)
			}
			if mysql != tt.wantMySQL || mongo != tt.wantMongo {
				t.Fatalf(
					"ParseDBFlag(%q) = (mysql=%v,mongo=%v), want (mysql=%v,mongo=%v)",
					tt.val,
					mysql,
					mongo,
					tt.wantMySQL,
					tt.wantMongo,
				)
			}
		})
	}
}

func TestInitCommandInEmptyDir(t *testing.T) {
	baseDir := t.TempDir()
	projectDir := filepath.Join(baseDir, "init-demo")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(projectDir, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git dir: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir project dir: %v", err)
	}

	origModule := initModuleNameFlag
	origBinary := initBinaryNameFlag
	origDB := initDBFlag
	t.Cleanup(func() {
		initModuleNameFlag = origModule
		initBinaryNameFlag = origBinary
		initDBFlag = origDB
	})

	initModuleNameFlag = ""
	initBinaryNameFlag = ""
	initDBFlag = "mysql"

	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("initCmd.RunE() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectDir, "app", "main.go")); err != nil {
		t.Fatalf("expected app/main.go to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "internal", "lib", "gorm", "gorm.go")); err != nil {
		t.Fatalf("expected mysql gorm file to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "internal", "lib", "mongodb", "mongodb.go")); !os.IsNotExist(err) {
		t.Fatalf("expected mongodb file to be absent in mysql mode, got err=%v", err)
	}
}
