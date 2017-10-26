# gilbert.nvim☄
![](https://travis-ci.org/NoahOrberg/gilbert.nvim.svg?branch=master)

## Description
`gilbert.nvim` is neovim plugin that easy file load, upload, update(patch) to gist using [NoahOrberg/gilbert](http://github.com/NoahOrberg/gilbert).

## Requirements
- only use MacOS
  - because this plugin is using `open` & `pbcopy` command.
  - Linux, and Windows does not execute Command.
- go
  - version `1.8` or `1.9`
- glide
  - But it is NOT using by [Installation - using only dein.vim](https://github.com/NoahOrberg/gilbert.nvim#using-only-deinvim) because `$XDG_CONFIG_HOME` is NOT included `$GOPATH`. And dependent packages are installed by `util/dep.sh` and it placed in `$GOPATH`. (cannot be vendoring)
- make

## Installation
0. Please set ENVIRONMENT VARIABLE because using `gilbert`.
``` sh
$ export GILBERT_GISTTOKEN=********
$ export GILBERT_GISTURL=https://api.github.com/gists
```

### Quickstart [using only dein.vim]
1. Please write your `init.vim`.
``` vim
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
2. Restart `nvim`.

Perhaps, some problems may happen ;(  
Therefore, will recommended another install-method.

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


## VARIABLE 
``` vim
g:gilbert#allow_open_by_browser=1 " allow open browser when `:GiUpload` or `:GiPatch`
g:gilbert#should_copy_url_to_clipboard=1 " allow copy URL to clipboard
```

## How To Use
- Upload current buffer
  - `<FILENAME>` is optional. For example when use it if buffer is `[No Name]`.
  - If this command is success, output `URL`.
``` vim
:GiUpload <FILENAME>
```
- Load new buffer.
  - Load Gist.
  - Open some buffer.
  - Auto save gist-files in `~/.gilbert/<gist_id>/<file_name>` when you load gist.
    - However, It's just a workspace(you can use `quickrun`, `syntax highliight` and so on).
``` vim
:GiLoad <GIST-ID or GIST-URL>
```
- Update gist
  - Upload all gist-file related from current buffer to gist.
  - Should be load by `:GiLoad <GIST-ID>` to current buffer or execute `:GiUpload` command from `NoName` buffer before execute this command.
  - And after execute this command, related buffer will be closed.
``` vim
:GiPatch
```

