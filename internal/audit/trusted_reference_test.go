package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTrustedReferenceIsRestrictedToConfiguredFileAndRoot(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, "config"))
	root := filepath.Join(dir, "repo")
	skill := filepath.Join(root, "skills", "swarm")
	target := filepath.Join(root, "docs", "adapter.md")
	if err := os.MkdirAll(skill, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("adapter"), 0644); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(root, "docs", "other.md")
	if err := os.WriteFile(other, []byte("other"), 0644); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "config", "skillshare")
	if err := os.MkdirAll(config, 0700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal([]trustedReferenceEntry{{Root: root, Target: "docs/adapter.md"}})
	if err := os.WriteFile(filepath.Join(config, "trusted-references.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if !trustedReference(skill, "../../docs/adapter.md") {
		t.Fatal("configured adapter denied")
	}
	if trustedReference(skill, "../../docs/other.md") {
		t.Fatal("unconfigured file allowed")
	}
	outside := filepath.Join(dir, "outside.md")
	if err := os.WriteFile(outside, []byte("private"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "docs", "link.md")); err != nil {
		t.Fatal(err)
	}
	if trustedReference(skill, "../../docs/link.md") {
		t.Fatal("symlink escape allowed")
	}
	if trustedReference(dir, "repo/docs/adapter.md") {
		t.Fatal("untrusted skill allowed")
	}
}

func TestTrustedTemplateSuppressesOnlyItsLiteralExample(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, "config"))
	root := filepath.Join(dir, "repo")
	skill := filepath.Join(root, "skills", "why")
	if err := os.MkdirAll(filepath.Join(skill, "references"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skill, "SKILL.md"), []byte("# Fixture\n"), 0644); err != nil {
		t.Fatal(err)
	}
	content := []byte("Source: [PR #123](url)\n[real missing](missing.md)\n")
	if err := os.WriteFile(filepath.Join(skill, "references", "prompt.md"), content, 0644); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "config", "skillshare")
	if err := os.MkdirAll(config, 0700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal([]trustedReferenceEntry{{Root: root, Target: "skills/why/references/prompt.md"}})
	if err := os.WriteFile(filepath.Join(config, "trusted-references.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	scan, err := ScanSkillWithRules(skill, nil)
	if err != nil {
		t.Fatal(err)
	}
	var placeholder, missing bool
	for _, f := range scan.Findings {
		if f.RuleID == "dangling-link" && strings.Contains(f.Message, `"url"`) {
			placeholder = true
		}
		if f.RuleID == "dangling-link" && strings.Contains(f.Message, "missing.md") {
			missing = true
		}
	}
	if placeholder || !missing {
		t.Fatalf("placeholder=%v missing=%v findings=%+v", placeholder, missing, scan.Findings)
	}
}
