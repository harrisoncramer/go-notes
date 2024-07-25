# Go Notes

This is a very simple Charm CLI tool for writing files like diary entries or todo items. 

It writes the files to a sqlite3 database directly from the terminal, using an editor like Vim or Neovim.

## Installation

```bash
go install github.com/harrisoncramer/go-notes/cmd/go-notes@latest`
```

## Usage

```bash
$ go build .
$ ./go-notes "My Diary"
```

You can make as many databases as you want.
