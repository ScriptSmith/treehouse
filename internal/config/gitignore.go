package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kunchenguid/treehouse/v3/internal/git"
)

// EnsureGitignore adds treehouseDir to the .gitignore of the enclosing git
// repo, if treehouseDir is inside a git repo. It is a no-op if the directory
// is not inside a repo or if git already ignores it through any ignore source
// (.gitignore, .git/info/exclude, or the user's core.excludesFile).
func EnsureGitignore(treehouseDir string) error {
	// Walk up from treehouseDir to find an existing ancestor for the git check,
	// since the directory itself may not exist yet.
	checkDir := treehouseDir
	for {
		if info, err := os.Stat(checkDir); err == nil && info.IsDir() {
			break
		}
		parent := filepath.Dir(checkDir)
		if parent == checkDir {
			return nil
		}
		checkDir = parent
	}

	repoRoot, err := git.FindRepoRootFrom(checkDir)
	if err != nil {
		// Not inside a git repo — nothing to do.
		return nil
	}

	rel, err := filepath.Rel(repoRoot, treehouseDir)
	if err != nil {
		return nil
	}

	// Use forward slashes for .gitignore and prefix with /
	entry := "/" + filepath.ToSlash(rel)

	// Skip the append when git already ignores the directory, e.g. via a
	// global excludes file, .git/info/exclude, or a broader .gitignore
	// pattern. check-ignore only matches directory patterns like
	// ".treehouse/" against existing paths, so probe a file inside the
	// directory instead of the directory itself.
	if ignored, err := git.IsIgnored(repoRoot, filepath.Join(rel, ".treehouse-ignore-probe")); err == nil && ignored {
		return nil
	}

	gitignorePath := filepath.Join(repoRoot, ".gitignore")
	existing, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}

	f, err := os.OpenFile(gitignorePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	prefix := ""
	if len(existing) > 0 && !strings.HasSuffix(string(existing), "\n") {
		prefix = "\n"
	}
	_, err = f.WriteString(prefix + entry + "\n")
	return err
}
