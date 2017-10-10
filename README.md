# gilbert.nvim

## Description
`gilbert.nvim` is neovim plugin that easy file upload to gist using [NoahOrberg/gilbert](http://github.com/NoahOrberg/gilbert).

## Requirements
- only use MacOS
  - because this plugin is using `open` command(to open Upload gist).
  - Linux, and Windows does not execute Command.
- go
- glide
  - But it is NOT using now because `$XDG_CONFIG_HOME` is NOT included `$GOPATH`. And dependent packages are installed by `util/dep.sh` and it placed `$GOPATH`. (cannot be vendoring)
- make

## Installation
0. Please set ENVIRONMENT VARIABLE because using `gilbert`.
``` sh
$ export GIST_TOKEN=<YOUR TOKEN HERE>
```
1. Please write your `init.vim`.
``` vim
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
2. Restart `nvim`.

## Variable
``` vim
g:gilbert#is_allow_open_brower=1 " allow open browser when `:GiUpload` or `:GiPatch`
```

## How To Use
- Upload current buffer
  - `<FILENAME>` is optional, when use it if buffer is `[No Name]`.
  - If this command is success, output `URL`.
  - If it is success and `g:gilbert#is_allow_open_brower==1`, Open your browser.
``` vim
:GiUpload <FILENAME>
```
- Load current buffer.
  - Load Gist of only one file.
``` vim
:GiLoad <GIST-ID>
```
- Update gist
  - Upload gist-file from current buffer to gist.
  - If it is success and `g:gilbert#is_allow_open_brower==1`, Open your browser.
  - Should be load by `:GiLoad <GIST-ID>` to current buffer before execute OR already `:GiUpload` command from `NoName` buffer.
``` vim
:GiPatch
```

