# gorace
Web race condition testing tool written in Go
NOT FINISHED YET

Works like ffuf, but instead of instantly activating a request, makes a queue that fire all the requests at the same time

<hr>

## Installation
### Go Installer
```bash
go install NOT WORKING YET
```
### AUR package
```bash
yay -S NOT WORKING YET
```
### Manual (Building)
```bash
git clone https://github.com/MatheusBrunheroto/gorace.git
cd gorace
chmod +x install.sh
./install.sh
NOT WORKING
```
<hr>

## Usage
Usage Example:

```gorace -u 'https://website.com' -H '{h1_key:h1_value, h2_key:WORDLIST1}' -c {WORDLIST2:WORDLIST3} -t 50```

### Manual Page
```text
GORACE(1)

NAME
    gorace — tool for testing web race conditions through concurrent HTTP requests

SYNOPSIS
    gorace [OPTIONS]

DESCRIPTION
    gorace is a command-line tool designed to test for race conditions in web applications.
    It sends multiple concurrent HTTP requests to a target endpoint in order to detect timing vulnerabilities that occur when requests are processed simultaneously.

    The tool allows the user to configure HTTP methods, headers, cookies, request bodies, and the number of concurrent workers.
    All workers are synchronized so that requests are released at the same time, increasing the likelihood of triggering race conditions.

OPTIONS
    -U, --url URL
        Target URL.

    -X, --method METHOD
        HTTP method to use (GET, POST, PUT, PATCH, DELETE).

    -d, --data DATA
        Request body data.

    -H, --header HEADER
        Custom HTTP header. Can be used multiple times.

    -C, --cookie COOKIE
        Cookie to include in the request.

    -t, --threads NUMBER
        Number of concurrent workers.

    -w, --wordlist FILE
        Wordlist used to generate multiple requests.

    -i, --interactive
        Interactive mode for entering request parameters.

    -h, --help
        Display help information.

EXAMPLES
    Send 100 concurrent POST requests:

        gorace -U https://target.com/redeem -X POST -d "coupon=FREE100" -t 100

    Send requests with custom headers and cookies:

        gorace -U https://target.com/api -X POST -d "id=1" \
        -H "Authorization: Bearer TOKEN" \
        -C "session=abcd123"

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
