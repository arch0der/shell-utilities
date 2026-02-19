// gitignore - generate .gitignore files for common project types
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

var templates = map[string]string{
	"go": `# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.sum
vendor/
bin/
dist/
`,
	"python": `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
*.egg-info/
.installed.cfg
*.egg
.env
.venv
env/
venv/
ENV/
.pytest_cache/
.mypy_cache/
.ruff_cache/
htmlcov/
.coverage
*.cover
`,
	"node": `# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.npm
.yarn/cache
.yarn/unplugged
.pnp.*
dist/
build/
.cache/
*.tsbuildinfo
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
`,
	"rust": `# Rust
/target/
Cargo.lock
**/*.rs.bk
*.pdb
`,
	"java": `# Java
*.class
*.log
*.ctxt
.mtj.tmp/
*.jar
*.war
*.nar
*.ear
*.zip
*.tar.gz
*.rar
hs_err_pid*
.gradle/
build/
!gradle/wrapper/gradle-wrapper.jar
!**/src/main/**/build/
!**/src/test/**/build/
`,
	"macos": `# macOS
.DS_Store
.AppleDouble
.LSOverride
Icon
._*
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk
`,
	"windows": `# Windows
Thumbs.db
Thumbs.db:encryptable
ehthumbs.db
ehthumbs_vista.db
*.stackdump
[Dd]esktop.ini
$RECYCLE.BIN/
*.cab
*.msi
*.msix
*.msm
*.msp
*.lnk
`,
	"linux": `# Linux
*~
.fuse_hidden*
.directory
.Trash-*
.nfs*
`,
	"vim": `# Vim
[._]*.s[a-v][a-z]
!*.svg
[._]*.sw[a-p]
[._]s[a-rt-v][a-z]
[._]ss[a-gi-z]
[._]sw[a-p]
Session.vim
Sessionx.vim
.netrwhist
*~
tags
[._]*.un~
`,
	"vscode": `# VS Code
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
!.vscode/*.code-snippets
.history/
*.vsix
`,
	"terraform": `# Terraform
.terraform/
.terraform.lock.hcl
*.tfstate
*.tfstate.*
crash.log
crash.*.log
*.tfvars
*.tfvars.json
override.tf
override.tf.json
*_override.tf
*_override.tf.json
.terraformrc
terraform.rc
`,
	"docker": `# Docker
.docker/
docker-compose.override.yml
.dockerignore
`,
	"general": `# General
*.log
*.tmp
*.temp
*.bak
*.backup
*.swp
*.swo
*~
.env
.env.*
!.env.example
secrets.*
*.pem
*.key
*.cert
*.p12
.idea/
.vscode/
*.iml
`,
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-l" || os.Args[1] == "--list" {
		keys := make([]string, 0, len(templates))
		for k := range templates { keys = append(keys, k) }
		sort.Strings(keys)
		fmt.Println("Available templates:")
		for _, k := range keys { fmt.Printf("  %s\n", k) }
		return
	}

	var parts []string
	write := false
	for _, arg := range os.Args[1:] {
		if arg == "-w" || arg == "--write" { write = true; continue }
		t, ok := templates[arg]
		if !ok { fmt.Fprintf(os.Stderr, "gitignore: unknown template %q (use -l to list)\n", arg); os.Exit(1) }
		parts = append(parts, fmt.Sprintf("# === %s ===\n%s", strings.ToUpper(arg), t))
	}

	content := fmt.Sprintf("# .gitignore — generated %s\n\n%s", time.Now().Format("2006-01-02"), strings.Join(parts, "\n"))
	if write {
		if err := os.WriteFile(".gitignore", []byte(content), 0644); err != nil {
			fmt.Fprintln(os.Stderr, err); os.Exit(1)
		}
		fmt.Println("Written: .gitignore")
	} else {
		fmt.Print(content)
	}
}
