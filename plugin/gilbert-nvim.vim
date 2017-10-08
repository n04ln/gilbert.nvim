if exists('g:loaded_gilbert')
  finish
endif
let g:loaded_gilbert = 1

function! s:RequireGilbert(host) abort
  return jobstart(['gilbert.nvim'], { 'rpc': v:true })
endfunction

call remote#host#Register('gilbert.nvim', '0', function('s:RequireGilbert'))
call remote#host#RegisterPlugin('gilbert.nvim', '0', [
  \ {'type': 'command', 'name': 'GiUpload', 'sync': 0, 'opts': {'nargs': '?'}},
  \ {'type': 'command', 'name': 'GiLoad', 'sync': 0, 'opts': {'nargs': '1'}},
  \ {'type': 'command', 'name': 'GiPatch', 'sync': 0, 'opts': {'nargs': '1'}},
  \ ])
