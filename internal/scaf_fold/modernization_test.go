package scaf_fold

import (
	"os"
	"path/filepath"
	"testing"
)

// 默认生成不应强制用户配置或运行外部数据库。
func TestGenerateStandaloneProject(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "standalone")
	err := Generate(dir, TemplateData{ModuleName: "example.com/standalone", BinaryName: "standalone", ProjectName: "standalone"})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"app/main.go", "AGENTS.md", ".agents/PLANS.md", "docs/README.md", "docs/sqls/README.md", "docs/sqls/schema/README.md", "docs/sqls/unreleased/README.md", "docs/sqls/releases/README.md", "docs/design/README.md", "docs/design/architecture/models.md", "docs/design/details/development/code-style.md", "internal/models/AGENTS.md", "internal/models/do/README.md", "internal/models/dto/README.md", "internal/models/vo/README.md", ".starter.json"} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("generated project missing %s: %v", path, err)
		}
	}
	for _, path := range []string{"docs/design/architecture/packages.md", "internal/controller/AGENTS.md", "internal/service/AGENTS.md", "internal/repository/AGENTS.md"} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("standalone project missing package guidance %s: %v", path, err)
		}
	}
}

// 选择数据库只选择基础设施，不应自动暴露示例业务。
func TestDatabaseSelectionDoesNotGenerateExampleBusiness(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "database")
	if err := Generate(dir, TemplateData{ModuleName: "example.com/database", BinaryName: "database", ProjectName: "database", MySQL: true, MongoDB: true}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"internal/controller/example_controller", "internal/service/example_srv", "docs/schema/users.sql"} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Errorf("unrequested example generated: %s (stat: %v)", path, err)
		}
	}
}

func TestGenerateRejectsSymlinkOutput(t *testing.T) {
	root := t.TempDir()
	destination := filepath.Join(root, "destination")
	if err := os.Mkdir(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(destination, link); err != nil {
		t.Fatal(err)
	}
	err := Generate(link, TemplateData{ModuleName: "example.com/safe", BinaryName: "safe", ProjectName: "safe", MySQL: true})
	if err == nil {
		t.Fatal("generation accepted a symlink output")
	}
	entries, err := os.ReadDir(destination)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatal("rejected generation wrote through symlink")
	}
}
