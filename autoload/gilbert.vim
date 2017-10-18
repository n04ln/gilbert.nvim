" flush undo history
function! gilbert#clear_undo() abort
    let old_undolevels = &undolevels
    setlocal undolevels=-1
    execute "normal! a \<BS>\<ESC>"
    let &l:undolevels = old_undolevels
endfunction
