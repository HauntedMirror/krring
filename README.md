<div align="right">

![golangci-lint](https://github.com/HauntedMirror/krring/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/HauntedMirror/krring/actions/workflows/release.yml/badge.svg)

</div>


<div align="center">

# 🥀🖤🎀 krring

`ping` command implementation in Go but with colorful output and krrg ascii art


![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/HauntedMirror/krring?style=flat-square)](https://github.com/HauntedMirror/krring/releases/latest)

</div>

## Features

- [x] Colorful and fun output.
- [x] Cross-platform support: Windows, macOS, and Linux (also WSL)
- [x] It works with a single executable file, so it can be installed easily.
- [x] Support IPv4 and IPv6.

## Usage

Simply specify the target host name or IP address in the first argument e.g. `krring github.com` or `krring 13.114.40.48`.
You can change the number of transmissions by specifying the `-c` option.

```
Usage:
  krring [OPTIONS] HOST

`ping` command but with krrg

Application Options:
  -c, --count=     Stop after <count> replies (default: 20)
  -P, --privilege  Enable privileged mode
  -V, --version    Show version

Help Options:
  -h, --help       Show this help message
```

## Installation

### Requirements

- Terminals compatible with True Color

### Download executable binaries

You can download executable binaries from the latest release page.

> [![Latest Release](https://img.shields.io/github/v/release/HauntedMirror/krring?style=flat-square)](https://github.com/HauntedMirror/krring/releases/latest)

### Build from source

To build from source, clone this repository then run `make build` or `go install`. Develo*ping* on `go1.18.3 linux/amd64`.

## LICENSE

[MIT](./LICENSE)

## Credits

Model VTuber: [枢木みるは](https://www.youtube.com/@krrg_mrh)

Pixel Art: [衣乃環ゆひ](https://coconala.com/users/3868492)

## Author
[HauntedMirror](https://twitter.com/HauntedMirror)

Original pingu: [Sheepla](https://github.com/sheepla)
