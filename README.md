# Filesystem MCP

This repository provides an implementation of the MCP to offer a suite of tools for interacting with the file system, such as listing directory entries, reading and writing files, retrieving file information, and renaming and copying files or directories.

It allows users to run using stdio or SSE server on a local machine.

## Table of Contents

- [Installation](#installation)
  - [Installing Locally by Cloning the Repository](#installing-locally-by-cloning-the-repository)
  - [Using Docker](#using-docker)
- [How to Use](#how-to-use)
  - [Example of usage with PydanticAI in Python](#example-of-usage-with-pydanticai-in-python)
    - [Using SSE server](#using-sse-server)
    - [Using stdio](#using-stdio)
- [Tool Descriptions](#tool-descriptions)

## Installation

1. Ensure that you have [Go](https://golang.org/doc/install) installed on your system. Use Go version 1.24.
2. Run the following command to install the package using `go install`:

```bash
go install github.com/lealre/fs-mcp@latest
```

3. Run the help command:

```bash
fs-mcp -h
```

- The `-dir` flag specifies the base directory that the server will serve. It is required.
- The `-port` flag specifies the port on which the server will listen. It is optional, with the default being `8080`.
- The `-t` flag specifies the transport type. It can be either `stdio` or `http` (default is `stdio`).

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
./fs-mcp -t http -dir /your/directory/path
```

- This will run the server at `http://localhost:8080`.

### Using Docker

This MCP can only be run in Docker using the SSE option. To do so, create a volume pointing to the directory you want to expose for file system operations.

You can pull the image from Docker Hub and run it as follows:

```shell
docker pull lealre/fs-mcp:latest
```

Or build the image from the cloned repository:

```shell
docker build -t fs-mcp .
```

Run the container, mapping your desired directory and port:

```shell
docker run -p 8080:8080 -v /your/directory/path:/baseDir lealre/fs-mcp -volume "/your/directory/path:/baseDir"
```

- Replace `/your/directory/path` with the actual path you want to serve.
- The `-p` flag maps your local machine's port to the Docker container's port.
- The `-v` flag specifies the path to the base directory similarly in the `--volume` flag to ensure Docker has access to your files.

This setup will start the server at `http://localhost:8080`, serving the specified directory.

> [!IMPORTANT]
> For this to work properly, ensure the paths in `-v` and `-volume` match exactly.

The container’s base folder volume name (`/baseDir`) can be customized, and the port can be changed using the `-port` flag (along with adjusting `-p` in Docker). Example:

```shell
docker run -p 8081:8081 -v /your/directory/path:/baseDir lealre/fs-mcp -volume "/your/directory/path:/baseDir" -port 8081
```

## How to Use

Once the installation is complete, you can use the server by running:

```bash
fs-mcp -t http -dir /your/directory/path
```

This will start the MCP server at `http://localhost:8080`, restricting the file system operations to be under this specific path.

For now, it only accepts one path to the server.

To change the server port, you can pass it as a flag:

```bash
fs-mcp -dir /your/directory/path -port 3000
```

To run it using stdio, just omit the flag `-t`.

### Example of usage with PydanticAI in Python

#### Using SSE server

- Run the MCP server:

```bash
fs-mcp -dir /your/directory/path
```

- Create the Python client (using OpenAI in this case):

```python
# script.py

import asyncio
import sys
from pydantic_ai import Agent
from pydantic_ai.mcp import MCPServerSSE

server = MCPServerSSE('http://localhost:8080/sse')
agent = Agent('openai:gpt-4o', mcp_servers=[server])

async def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py 'your prompt/query here'")
        sys.exit(1)

    query = sys.argv[1]

    async with agent.run_mcp_servers():
        result = await agent.run(query)

    print(result.output)

if __name__ == '__main__':
    asyncio.run(main())
```

- Then you can run:

```bash
python script.py "List all the entries for the path in /your/directory/path/somesubpath"
```

#### Using stdio

- Create the Python client (using OpenAI in this case):

```python
# script.py

import asyncio
import sys
from pydantic_ai import Agent
from pydantic_ai.mcp import MCPServerStdio

base_path = "/your/directory/path"
server = MCPServerStdio(
    'fs-mcp',
    args=['-dir', base_path]
)
agent = Agent('openai:gpt-4o', mcp_servers=[server])

async def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py 'your prompt/query here'")
        sys.exit(1)

    query = sys.argv[1]

    async with agent.run_mcp_servers():
        result = await agent.run(query)

    print(result.output)

if __name__ == '__main__':
    asyncio.run(main())
```

- Then you can run:

```bash
python script.py "List all the entries for the path in /your/directory/path/somesubpath"
```

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

This project uses the [mcp-go library](https://pkg.go.dev/github.com/mark3labs/mcp-go/mcp) to implement core functionality.
