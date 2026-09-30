package scaf_fold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReferenceSyncPreservesLocalEditsAndUnmanagedFiles(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	target := filepath.Join(t.TempDir(), "target")
	data := TemplateData{ModuleName: "example.com/reference", BinaryName: "reference", ProjectName: "reference"}
	if err := Generate(source, data); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "notes.txt"), []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("local edits"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err == nil {
		t.Fatal("overwrote a modified managed file")
	}
	for name, want := range map[string]string{"notes.txt": "keep me", "README.md": "local edits"} {
		got, err := os.ReadFile(filepath.Join(target, name))
		if err != nil || string(got) != want {
			t.Fatalf("%s changed: %s %v", name, got, err)
		}
	}
}

func TestReferenceCheckDoesNotWrite(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	target := filepath.Join(t.TempDir(), "target")
	if err := Generate(source, TemplateData{ModuleName: "example.com/check", BinaryName: "check", ProjectName: "check"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	changes, err := SyncReference(source, target, false)
	if err != nil || len(changes) == 0 {
		t.Fatalf("check=%v %v", changes, err)
	}
	entries, err := os.ReadDir(target)
	if err != nil || len(entries) != 0 {
		t.Fatal("check wrote files")
	}
}

func TestReferenceSyncRejectsSymlinkManifest(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	target := t.TempDir()
	if err := Generate(source, TemplateData{ModuleName: "example.com/sync", BinaryName: "sync", ProjectName: "sync"}); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(target, ".starter.json")
	if err := os.Rename(manifest, filepath.Join(target, "saved.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("saved.json", manifest); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err == nil {
		t.Fatal("accepted symlink manifest")
	}
}

func TestReferenceSyncRejectsLocalDeletion(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	target := t.TempDir()
	if err := Generate(source, TemplateData{ModuleName: "example.com/sync", BinaryName: "sync", ProjectName: "sync"}); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(target, "README.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncReference(source, target, true); err == nil {
		t.Fatal("overwrote a local deletion")
	}
}
