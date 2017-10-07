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
export GIST_TOKEN=<YOUR TOKEN HERE>
```
1. Please write your `init.vim`
```
call dein#add('NoahOrberg/gilbert.nvim', {'build' : 'make'})
```
2. Restart `nvim`

## How To Use
- Upload current buffer
  - <FILENAME> is optional, when use it if buffer is `[No Name]`.
```
:GiUpload <FILENAME>
```
