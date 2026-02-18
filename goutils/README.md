# goutils — 50 Unix CLI Utilities in Go

A collection of 50 standard Unix utilities reimplemented in pure Go with no external dependencies.

## Build All

```bash
# Build all binaries into ./bin/
mkdir -p bin
for d in cmd/*/; do
  name=$(basename $d)
  go build -o bin/$name ./$d
done
```

## Utilities

### File & Directory Operations

| Utility | Usage | Description |
|---------|-------|-------------|
| `cp` | `cp [-r] <src> <dst>` | Copy files/directories |
| `mv` | `mv <src> <dst>` | Move or rename files |
| `rm` | `rm [-r] [-f] <file>...` | Remove files/directories |
| `mkdir` | `mkdir [-p] <dir>...` | Create directories |
| `find` | `find [path] [-name pat] [-type f\|d] [-size n]` | Search for files |
| `tree` | `tree [dir] [-a] [-L depth]` | Display directory as a tree |
| `du` | `du [-h] [-s] [path...]` | Disk usage |
| `touch` | `touch <file>...` | Create files or update timestamps |

### Text Processing

| Utility | Usage | Description |
|---------|-------|-------------|
| `cat` | `cat [-n] [file...]` | Concatenate and print files |
| `grep` | `grep [-i] [-n] [-v] [-r] [-c] <pattern> [file...]` | Search with regex |
| `wc` | `wc [-l] [-w] [-c] [file...]` | Count lines, words, chars |
| `head` | `head [-n N] [file...]` | Print first N lines |
| `tail` | `tail [-n N] [-f] [file]` | Print last N lines, follow |
| `sed` | `sed [-i] 's/pat/repl/[g]' [file...]` | Find and replace |
| `sort` | `sort [-r] [-n] [-u] [-k field] [file...]` | Sort lines |
| `uniq` | `uniq [-c] [-d] [-u] [file]` | Filter duplicate lines |
| `cut` | `cut -f fields [-d delim] [file...]` | Extract fields/chars |

### System & Process

| Utility | Usage | Description |
|---------|-------|-------------|
| `ps` | `ps [-a]` | List running processes |
| `kill` | `kill [-s signal] <pid>...` | Send signals to processes |
| `env` | `env [NAME=VAL...] [cmd]` | Print/set environment variables |
| `which` | `which [-a] <cmd>...` | Locate command in PATH |
| `uptime` | `uptime` | Show system uptime & load |

### Networking

| Utility | Usage | Description |
|---------|-------|-------------|
| `ping` | `ping [-c count] [-i interval] <host>` | ICMP ping (needs root/CAP_NET_RAW) |
| `curl` | `curl [-X method] [-H hdr] [-d data] [-o file] [-i] <url>` | HTTP requests |
| `wget` | `wget [-O file] [-q] <url>` | Download files |
| `netstat` | `netstat [-l] [-t] [-u]` | Show open connections (Linux) |
| `dns` | `dns [-type A\|MX\|NS\|TXT\|CNAME] <host>` | DNS lookup |

### Output & I/O

| Utility | Usage | Description |
|---------|-------|-------------|
| `echo` | `echo [-n] [-e] [string...]` | Print arguments; `-e` enables `\n`, `\t` escapes |
| `tee` | `tee [-a] [file...]` | Read stdin, write to stdout + files simultaneously |
| `yes` | `yes [string]` | Repeatedly output a string until killed |

### Path & File Info

| Utility | Usage | Description |
|---------|-------|-------------|
| `basename` | `basename <path> [suffix]` | Strip directory from filename |
| `dirname` | `dirname <path>...` | Strip filename from path |
| `pwd` | `pwd [-P]` | Print working directory (`-P` resolves symlinks) |
| `ln` | `ln [-s] [-f] <target> <link>` | Create hard or symbolic links |
| `stat` | `stat <file>...` | Display detailed file metadata |

### File Comparison & Checksums

| Utility | Usage | Description |
|---------|-------|-------------|
| `diff` | `diff [-u] [-i] <file1> <file2>` | Compare files line by line |
| `md5sum` | `md5sum [-c] [file...]` | Compute or verify MD5 checksums |
| `sha256sum` | `sha256sum [-c] [file...]` | Compute or verify SHA-256 checksums |

