package scaf_fold

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func checksum(body []byte) string { sum := sha256.Sum256(body); return hex.EncodeToString(sum[:]) }
func readManifest(root *os.Root) (Manifest, error) {
	if err := safeManagedPath(root, ".starter.json"); err != nil {
		return Manifest{}, err
	}
	body, err := root.ReadFile(".starter.json")
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return m, err
	}
	if m.Schema != 1 {
		return m, fmt.Errorf("unsupported starter manifest schema %d", m.Schema)
	}
	return m, nil
}
func manifestBytes(m Manifest) ([]byte, error) {
	body, err := json.Marshal(m, json.Deterministic(true), jsontext.WithIndent("  "))
	return append(body, '\n'), err
}

// RefreshManifest 仅用于维护时依赖解析之后，记录 go.mod 和 go.sum 的最终内容。
func RefreshManifest(dir string) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	m, err := readManifest(root)
	if err != nil {
		return err
	}
	for _, name := range []string{"go.mod", "go.sum"} {
		body, err := root.ReadFile(name)
		if err != nil {
			return err
		}
		m.Files[name] = checksum(body)
	}
	body, err := manifestBytes(m)
	if err != nil {
		return err
	}
	return root.WriteFile(".starter.json", body, 0o644)
}

// SyncReference 更新固定参考工程。任何已管理文件的本地修改都会阻止整次更新。
// 普通业务项目应使用 new/init；此接口不接管未知文件。
func SyncReference(source, target string, write bool) ([]string, error) {
	info, err := os.Lstat(target)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("reference target must be a real directory")
	}
	src, err := os.OpenRoot(source)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	dst, err := os.OpenRoot(target)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	desired, err := readManifest(src)
	if err != nil {
		return nil, err
	}
	// 仓库法律信息与独立发布历史由各仓维护。
	delete(desired.Files, "LICENSE")
	delete(desired.Files, "changeLog.md")
	previous, err := readManifest(dst)
	if errors.Is(err, os.ErrNotExist) {
		previous = Manifest{Schema: 1, Files: map[string]string{}}
	} else if err != nil {
		return nil, err
	}
	delete(previous.Files, "LICENSE")
	delete(previous.Files, "changeLog.md")
	paths := maps.Clone(previous.Files)
	maps.Copy(paths, desired.Files)
	changes := []string{}
	payloads := map[string][]byte{}
	for _, path := range slices.Sorted(maps.Keys(paths)) {
		if err := safeManagedPath(dst, path); err != nil {
			return nil, err
		}
		current, readErr := dst.ReadFile(path)
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			return nil, readErr
		}
		old, managed := previous.Files[path]
		if managed && errors.Is(readErr, os.ErrNotExist) {
			return nil, fmt.Errorf("preserve locally deleted file %s; merge it before syncing", path)
		}
		if readErr == nil && (!managed || checksum(current) != old) {
			return nil, fmt.Errorf("preserve local or unowned file %s; merge it before syncing", path)
		}
		want, exists := desired.Files[path]
		if !exists {
			if readErr == nil {
				changes = append(changes, "delete "+path)
			}
			continue
		}
		body, err := src.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if checksum(body) != want {
			return nil, fmt.Errorf("source checksum mismatch: %s", path)
		}
		if readErr != nil || checksum(current) != want {
			changes = append(changes, "write "+path)
			payloads[path] = body
		}
	}
	body, err := manifestBytes(desired)
	if err != nil {
		return nil, err
	}
	oldBody, err := dst.ReadFile(".starter.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if string(body) != string(oldBody) {
		changes = append(changes, "write .starter.json")
		payloads[".starter.json"] = body
	}
	if !write {
		return changes, nil
	}
	for _, change := range changes {
		action, path, _ := strings.Cut(change, " ")
		if err := safeManagedPath(dst, path); err != nil {
			return nil, err
		}
		if action == "delete" {
			if err := dst.Remove(path); err != nil {
				return nil, err
			}
			continue
		}
		if err := dst.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		if err := dst.WriteFile(path, payloads[path], 0o644); err != nil {
			return nil, err
		}
	}
	return changes, nil
}

func safeManagedPath(root *os.Root, path string) error {
	if !fs.ValidPath(path) || path == ".git" || strings.HasPrefix(path, ".git/") {
		return fmt.Errorf("invalid managed path %q", path)
	}
	current := ""
	for part := range strings.SplitSeq(path, "/") {
		current = filepath.Join(current, part)
		info, err := root.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing managed symlink %s", current)
		}
	}
	return nil
}
