// genbash - generate boilerplate bash scripts with best practices
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const scriptTemplate = `#!/usr/bin/env bash
# {{NAME}} - {{DESC}}
# Generated: {{DATE}}
# Usage: {{NAME}} [options] [args]

set -euo pipefail
IFS=$'\n\t'

# ─── Constants ────────────────────────────────────────────────────────────────
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="$(basename "$0")"
readonly VERSION="0.1.0"

# ─── Defaults ─────────────────────────────────────────────────────────────────
VERBOSE=false
DRY_RUN=false
LOG_FILE=""
{{EXTRA_VARS}}
# ─── Logging ──────────────────────────────────────────────────────────────────
log()     { echo "[$(date '+%H:%M:%S')] $*" >&2; }
info()    { echo "[INFO]  $*" >&2; }
warn()    { echo "[WARN]  $*" >&2; }
error()   { echo "[ERROR] $*" >&2; exit 1; }
debug()   { [[ "$VERBOSE" == true ]] && echo "[DEBUG] $*" >&2 || true; }
success() { echo "[OK]    $*" >&2; }

# ─── Helpers ──────────────────────────────────────────────────────────────────
usage() {
  cat <<USAGE
Usage: $SCRIPT_NAME [options] [args]

Options:
  -h, --help        Show this help message
  -v, --verbose     Enable verbose output
  -n, --dry-run     Show actions without executing
  -V, --version     Print version
{{EXTRA_OPTS_HELP}}
Examples:
  $SCRIPT_NAME --verbose
  $SCRIPT_NAME --dry-run input.txt

USAGE
}

require_cmd() { command -v "$1" &>/dev/null || error "Required command not found: $1"; }
require_file() { [[ -f "$1" ]] || error "File not found: $1"; }
require_dir()  { [[ -d "$1" ]] || error "Directory not found: $1"; }

cleanup() {
  local exit_code=$?
  debug "Cleaning up (exit code: $exit_code)"
  # Add cleanup here
  exit "$exit_code"
}
trap cleanup EXIT INT TERM

# ─── Argument Parsing ─────────────────────────────────────────────────────────
parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)    usage; exit 0 ;;
      -v|--verbose) VERBOSE=true ;;
      -n|--dry-run) DRY_RUN=true ;;
      -V|--version) echo "$SCRIPT_NAME v$VERSION"; exit 0 ;;
{{EXTRA_OPTS_PARSE}}      --) shift; break ;;
      -*) error "Unknown option: $1" ;;
      *)  POSITIONAL+=("$1") ;;
    esac
    shift
  done
  POSITIONAL=("${POSITIONAL[@]+"${POSITIONAL[@]}"}")
}

POSITIONAL=()
parse_args "$@"
set -- "${POSITIONAL[@]+"${POSITIONAL[@]}"}"

# ─── Main ─────────────────────────────────────────────────────────────────────
main() {
  debug "Starting $SCRIPT_NAME v$VERSION"
  [[ "$DRY_RUN" == true ]] && info "Dry-run mode enabled"

  # TODO: implement main logic here
  info "Done."
}

main "$@"
`

func main() {
	name := "myscript"
	desc := "does something useful"
	if len(os.Args) > 1 { name = os.Args[1] }
	if len(os.Args) > 2 { desc = strings.Join(os.Args[2:], " ") }

	out := scriptTemplate
	out = strings.ReplaceAll(out, "{{NAME}}", name)
	out = strings.ReplaceAll(out, "{{DESC}}", desc)
	out = strings.ReplaceAll(out, "{{DATE}}", time.Now().Format("2006-01-02"))
	out = strings.ReplaceAll(out, "{{EXTRA_VARS}}", "")
	out = strings.ReplaceAll(out, "{{EXTRA_OPTS_HELP}}", "")
	out = strings.ReplaceAll(out, "{{EXTRA_OPTS_PARSE}}", "")

	fname := name + ".sh"
	if err := os.WriteFile(fname, []byte(out), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err); os.Exit(1)
	}
	fmt.Printf("Generated: %s\n", fname)
}
