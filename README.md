# Filesystem MCP Server in Go

This repository provides a server implementation of a MCP to offer a suite of tools that allow interaction with the file system, such as listing directory entries, reading and writing files, retrieving file information, renaming, and copying files or directories.

The server runs on a local machine and listens for commands, making it a powerful utility for automated scripts or remote file management tasks.

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