### Binary & Encoding

| Utility | Usage | Description |
|---------|-------|-------------|
| `xxd` | `xxd [-c cols] [-l limit] [-r] [file]` | Hex dump; `-r` reverses back to binary |
| `base64` | `base64 [-d] [-w cols] [file]` | Encode or decode base64 |

### Text Transformation

| Utility | Usage | Description |
|---------|-------|-------------|
| `tr` | `tr [-d] [-s] <set1> [set2]` | Translate or delete characters |
| `fmt` | `fmt [-w width] [file...]` | Word-wrap text to a given width |
| `nl` | `nl [-b a\|t\|n] [-n ln\|rn\|rz] [-w N] [file...]` | Number lines |
| `tac` | `tac [file...]` | Print file lines in reverse order |
| `rev` | `rev [file...]` | Reverse characters on each line |

### Process Control

| Utility | Usage | Description |
|---------|-------|-------------|
| `timeout` | `timeout [-s signal] <duration> <cmd> [args]` | Run command with time limit |
| `xargs` | `xargs [-n N] [-I str] [-P N] [-0] [-t] <cmd>` | Build commands from stdin |

### Misc

| Utility | Usage | Description |
|---------|-------|-------------|
| `cal` | `cal [month] [year]` | Print a calendar |
| `watcher` | `watcher [-r] <dir>` | Watch directory for changes |

## Notes

- `ping` requires root or `CAP_NET_RAW` capability on Linux
- `netstat` and `uptime` read from `/proc` — Linux only
- `ps` reads from `/proc` — Linux only
- All other utilities are cross-platform (Linux, macOS, Windows)

## Examples

```bash
# Find all .go files larger than 1KB
./bin/find . -name "*.go" -size 1024

# Show directory tree 2 levels deep
./bin/tree . -L 2

# Count lines in all Go files
./bin/grep -r "" --count . 

# Replace "foo" with "bar" in place
./bin/sed -i 's/foo/bar/g' file.txt

# Watch current directory for changes
./bin/watcher -r .

# Download a file quietly
./bin/wget -q -O output.html https://example.com

# DNS lookup for MX records
./bin/dns -type MX gmail.com

# Show calendar for March 2025
./bin/cal 3 2025
```

```bash
# Word-wrap a long text file at 80 columns
./bin/fmt -w 80 essay.txt

# Hex dump first 64 bytes of a binary
./bin/xxd -l 64 binary.bin

# Reverse a hex dump back to binary
./bin/xxd binary.bin | ./bin/xxd -r > copy.bin

# SHA-256 checksum of multiple files
./bin/sha256sum file1 file2 > checksums.txt
./bin/sha256sum -c checksums.txt

# Run diff in unified format
./bin/diff -u old.txt new.txt

# Convert uppercase to lowercase
echo "HELLO WORLD" | ./bin/tr '[:upper:]' '[:lower:]'

# Delete all digits from input
echo "abc123def" | ./bin/tr -d '[:digit:]'

# Number all lines in a file
./bin/nl -b a file.txt

# Base64 encode a file
./bin/base64 image.png > image.b64
./bin/base64 -d image.b64 > restored.png

# Run find in parallel with xargs
./bin/find . -name "*.log" | ./bin/xargs -P 4 -n 1 gzip

# Replace {} placeholder
echo "/tmp/file.txt" | ./bin/xargs -I{} cp {} {}.bak

# Timeout a slow command after 5 seconds
./bin/timeout 5s sleep 100

# Write to stdout and a log file simultaneously
./bin/yes "test line" | head -5 | ./bin/tee output.log

# Reverse lines in a file (like tac)
./bin/tac logfile.txt

# Reverse characters on each line
echo "hello" | ./bin/rev   # → olleh
```

## Notes

- `stat` uses Linux-specific syscall fields (Inode, UID, GID). On macOS/Windows, it gracefully skips those.
- `timeout` exit code 124 means the command was killed due to timeout (POSIX convention).
- `xargs -P` runs commands in parallel using goroutines.
- `diff` uses a simple O(n×m) LCS algorithm; works well for small-to-medium files.
- `tr` supports character ranges (`a-z`), POSIX classes (`[:upper:]`, `[:lower:]`, `[:digit:]`, `[:alpha:]`, `[:space:]`).
