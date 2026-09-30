package scaf_fold

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io/fs"

	"github.com/SisyphusSQ/go-web-starter/v2/vars"
)

const HarnessRevision = "b20e5e8ece6a529c7d74aa0a8b1de77bd06c374c"

// Manifest 标识模板来源和生成选项；不作为覆盖用户文件的授权。
type Manifest struct {
	GeneratorVersion string            `json:"generator_version"`
	Schema           int               `json:"schema"`
	TemplateDigest   string            `json:"template_digest"`
	HarnessRevision  string            `json:"harness_revision"`
	Options          TemplateData      `json:"options"`
	Files            map[string]string `json:"files"`
}

func buildManifest(data TemplateData, files map[string][]byte) ([]byte, error) {
	hash := sha256.New()
	if err := fs.WalkDir(templateFS, templateRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		body, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintf(hash, "%s\x00", path)
		hash.Write(body)
		return nil
	}); err != nil {
		return nil, err
	}
	m := Manifest{GeneratorVersion: vars.ReleaseVersion, Schema: 1, TemplateDigest: hex.EncodeToString(hash.Sum(nil)), HarnessRevision: HarnessRevision, Options: data, Files: map[string]string{}}
	for path, body := range files {
		sum := sha256.Sum256(body)
		m.Files[path] = hex.EncodeToString(sum[:])
	}
	body, err := json.Marshal(m, json.Deterministic(true), jsontext.WithIndent("  "))
	if err != nil {
		return nil, err
	}
	return append(body, '\n'), nil
}
