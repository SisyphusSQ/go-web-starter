package scaf_fold

import (
	"errors"
	"fmt"
	"go/format"
	goversion "go/version"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const (
	templateRoot             = "_template"
	stableGoVersion          = "1.27.1"
	minimumTemplateGoVersion = "1.27.0"
)

type TemplateData struct {
	ModuleName    string
	BinaryName    string
	ProjectName   string
	GoVersion     string
	AppVersion    string
	MySQL         bool
	MongoDB       bool
	Redis         bool
	Cron          bool
	Lark          bool
	Prometheus    bool
	JWT           bool
	Examples      bool
	IssueProvider string
	IssuePrefix   string
}

var (
	binaryNamePattern         = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	projectNamePattern        = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	goVersionPattern          = regexp.MustCompile(`^\d+\.\d+(?:\.\d+)?$`)
	mysqlOnlyTemplatePrefixes = []string{
		templateRoot + "/internal/lib/gorm",
		templateRoot + "/internal/lib/log/silent.go.tmpl",
		templateRoot + "/internal/repository/mysql",
		templateRoot + "/internal/models/do/mysql",
		templateRoot + "/internal/controller/example_controller/user_handler.go.tmpl",
		templateRoot + "/internal/service/example_srv/user_service.go.tmpl",
	}
	mongoOnlyTemplatePrefixes = []string{
		templateRoot + "/internal/lib/mongodb",
		templateRoot + "/internal/repository/mongo",
		templateRoot + "/internal/models/do/mongo",
		templateRoot + "/internal/controller/example_controller/user_mongo_handler.go.tmpl",
		templateRoot + "/internal/service/example_srv/user_mongo_service.go.tmpl",
	}
)

func (d TemplateData) Validate() error {
	if d.AppVersion != "" && d.AppVersion != "dev" && !semver.IsValid(d.AppVersion) {
		return fmt.Errorf("invalid application version %q", d.AppVersion)
	}
	if err := validateModulePath(d.ModuleName); err != nil {
		return err
	}
	if err := validateBinaryName(d.BinaryName); err != nil {
		return err
	}
	if err := validateProjectName(d.ProjectName); err != nil {
		return err
	}

	goVersion := strings.TrimSpace(d.GoVersion)
	if goVersion == "" {
		goVersion = defaultGoVersion()
	}
	if err := validateGoVersion(goVersion); err != nil {
		return err
	}
	if d.JWT && !d.Redis {
		return fmt.Errorf("jwt requires redis")
	}
	if d.Examples && !d.MySQL && !d.MongoDB {
		return fmt.Errorf("user examples require mysql or mongodb")
	}
	switch d.IssueProvider {
	case "", "linear", "github", "gitlab", "repo", "other":
	default:
		return fmt.Errorf("unsupported issue provider %q", d.IssueProvider)
	}
	if d.IssuePrefix != "" && !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`).MatchString(d.IssuePrefix) {
		return fmt.Errorf("invalid issue prefix")
	}

	return nil
}

func ParseDBFlag(val string) (mysql, mongodb bool, err error) {
	return parseDBFlag(val)
}

func parseDBFlag(val string) (mysql, mongodb bool, err error) {
	if strings.EqualFold(strings.TrimSpace(val), "none") {
		return false, false, nil
	}
	for _, rawToken := range strings.Split(val, ",") {
		token := strings.ToLower(strings.TrimSpace(rawToken))
		if token == "" {
			continue
		}

		switch token {
		case "mysql":
			mysql = true
		case "mongodb":
			mongodb = true
		default:
			return false, false, fmt.Errorf(
				"invalid db value %q: allowed values are none,mysql,mongodb",
				rawToken,
			)
		}
	}

	if !mysql && !mongodb {
		return false, false, fmt.Errorf("at least one database must be selected: mysql,mongodb")
	}

	return mysql, mongodb, nil
}

func Generate(outputDir string, data TemplateData) error {
	outputDir = filepath.Clean(outputDir)
	data.applyDefaults()
	if err := data.Validate(); err != nil {
		return fmt.Errorf("invalid template data: %w", err)
	}

	files := make(map[string][]byte)
	if err := fs.WalkDir(templateFS, templateRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == templateRoot {
			return nil
		}
		if shouldSkipTemplate(path, data) {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		rel := strings.TrimSuffix(strings.TrimPrefix(path, templateRoot+"/"), ".tmpl")
		raw, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		rendered, err := renderTemplate(path, raw, data)
		if err != nil {
			return err
		}
		if filepath.Ext(rel) == ".go" {
			rendered, err = format.Source(rendered)
			if err != nil {
				return fmt.Errorf("format %s: %w", path, err)
			}
		}
		files[rel] = rendered
		return nil
	}); err != nil {
		return fmt.Errorf("render project: %w", err)
	}
	manifest, err := buildManifest(data, files)
	if err != nil {
		return err
	}
	files[".starter.json"] = manifest
	if err := prepareOutputDir(outputDir); err != nil {
		return err
	}
	root, err := os.OpenRoot(outputDir)
	if err != nil {
		return err
	}
	defer root.Close()
	for _, path := range slices.Sorted(maps.Keys(files)) {
		if err := root.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		file, err := root.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return fmt.Errorf("create %s: %w", path, err)
		}
		_, writeErr := file.Write(files[path])
		closeErr := file.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}

func (d *TemplateData) applyDefaults() {
	if d.AppVersion == "" {
		d.AppVersion = "dev"
	}
	if strings.TrimSpace(d.GoVersion) == "" {
		d.GoVersion = defaultGoVersion()
	}
	if d.IssueProvider == "" {
		d.IssueProvider = "linear"
	}
}

func defaultGoVersion() string {
	return stableGoVersion
}

func shouldSkipTemplate(path string, data TemplateData) bool {
	path = filepath.ToSlash(path)
	rel := strings.TrimPrefix(path, templateRoot+"/")
	if !data.Examples && strings.Contains(rel, "/example_") {
		return true
	}
	if data.IssueProvider != "repo" && strings.HasPrefix(rel, "docs/issues") {
		return true
	}
	if strings.HasPrefix(rel, "internal/controller/comm_controller") || strings.HasPrefix(rel, "internal/repository/mysql/my_common") || strings.HasPrefix(rel, "internal/models/do/mysql/common") {
		return true
	}
	optional := []struct {
		enabled bool
		paths   []string
	}{
		{data.Redis, []string{"internal/lib/redis"}},
		{data.Cron, []string{"internal/cron"}},
		{data.Lark, []string{"internal/service/common_srv/lark_service.go.tmpl", "internal/models/dto/lark_dto", "internal/lib/log/lark_logger.go.tmpl"}},
		{data.Prometheus, []string{"internal/service/common_srv/prometheus_service.go.tmpl"}},
		{data.JWT, []string{"utils/jwt.go.tmpl", "utils/jwt_test.go.tmpl"}},
	}
	for _, feature := range optional {
		if !feature.enabled && anyTemplatePrefixMatches(rel, feature.paths) {
			return true
		}
	}
	if !data.MySQL && anyTemplatePrefixMatches(path, mysqlOnlyTemplatePrefixes) {
		return true
	}
	if !data.MongoDB && anyTemplatePrefixMatches(path, mongoOnlyTemplatePrefixes) {
		return true
	}
	return false
}

func anyTemplatePrefixMatches(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if path == prefix {
			return true
		}
		if strings.HasSuffix(prefix, ".tmpl") {
			continue
		}
		if strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func prepareOutputDir(outputDir string) error {
	info, err := os.Lstat(outputDir)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output directory must not be a symlink: %s", outputDir)
		}
		if !info.IsDir() {
			return fmt.Errorf("output path is not a directory: %s", outputDir)
		}

		entries, err := os.ReadDir(outputDir)
		if err != nil {
			return fmt.Errorf("read output directory %s: %w", outputDir, err)
		}
		if hasVisibleEntries(entries) {
			return fmt.Errorf("output directory is not empty: %s", outputDir)
		}

		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat output directory %s: %w", outputDir, err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", outputDir, err)
	}

	return nil
}

func hasVisibleEntries(entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		return true
	}

	return false
}

func validateModulePath(moduleName string) error {
	v := strings.TrimSpace(moduleName)
	if v == "" {
		return fmt.Errorf("module cannot be empty")
	}
	if strings.ContainsAny(v, " \t\r\n") {
		return fmt.Errorf("module cannot contain whitespace: %q", moduleName)
	}
	if strings.Contains(v, `\`) {
		return fmt.Errorf("module cannot contain backslash: %q", moduleName)
	}
	if strings.ContainsAny(v, `"'`) {
		return fmt.Errorf("module cannot contain quotes: %q", moduleName)
	}
	if err := module.CheckPath(v); err != nil {
		return fmt.Errorf("invalid module path %q: %w", moduleName, err)
	}

	return nil
}

