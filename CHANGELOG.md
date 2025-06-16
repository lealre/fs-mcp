
<a name="v0.1.4"></a>
## [v0.1.4](https://github.com/lealre/fs-mcp/compare/v0.1.3...v0.1.4) (2025-06-16)

* Add option to run using Docker ([#3](https://github.com/lealre/fs-mcp/issues/3))
* Update `CHANGELOG.md` for `v0.1.3` release

<a name="v0.1.3"></a>
## [v0.1.3](https://github.com/lealre/fs-mcp/compare/v0.1.2...v0.1.3) (2025-06-14)

* Add a new struct as the return type for file system operations ([#2](https://github.com/lealre/fs-mcp/issues/2))
* Update help text to include new transport flag (-t)
* Update description in `README.md` to include new transport type flag
* Add some tests to `copyFileOrDir`
* Update ``CHANGELOG.md`

<a name="v0.1.2"></a>
## [v0.1.2](https://github.com/lealre/fs-mcp/compare/v0.1.1...v0.1.2) (2025-06-09)

* Add `stdio` as the default transport option ([#1](https://github.com/lealre/fs-mcp/issues/1))
* Update `README.md` to mention  default transport as `stdio`
* Add new default transport as `stdio`
* Add test to `renamePath` function from `fs_operations.go`

<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/lealre/fs-mcp/compare/v0.1.0...v0.1.1) (2025-06-05)

* Add `CHANGELOG.md` reflecting `v0.1.1`
* Update changelog templae to keep it simple for now
* Add tests to `getFileInfo`
* Add changelog configurations

<a name="v0.1.0"></a>
## v0.1.0 (2025-06-04)

* Update message when using with no `-dir` arg
* Add tests to readFile
* Fix prefix in tests to list listEntries
* Add mote tests cases to listEntries
* Add test to listEntries
* Add How to use section to README.md
* Add installation section to README.md
* Add '-h' flag description
* Extend it to support the `--port` and `--dir` flags
* Add a brief description of the repository and its tools
* Fix depth default value for `listEntries` tools and expand to copy Files or Dirs
* Add copy file operation
* Check if path exists before verifying whether it's a file or directory
* Add tool to rename a path
* Start adding tests for filesystem operations
* Add `.vscode` and `ttodo.md` in `.gitignore`
* Add a safe path check that allows being just in a specific path
* Add depth to list entries and add the raw implementation of baseEntry. Correct the English.
* Add tool to retrieve information from a specific file
* Change to not return error when a path not found
* Refactor code and add logs
* Add compiled Go binary to `.gitignore`
* Add two tools: read and write files
* Basic SSE MCP server with one tool
* Add function to list entries from a path
* Initial commit
