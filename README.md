# Git CLI Tool for macOS

## Project Structure
```
git-cli-tool/
├── cmd/
│   └── main.go
├── internal/
│   ├── cli/
│   │   ├── input.go
│   │   └── menu.go
│   └── git/
│       ├── commands.go
│       ├── squash.go
│       └── utils.go
├── go.mod
└── README.md
```

# How to Run 


## To build the binary file and to move it to your bin folder.

```bash
cd <your-directory-where-you-placed-it>
go build -o git-cli ./cmd/main.go #To build the binary file.
sudo mv git-cli /usr/local/bin/
```

## Squashing the commits

```bash
cd /path/to/your/other/repository
git-cli
```

## How to delete the binary now?

```bash
sudo rm /usr/local/bin/git-cli  #To delete the file from the entire system.
```
