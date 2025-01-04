# gopkgviewer - Go dependency visualization

<p align="center">
  <img src="./frontend/public/favicon.png" width="150">
   <br />
   <strong>Status: </strong>Maintained
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/tag/antonhancharyk/gopkgviewer" alt="GitHub tag">
  <img src="https://goreportcard.com/badge/github.com/antonhancharyk/gopkgviewer" alt="Go Report Card">
  <img src="https://github.com/antonhancharyk/gopkgviewer/actions/workflows/release.yml/badge.svg" alt="Build Status">
</p>

**gopkgviewer** is an interactive tool designed to visualize and analyze Go project dependencies. It provides a rich, web-based interface for better understanding of how your project connects its components and external libraries.

## Features

- Interactive web-based visualization of Go dependencies
- Toggle dependencies by type
- Focus on specific dependencies for analysis

## Installation - 2 options

### Install via `go install`

```bash
go install github.com/antonhancharyk/gopkgviewer@latest
```

### Download the Release

From the latest release from the [Releases Page](https://github.com/antonhancharyk/gopkgviewer/releases).


## Usage

Navigate to your Go project directory and run:

```bash
cd my-go-project
gopkgviewer
```

This will start a web server with the dependency visualization available in your browser.

### Available Flags

```plaintext
--gomod value           Path to go.mod
--addr value            Address to listen on (default: :0)
--skip-browser          Do not open browser (default: true)
--help, -h              Show help
--version, -v           Print the version
```

## Alternatives

- [go-callvis](https://github.com/ondrajz/go-callvis) - Great tool for visualizing of call, but panic on Go >= 1.21
- [godepgraph](https://github.com/kisielk/godepgraph) - Same idea, but output is static image
- [depgraph](https://github.com/becheran/depgraph)
- [gomod](https://github.com/Helcaraxan/gomod)

## License

© 2025 [Anton Hancharyk](https://github.com/antonhancharyk)  
This project is [GPL-3.0 license](./LICENSE) licensed.
