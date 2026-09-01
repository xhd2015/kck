package run

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsurePickupSkillCached_WritesThenSkips(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "SKILL.md")
	content := "# hello cache\n"
	var writes int
	popts := &PickupOpts{
		SkillContent:   content,
		CacheSkillPath: cachePath,
		WriteFile: func(path string, data []byte, perm os.FileMode) error {
			writes++
			return os.WriteFile(path, data, perm)
		},
		MkdirAll:  os.MkdirAll,
		Rename:    os.Rename,
		ReadFile:  os.ReadFile,
		UserHomeDir: func() (string, error) { return dir, nil },
	}

	path1, err := ensurePickupSkillCached(popts)
	if err != nil {
		t.Fatalf("first ensure: %v", err)
	}
	if path1 != cachePath {
		t.Fatalf("path=%q want %q", path1, cachePath)
	}
	if writes != 1 {
		t.Fatalf("writes=%d want 1", writes)
	}
	got, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != content {
		t.Fatalf("content=%q", got)
	}
	sum := md5.Sum([]byte(content))
	if hex.EncodeToString(sum[:]) != md5Hex(got) {
		t.Fatalf("md5 mismatch")
	}

	path2, err := ensurePickupSkillCached(popts)
	if err != nil {
		t.Fatalf("second ensure: %v", err)
	}
	if path2 != cachePath {
		t.Fatalf("path=%q", path2)
	}
	if writes != 1 {
		t.Fatalf("second ensure must skip write; writes=%d", writes)
	}
}

func TestEnsurePickupSkillCached_RewritesOnMismatch(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(cachePath, []byte("stale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var writes int
	popts := &PickupOpts{
		SkillContent:   "fresh\n",
		CacheSkillPath: cachePath,
		WriteFile: func(path string, data []byte, perm os.FileMode) error {
			writes++
			return os.WriteFile(path, data, perm)
		},
		MkdirAll: os.MkdirAll,
		Rename:   os.Rename,
		ReadFile: os.ReadFile,
	}
	if _, err := ensurePickupSkillCached(popts); err != nil {
		t.Fatal(err)
	}
	if writes != 1 {
		t.Fatalf("writes=%d want 1", writes)
	}
	got, _ := os.ReadFile(cachePath)
	if string(got) != "fresh\n" {
		t.Fatalf("got %q", got)
	}
}
