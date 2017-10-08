# gilbert.nvim

## Desctiption
`gilbert.nvim` is neovim plugin that easy file upload to gist using [NoahOrberg/gilbert](http://github.com/NoahOrberg/gilbert).

## Requirements
- go
- glide
  - But it is NOT using now because `$XDG_CONFIG_HOME` is NOT included `$GOPATH`. And dependent packages are installed by `util/dep.sh` and it placed `$GOPATH`. (cannot be vendoring)
- make

## Installation
0. Please set ENVIRONMENT VARIABLE because using `gilbert`.
```
$ export GIST_TOKEN=<YOUR TOKEN HERE>
```
1. Please write your `init.vim`.
```
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
2. Restart `nvim`.

## How To Use
- Upload current buffer
  - `<FILENAME>` is optional, when use it if buffer is `[No Name]`.
  - If this command is success, output `URL`.
```
:GiUpload <FILENAME>
```
- Load current buffer
```
:GiLoad <GIST_ID>
```
