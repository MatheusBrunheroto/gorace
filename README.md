# gorace
A race condition testing tool written in Go. gorace queues HTTP requests and releases them through multiple timing modes (flood, cascade, round-cascade and more), instead of firing them one at a time like typical fuzzers.

<!-- Demo video -->
<!-- 
https://github.com/user-attachments/assets/YOUR-VIDEO-ID-HERE
-->

<hr>

## Installation

### Go Installer
```bash
go install github.com/MatheusBrunheroto/gorace@latest
```

### AUR package
```bash
yay -S NOT WORKING YET
```

### Manual (Building from source)
```bash
git clone https://github.com/MatheusBrunheroto/gorace.git
cd gorace
chmod +x install.sh
./install.sh
```
This will compile the project and install the `gorace` binary to `/usr/local/bin`, making it available globally in your terminal.

**Requirements:** Go 1.25 or higher.

<hr>

## Usage

Usage Example:
```bash
gorace -u 'https://x.com' -d 'FUZZ=payload' -w 'FUZZ=./wordlist.txt' --threads 10 --delay 100 \
       -u 'https://y.com' --threads 10 --delay 1000 --verbose 3 --match e
```

Every `-u` starts a new target. Any flags set before the next `-u` (headers, cookies, data, wordlist, threads, delay) apply only to that target — so a single command can test multiple websites at once, each with its own configuration.

<!-- Usage gifs -->
<!--
### Flood mode (verbose 3)
![flood mode](./docs/gifs/flood-verbose3.gif)

### Sequential mode (verbose 3)
![sequential mode](./docs/gifs/sequential-verbose3.gif)

### Cascade mode (verbose 3)
![cascade mode](./docs/gifs/cascade-verbose3.gif)

### Round-Cascade mode (verbose 3)
![round-cascade mode](./docs/gifs/round-cascade-verbose3.gif)
-->

### Manual Page
```text
GORACE(1)

NAME
    gorace — tool for testing web race conditions through concurrent HTTP requests

SYNOPSIS
    gorace [OPTIONS]

DESCRIPTION
    gorace is a command-line tool designed to test for race conditions in web applications.
    It sends multiple concurrent HTTP requests to one or more target endpoints in order to
    detect timing vulnerabilities that occur when requests are processed simultaneously.

    Multiple targets can be defined in a single command by repeating the -u flag. Every flag
    set before a -u applies to that specific target, allowing different URLs, headers,
    cookies, request bodies, thread counts and delays to be tested together in the same run.

    Requests are released according to the selected mode (flood, sequential, cascade, or
    round-cascade), controlling the timing strategy used to trigger race conditions.

OPTIONS
    -u, --url URL
        Target URL. Can be used multiple times to define several targets in the same run.

    -X, --method METHOD
        HTTP method to use (GET, POST, PUT, PATCH, DELETE).

    -d, --data DATA
        Request body data. Supports FUZZ placeholders when used with -w.

    -H, --headers HEADERS
        Custom HTTP headers.

    -b, --cookies COOKIES
        Cookies to include in the request.

    -w, --wordlist WORDLIST
        Wordlist used to replace FUZZ placeholders and generate multiple requests.

    -t, --threads NUMBER
        Number of concurrent workers for the target.

    -D, --delay MILLISECONDS
        Delay applied between requests, depending on the selected mode.

    -m, --mode MODE
        Timing mode used to release requests: flood, sequential, cascade, round-cascade.

    -v, --verbose LEVEL
        Verbosity level (0-3). Higher levels show per-request details, including cache hits
        and errors.

    --match STRING
        Filters and highlights responses containing the given string.

    -h, --help
        Display help information.

EXAMPLES
    Basic race condition test against two targets:
        gorace -u '1.com' --threads 10 --delay 100 -u '2.com' --threads 10 --delay 1000 --verbose 3

    Fuzzing a request body with a wordlist:
        gorace -u 'x.com' -d 'FUZZ=payload' -w 'FUZZ=./wordlist.txt' \
               --threads 10 --delay 100 --verbose 2 --match 'file'

    Testing multiple different targets in the same run:
        gorace -u 'x.com' -d 'FUZZ=payload' -w 'FUZZ=./wordlist.txt' --threads 10 --delay 100 \
               -u 'y.com' --threads 10 --delay 1000 --verbose 3 --match e

EXIT STATUS
    0
        Successful execution.

    1
        Execution error.

SEE ALSO
    curl(1), ffuf(1), nmap(1)

AUTHOR
    Matheus Brunheroto

COPYRIGHT
    MIT License
```

<hr>

## Inspiration

This repository was created to explore the [Race Condition](https://portswigger.net/web-security/race-conditions) challenges from PortSwigger Academy while learning [Go](https://go.dev/). The CLI was inspired by one of the most powerful fuzzing tools available, [FFUF](https://github.com/ffuf/ffuf).