func validateBinaryName(binaryName string) error {
	v := strings.TrimSpace(binaryName)
	if v == "" {
		return fmt.Errorf("binary cannot be empty")
	}
	if strings.ContainsAny(v, " \t\r\n") {
		return fmt.Errorf("binary cannot contain whitespace: %q", binaryName)
	}
	if strings.ContainsAny(v, `/\`) {
		return fmt.Errorf("binary cannot contain path separators: %q", binaryName)
	}
	if strings.ContainsAny(v, `"'`) {
		return fmt.Errorf("binary cannot contain quotes: %q", binaryName)
	}
	if v == "." || v == ".." {
		return fmt.Errorf("binary cannot be %q", binaryName)
	}
	if !binaryNamePattern.MatchString(v) {
		return fmt.Errorf(
			"binary contains invalid characters: %q (allowed: letters, digits, dot, underscore, hyphen)",
			binaryName,
		)
	}

	return nil
}

func validateProjectName(projectName string) error {
	v := strings.TrimSpace(projectName)
	if v == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if v == "." || v == ".." {
		return fmt.Errorf("project name cannot be %q", projectName)
	}
	if strings.ContainsAny(v, `/\`) {
		return fmt.Errorf("project name cannot contain path separators: %q", projectName)
	}
	if strings.ContainsAny(v, " \t\r\n") {
		return fmt.Errorf("project name cannot contain whitespace: %q", projectName)
	}
	if strings.ContainsAny(v, `"'`) {
		return fmt.Errorf("project name cannot contain quotes: %q", projectName)
	}
	if !projectNamePattern.MatchString(v) {
		return fmt.Errorf(
			"project name contains invalid characters: %q (allowed: letters, digits, dot, underscore, hyphen)",
			projectName,
		)
	}

	return nil
}

func validateGoVersion(goVersion string) error {
	v := strings.TrimSpace(goVersion)
	if v == "" {
		return fmt.Errorf("go version cannot be empty")
	}
	if !goVersionPattern.MatchString(v) {
		return fmt.Errorf(
			"go version contains invalid format: %q (example: 1.27.1)",
			goVersion,
		)
	}
	if goversion.Compare("go"+v, "go"+minimumTemplateGoVersion) < 0 {
		return fmt.Errorf("go version %q is below the template minimum %s", goVersion, minimumTemplateGoVersion)
	}

	return nil
}
