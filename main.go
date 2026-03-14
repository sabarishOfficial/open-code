package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const DefaultWorkDir = "~/Desktop/GitHub"

// Editors shown in the fzf picker
var Editors = []string{
	"cursor",
	"code",
	"vim",
	"nvim",
	"nano",
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	workDir := resolveWorkDir()
	fmt.Printf("📁 Scanning: %s\n\n", workDir)

	// 1. List folders
	folders, err := listFolders(workDir)
	if err != nil {
		fatalf("Cannot read directory %q: %v", workDir, err)
	}
	if len(folders) == 0 {
		fatalf("No folders found in %s", workDir)
	}

	// 2. Pick folder via fzf
	folder, err := fzfPick(folders, "📂 Project > ")
	if err != nil {
		fmt.Println("No folder selected. Bye!")
		os.Exit(0)
	}

	projectPath := filepath.Join(workDir, folder)
	fmt.Printf("✅ Project : %s\n\n", projectPath)

	// 3. Pick editor via fzf
	editor, err := fzfPick(Editors, "🖊️  Editor  > ")
	if err != nil {
		fmt.Println("No editor selected. Bye!")
		os.Exit(0)
	}

	// 4. Open
	fmt.Printf("🚀 Opening with %s...\n", editor)
	if err := openWith(editor, projectPath); err != nil {
		fatalf("Could not launch %q: %v\n\nMake sure %q is installed and in your PATH.", editor, err, editor)
	}
}

func resolveWorkDir() string {
	if len(os.Args) > 1 {
		return expandHome(os.Args[1])
	}
	if env := os.Getenv("DEVOPEN_DIR"); env != "" {
		return expandHome(env)
	}
	return expandHome(DefaultWorkDir)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func listFolders(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func fzfPick(items []string, prompt string) (string, error) {
	// Fallback if fzf is not installed
	if _, err := exec.LookPath("fzf"); err != nil {
		return fallbackSelect(items, prompt)
	}

	input := strings.Join(items, "\n")

	cmd := exec.Command("fzf",
		"--prompt", prompt,
		"--height", "50%",
		"--layout", "reverse",
		"--border", "rounded",
		"--margin", "1,2",
		"--info", "inline",
		"--header", "Use ↑↓ arrows · Enter to select · Esc to quit",
	)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	selected := strings.TrimSpace(buf.String())
	if selected == "" {
		return "", fmt.Errorf("empty selection")
	}
	return selected, nil
}

func fallbackSelect(items []string, prompt string) (string, error) {
	fmt.Println("⚠️  fzf not found — using plain list. Install fzf for TUI mode.")
	fmt.Println(prompt)
	for i, item := range items {
		fmt.Printf("  [%d] %s\n", i+1, item)
	}
	fmt.Print("\nEnter number: ")
	var n int
	if _, err := fmt.Scan(&n); err != nil || n < 1 || n > len(items) {
		return "", fmt.Errorf("invalid selection")
	}
	return items[n-1], nil
}

func openWith(editor, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}