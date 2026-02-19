# Go Coreutils

A Go implementation of GNU coreutils as a multi-call binary (similar to BusyBox).

## Building

```bash
go build -o coreutils .
```

## Usage

### As multi-call binary
```bash
./coreutils <command> [args...]
```

### Via symlinks (recommended)
```bash
# Create symlinks for all tools
./coreutils --list  # to see available commands

# Example: create individual symlinks
ln -s coreutils echo
ln -s coreutils ls
ln -s coreutils cat
# etc.
```

### Example build script
```bash
#!/bin/bash
go build -o coreutils .
for cmd in $(./coreutils 2>&1 | grep "  " | awk '{print $1}'); do
    ln -sf coreutils "$cmd"
done
```

## Implemented Utilities

| Utility | Description |
|---------|-------------|
| arch | Print machine hardware name |
| b2sum | Compute BLAKE2 checksums (uses SHA-512 without external deps) |
| base32 | Base32 encode/decode |
| base64 | Base64 encode/decode |
| basename | Strip directory and suffix from filenames |
| basenc | Encode/decode with various encodings |
| cat | Concatenate and print files |
| chcon | Change SELinux security context (stub) |
| chgrp | Change group ownership |
| chmod | Change file permissions |
| chown | Change file owner and group |
| chroot | Run command with different root directory |
| cksum | Print CRC checksum and byte counts |
| comm | Compare two sorted files line by line |
| cp | Copy files and directories |
| csplit | Split file into sections determined by patterns |
| cut | Remove sections from lines of files |
| date | Print or set system date and time |
| dd | Convert and copy a file |
| df | Report disk space usage |
| dir | List directory contents (like ls -C -b) |
| dircolors | Color setup for ls |
| dirname | Strip last component from filename |
| du | Estimate file space usage |
| echo | Display a line of text |
| env | Run a program in a modified environment |
| expr | Evaluate expressions |
| factor | Print prime factors |
| false | Do nothing, unsuccessfully |
| fmt | Simple optimal text formatter |
| fold | Wrap each input line to fit in specified width |
| groups | Print group memberships |
| head | Output the first part of files |
| hostid | Print the numeric identifier for the current host |
| hostname | Show or set system hostname |
| id | Print real and effective user and group IDs |
| install | Copy files and set attributes |
| join | Join lines of two files on a common field |
| kill | Send signals to processes |
| link | Call the link function to create a link |
| ln | Make links between files |
| logname | Print user's login name |
| ls | List directory contents |
| md5sum | Compute MD5 checksums |
| mkdir | Make directories |
| mkfifo | Make FIFOs (named pipes) |
| mknod | Make block or character special files |
| mktemp | Create a temporary file or directory |
| mv | Move (rename) files |
| nice | Run a program with modified scheduling priority |
| nl | Number lines of files |
| nohup | Run a command immune to hangups |
| nproc | Print the number of processing units |
| numfmt | Reformat numbers |
| od | Dump files in octal and other formats |
| paste | Merge lines of files |
| pathchk | Check whether file names are valid or portable |
| pinky | Lightweight finger |
| pr | Paginate or columnate files for printing |
| printenv | Print all or part of environment |
| printf | Format and print data |
| ptx | Produce a permuted index of file contents |
| pwd | Print name of current working directory |
| readlink | Print value of a symbolic link |
| realpath | Print the resolved path |
| rm | Remove files or directories |
| rmdir | Remove empty directories |
| runcon | Run command with specified SELinux security context (stub) |
| seq | Print a sequence of numbers |
| sha1sum | Compute SHA-1 checksums |
| sha224sum | Compute SHA-224 checksums |
| sha256sum | Compute SHA-256 checksums |
| sha384sum | Compute SHA-384 checksums |
| sha512sum | Compute SHA-512 checksums |
| shred | Overwrite a file to hide its contents |
| shuf | Generate random permutations |
| sleep | Delay for a specified amount of time |
| sort | Sort lines of text files |
| split | Split a file into pieces |
| stat | Display file or file system status |
| stdbuf | Run a command with modified I/O stream buffering |
| stty | Change and print terminal line settings |
| sum | Checksum and count the blocks in a file |
| sync | Flush file system buffers |
| tac | Concatenate and print files in reverse |
| tail | Output the last part of files |
| tee | Read from stdin and write to stdout and files |
| test | Evaluate expression (also: [) |
| timeout | Run a command with a time limit |
| touch | Change file timestamps |
| tr | Translate or delete characters |
| true | Do nothing, successfully |
| truncate | Shrink or extend the size of a file |
| tsort | Perform topological sort |
| tty | Print the file name of the terminal |
| uname | Print system information |
| unexpand | Convert spaces to tabs |
| uniq | Report or omit repeated lines |
| unlink | Call the unlink function to remove the specified file |
| uptime | Tell how long the system has been running |
| users | Print the user names of users currently logged in |
| vdir | List directory contents verbosely (like ls -l -b) |
| wc | Print newline, word, and byte counts for each file |
| who | Show who is logged on |
| whoami | Print effective userid |
| yes | Output a string repeatedly until killed |

## Notes

- **Platform**: Primarily targets Linux. Some features (chroot, stty, uptime) use Linux-specific syscalls.
- **b2sum**: Uses SHA-512 as a stand-in for BLAKE2b. For true BLAKE2b, add `golang.org/x/crypto` to go.mod and update b2sum.go.
- **chcon/runcon**: Stubbed — require SELinux kernel support.
- **stty**: Limited terminal settings support.
- **users/who**: Limited utmp parsing; shows current user as fallback.
