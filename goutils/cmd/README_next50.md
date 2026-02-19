# goutils — Next 50

Each util lives in its own directory with a single `<name>.go` file.
All utils use only the Go standard library — no external dependencies.

| # | Util | Description |
|---|------|-------------|
| 1 | **abs** | Absolute value of numbers (args or stdin) |
| 2 | **anagram** | Check if two words are anagrams; group anagram sets from wordlist |
| 3 | **ascii** | Display ASCII table or look up char/codepoint info |
| 4 | **banner** | Print large 5-row ASCII-art text banners |
| 5 | **bitflip** | Flip all bits or specific bit positions in an integer |
| 6 | **bytecount** | Count bytes, chars, words, lines (readable `wc`) |
| 7 | **charfreq** | Character frequency analysis with histogram |
| 8 | **checksum** | MD5 + SHA1 + SHA256 + SHA512 + CRC32 in one pass |
| 9 | **chunk** | Split stdin lines into chunks of N, separated by delimiter |
| 10 | **cipher** | Classical ciphers: rot13, rot47, caesar, atbash, vigenere |
| 11 | **clip** | Copy stdin to clipboard (pbcopy/xclip/xsel) or paste |
| 12 | **colorize** | Highlight regex pattern matches with ANSI colours |
| 13 | **countdown** | Live countdown timer from N seconds |
| 14 | **dateadd** | Add/subtract durations from dates (days, weeks, months, years) |
| 15 | **dedent** | Remove common leading whitespace from all lines |
| 16 | **dice** | Roll dice in standard notation: `2d6`, `1d20+5`, `4d6kh3` |
| 17 | **dotenv** | Load `.env` file and print exports, or inject vars into a command |
| 18 | **duration** | Parse & convert durations between units; humanize seconds |
| 19 | **eol** | Detect and convert line endings (LF ↔ CRLF ↔ CR) |
| 20 | **escape** | Escape/unescape: shell, regex, SQL, Go strings, XML |
| 21 | **fieldmap** | Rearrange, rename, or filter CSV/TSV fields |
| 22 | **fileage** | Show how old files are in human-readable form |
| 23 | **filehead** | Show first N **bytes** of files (binary-safe) |
| 24 | **filetail** | Show last N **bytes** of files (binary-safe) |
| 25 | **flip** | Flip text upside-down (Unicode) or mirror left-right |
| 26 | **floatfmt** | Format floats: precision, scientific notation, comma separators |
| 27 | **fuzz** | Generate fuzzing inputs: strings, ints, floats, bytes, boundary |
| 28 | **histogram** | ASCII histogram from numeric stdin data |
| 29 | **htmlstrip** | Strip HTML tags; optionally decode entities |
| 30 | **humanize** | Convert raw numbers to human bytes/count/duration |
| 31 | **indent2tab** | Convert space-indented code to tabs (or vice versa) |
| 32 | **initcap** | Title-case text, respecting common lowercase words |
| 33 | **iprange** | Expand CIDR, list IP ranges, check membership |
| 34 | **isutf8** | Validate UTF-8; report bad byte positions |
| 35 | **jsonformat** | Pretty-print or minify JSON; optional ANSI color output |
| 36 | **jsonkeys** | List all dot-notation keys/paths in a JSON document |
| 37 | **keygen** | Generate hex, base64, tokens, UUIDs, PINs, passphrases |
| 38 | **kwsearch** | Keyword search with context lines and match highlighting |
| 39 | **linediff** | Side-by-side line diff of two files with colour |
| 40 | **linenum** | Add or remove line numbers |
| 41 | **linesplit** | Word-wrap or hard-wrap lines at a column limit |
| 42 | **lorem** | Generate Lorem Ipsum (words, sentences, paragraphs) |
| 43 | **matrix** | Matrix math: transpose, multiply, add, sub, stats |
| 44 | **morse** | Encode/decode Morse code |
| 45 | **palindrome** | Check if text is a palindrome; filter palindromes from stdin |
| 46 | **passgen** | Generate secure passwords with configurable rules |
| 47 | **pathinfo** | Dissect file paths into dir, base, stem, ext, abs |
| 48 | **pipe** | Run a sequence of commands as an explicit pipeline |
| 49 | **pluralize** | Pluralize English words (with irregular/invariant support) |
| 50 | **ratelimit** | Rate-limit stdin lines: N lines per second/interval |
| 51 | **roman** | Convert between integers and Roman numerals |
| 52 | **signame** | Convert between Unix signal numbers and names |
| 53 | **table** | Render CSV/TSV as a formatted ASCII or Unicode box table |
| 54 | **template** | Render Go text templates with env vars or JSON data |
| 55 | **textcount** | Count sentences, paragraphs, FK grade level, reading time |
| 56 | **timer** | Interactive stopwatch with lap support |
| 57 | **treeprint** | Render directory tree as Unicode tree (like `tree`) |
| 58 | **xmlfmt** | Pretty-print or minify XML |
| 59 | **zigzag** | Rail Fence cipher encode/decode |

## Build all

```bash
for dir in */; do
  [ -f "$dir"/*.go ] && (cd "$dir" && go build -o "${dir%/}" *.go && echo "✓ ${dir%/}")
done
```
