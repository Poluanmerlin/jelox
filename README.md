# Jelox Repo Analyzer

Jelox is a powerful tool designed to generate repositories and local directories wordlists that can be used for directory and file enumeration. It is also useful for web developers to generate wordlists for analyzing status codes on files and paths, providing insights into file structures and content. Whether you want to analyze a single repository, multiple repositories from a list, or a local directory, Jelox has got you covered.
## Features

- Clone and analyze single or multiple repositories.
- Process local directories to list files.
- Supports input via command line, list files, or standard input.
- Output results to a file or stdout.

## Installation

Ensure Go is installed on your system. Clone the repository and build the tool:

```bash
git clone https://github.com/yourusername/jelox.git
cd jelox
go build -o jelox main2.go
```

## Usage

```bash
go run main2.go [OPTIONS]
```

### Options

- `-o FILE`  : Output file (default: `wordlist.txt`).
- `-d DIR`   : Process a local directory.
- `-l FILE`  : List file containing repository URLs.
- `-u URL`   : Single repository URL to process.
- `-h`       : Show help message.

### Examples

Clone and analyze a single repository:
```bash
go run main2.go -u https://github.com/example/repo.git
```

Process a local directory:
```bash
go run main2.go -d /path/to/directory
```

Analyze multiple repositories from a list file:
```bash
go run main2.go -l repos.txt
```

Specify an output file:
```bash
go run main2.go -u https://github.com/example/repo.git -o output.txt
```

## Contributing

Contributions are welcome! Fork the repository, make changes, and submit a pull request.

## License

This project is licensed under the MIT License.

## Acknowledgements

Thanks to the Go community for the robust tools and libraries that made this project possible.
