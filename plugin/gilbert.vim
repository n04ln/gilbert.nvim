if exists('g:loaded_gilbert')
  finish
endif
let g:loaded_gilbert = 1

function! s:RequireGilbert(host) abort
  return jobstart(['gilbert.nvim'], { 'rpc': v:true })
endfunction

" g: gilbert # is_loaded_by_giload is a flag indicating whether or not it was loaded by `:GiLoad`.
"   e.g, {'gist_id' : 1}
let g:gilbert#is_loaded_by_giload={}
" g:gilbert#buffer_and_gist_id_info is buffer information
"   e.g, {'buffer_id' : 'gist_id'}
let g:gilbert#buffer_and_gist_id_info={}
" flush undo history
function! Gi_clear_undo() abort
    let old_undolevels = &undolevels
    setlocal undolevels=-1
    execute "normal! a \<BS>\<ESC>"
    let &l:undolevels = old_undolevels
endfunction

call remote#host#Register('gilbert.nvim', '0', function('s:RequireGilbert'))
call remote#host#RegisterPlugin('gilbert.nvim', '0', [
  \ {'type': 'command', 'name': 'GiUpload', 'sync': 1, 'opts': {'nargs': '?'}},
  \ {'type': 'command', 'name': 'GiLoad', 'sync': 1, 'opts': {'nargs': '1'}},
  \ {'type': 'command', 'name': 'GiPatch', 'sync': 1, 'opts': {}},
  \ ])
