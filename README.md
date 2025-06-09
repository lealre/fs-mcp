# Filesystem SSE Server in Go

This repository provides a server implementation to offer a suite of tools for interacting with the file system, such as listing directory entries, reading and writing files, retrieving file information, renaming, and copying files or directories.

It specifically runs an SSE server on a local machine.

## Table of Contents

- [Installation](#installation)
  - [Installing Locally by Cloning the Repository](#installing-locally-by-cloning-the-repository)
- [How to Use](#how-to-use)
  - [Example of usage with PydanticAI in Python](#example-of-usage-with-pydanticai-in-python)
    - [Using SSE server](#using-sse-server)
    - [Using stdio](#using-stdio)
- [Tool Descriptions](#tool-descriptions)

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
- The `-port` flag specifies the port on which the server will listen. It is optional, with the default being `8080`.

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
from pydantic_ai.mcp import MCPServerHTTP

server = MCPServerHTTP('http://localhost:8080/sse')
agent = Agent('openai:gpt-4o', mcp_servers=[server])

async def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py 'your prompt/query here'")
        sys.exit(1)

    query = sys.argv[1]

    async with agent.run_mcp_servers():
        result = await agent.run(query)

    print(result.data)

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

    print(result.data)

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
