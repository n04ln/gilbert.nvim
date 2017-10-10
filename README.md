# gilbert.nvim☄
![](https://travis-ci.org/NoahOrberg/gilbert.nvim.svg?branch=master)

## Description
`gilbert.nvim` is neovim plugin that easy file upload to gist using [NoahOrberg/gilbert](http://github.com/NoahOrberg/gilbert).

## Requirements
- only use MacOS
  - because this plugin is using `open` & `pbcopy` command(to open Upload gist).
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

### using only dein.vim
1. Please write your `init.vim`.
``` vim
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
2. Restart `nvim`.

### using console and dein.vim(RECOMMENDED)
1. Please `go get` this repository.
``` sh
$ go get github.com/NoahOrberg/gilbert.nvim
```
2. Please build it.
``` sh
$ make build
```
3. Please write your `init.vim`
``` vim
call dein#add('NoahOrberg/gilbert.nvim')
```
4. Restart `nvim`.


## Variable
``` vim
g:gilbert#allow_open_by_browser=1 " allow open browser when `:GiUpload` or `:GiPatch`
g:gilbert#should_copy_url_to_clipboard=1 " allow copy URL to clipboard
```

## How To Use
- Upload current buffer
  - `<FILENAME>` is optional, when use it if buffer is `[No Name]`.
  - If this command is success, output `URL`.
  - If it is success and `g:gilbert#is_allow_open_brower==1`, Open your browser.
``` vim
:GiUpload <FILENAME>
```
- Load new buffer.
  - Load Gist of only one file.
``` vim
:GiLoad <GIST-ID or GIST-URL>
```
- Update gist
  - Upload gist-file from current buffer to gist.
  - If it is success and `g:gilbert#is_allow_open_brower==1`, Open your browser.
  - Should be load by `:GiLoad <GIST-ID>` to current buffer before execute OR already `:GiUpload` command from `NoName` buffer.
``` vim
:GiPatch
```

