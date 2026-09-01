package run

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pickupskill "kck/skills/kck-pickup-a-session"
)

const pickupSkillCacheRel = ".cache/kck-pickup-a-session/SKILL.md"

// ensurePickupSkillCached writes the embedded kck-pickup-a-session skill to
// ~/.cache/kck-pickup-a-session/SKILL.md when missing or MD5-mismatched.
// Returns the absolute cache path.
func ensurePickupSkillCached(popts *PickupOpts) (string, error) {
	content := []byte(pickupSkillContent(popts))
	wantSum := md5Hex(content)

	path, err := pickupSkillCachePath(popts)
	if err != nil {
		return "", err
	}

	readFile := popts.ReadFile
	if readFile == nil {
		readFile = os.ReadFile
	}
	if existing, rerr := readFile(path); rerr == nil {
		if md5Hex(existing) == wantSum {
			return path, nil
		}
	}

	mkdirAll := popts.MkdirAll
	if mkdirAll == nil {
		mkdirAll = os.MkdirAll
	}
	dir := filepath.Dir(path)
	if err := mkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create skill cache dir: %w", err)
	}

	writeFile := popts.WriteFile
	if writeFile == nil {
		writeFile = os.WriteFile
	}
	rename := popts.Rename
	if rename == nil {
		rename = os.Rename
	}

	tmp := path + ".tmp"
	if err := writeFile(tmp, content, 0o644); err != nil {
		return "", fmt.Errorf("write skill cache: %w", err)
	}
	if err := rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("finalize skill cache: %w", err)
	}
	return path, nil
}

func pickupSkillContent(popts *PickupOpts) string {
	if popts != nil && popts.SkillContent != "" {
		return popts.SkillContent
	}
	return pickupskill.SkillMD
}

func pickupSkillCachePath(popts *PickupOpts) (string, error) {
	if popts != nil {
		if p := strings.TrimSpace(popts.CacheSkillPath); p != "" {
			abs, err := filepath.Abs(p)
			if err != nil {
				return "", fmt.Errorf("skill cache path: %w", err)
			}
			return abs, nil
		}
	}
	homeFn := os.UserHomeDir
	if popts != nil && popts.UserHomeDir != nil {
		homeFn = popts.UserHomeDir
	}
	home, err := homeFn()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", fmt.Errorf("cannot resolve home for %s skill cache", pickupskill.Name)
	}
	return filepath.Join(home, pickupSkillCacheRel), nil
}

func md5Hex(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}
