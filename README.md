# gilbert.nvim

## Desctiption
`gilbert.nvim` is neovim plugin that easy file upload to gist using [NoahOrberg/gilbert](http://github.com/NoahOrberg/gilbert) 

## Requirements
- go
- glide
- make

## Installation
0. Please set ENVIRONMENT VARIABLE because using `gilbert`
```
$ export GIST_TOKEN=<YOUR TOKEN HERE>
```
1. Please get dependent package because `$XDG_CONFIG_HOME` is NOT include `$GOPATH`
```
$ go get github.com/neovim/...
$ go get github.com/NoahOrberg/gilbert.nvim/...
```
2. Please write your `init.vim`
```
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
3. Restart `nvim`

## How To Use
- Upload current buffer
  - `<FILENAME>` is optional, when use it if buffer is `[No Name]`.
```
:GiUpload <FILENAME>
```
