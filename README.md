# Filesystem MCP Server in Go

This repository provides a server implementation of a MCP to offer a suite of tools that allow interaction with the file system, such as listing directory entries, reading and writing files, retrieving file information, renaming, and copying files or directories.

The server runs on a local machine and listens for commands, making it a powerful utility for automated scripts or remote file management tasks.

## Installation

1. Ensure that you have [Go](https://golang.org/doc/install) installed on your system (using Go version 1.24).
2. Run the following command to install the package using `go install`:

```bash
go install github.com/lealre/fs-mcp@latest
```

3. Run the help command:

```bash
fs-mcp -h
```

- The `-dir` flag specifies the base directory that the server will serve. It is required.
- The `-port` flag specifies the port on which the server will listen. It is optional, and the default is `8080`.

### Installing Locally by Cloning the Repository

1. Clone the repository to your local machine:

```bash
git clone https://github.com/lealre/fs-mcp.git
```

2. Navigate into the cloned directory:

```bash
cd fs-mcp
```

3. Build the project:

```bash
go build
```

4. Run the executable:

```bash
./fs-mcp
```

- This will run the server, and you can specify options such as `-dir` and `-port` with the command.

## Tool Descriptions

This project provides various tools to interact with the file system. Below are the descriptions of each tool:

- **listEntries**: List entries at a given path. Parameters:

  - `path` (string, required): Path for which to list all entries.
  - `depth` (number, optional): Depth of the directory tree (default is 3).

- **readFromFile**: Read the contents of a file at a given path. Parameters:

  - `path` (string, required): Path to the file to be read.

- **writeToFile**: Create or overwrite a file with the given content. Parameters:

  - `path` (string, required): Path to the file to write to.
  - `content` (string, required): Content to write to the file.

- **getFileInfo**: Retrieve file information including size, last modified time, detected MIME type, and file permissions. Parameters:

  - `path` (string, required): Path to the file to retrieve information from.

- **renamePath**: Renames a file or directory to a new name. Parameters:

  - `path` (string, required): Path to the file or directory to be renamed.
  - `newPathFinalName` (string, required): New name for the file or directory (just the name, not the full path).

- **copyFileOrDir**: Copies a file or directory to a new location. Parameters:
  - `path` (string, required): Path to the file or directory to be copied.
  - `destination` (string, required): Destination path where the file or directory will be copied.

[mcp-go docs](https://pkg.go.dev/github.com/mark3labs/mcp-go/mcp)
