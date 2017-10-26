package command

import (
	"errors"
	"strconv"

	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

const noName = "NoName"

type Gilbert struct{}

// g:gilbert#allow_open_by_browser をチェックし、1ならブラウザで開く
func checkAndOpenGist(v *nvim.Nvim, url string) error {
	var isAllow int
	if v.Var("gilbert#allow_open_by_browser", &isAllow); isAllow == 1 {
		if err := util.Exec(v, "open "+url); err != nil {
			return err
		}
	}
	return nil
}

// g:gilbert#should_copy_url_to_clipboard をチェック
func checkAndCopyGistURL(v *nvim.Nvim, url string) error {
	var isAllow int
	if v.Var("gilbert#should_copy_url_to_clipboard", &isAllow); isAllow == 1 {
		if err := util.Exec(v, "echo '"+url+"' | pbcopy"); err != nil {
			return err
		}
	}
	return nil
}

// gist_idの情報を得る
func getGistIDFromBufferID(v *nvim.Nvim, buf nvim.Buffer) (string, error) {
	var bufferMap map[string]string
	v.Var("gilbert#buffer_and_gist_id_info", &bufferMap)
	res, ok := bufferMap[strconv.Itoa(int(buf))]
	if !ok {
		return "", errors.New("this buffer is not gist")
	}
	return res, nil
}

// 指定gist_idのバッファと、vimscript変数からすべて消す
func deleteBuffersOfGistID(v *nvim.Nvim, gistID string) error {
	var bufferMap map[string]string
	v.Var("gilbert#buffer_and_gist_id_info", &bufferMap)
	b := v.NewBatch()
	for key, value := range bufferMap {
		if value == gistID {
			b.Command("bw! " + key)
			b.Command("call remove(g:gilbert#buffer_and_gist_id_info, '" + key + "')")
		}
	}
	if err := b.Execute(); err != nil {
		return err
	}
	return nil
}

func getBufferIDsFromGistID(v *nvim.Nvim, gistID string) ([]nvim.Buffer, error) {
	bs := []nvim.Buffer{}

	var bufferMap map[string]string
	v.Var("gilbert#buffer_and_gist_id_info", &bufferMap)
	for k, v := range bufferMap {
		if v == gistID {
			b, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			bs = append(bs, nvim.Buffer(b))
		}
	}

	if len(bs) == 0 {
		return nil, errors.New("No such buffer(getBufferIDsFromGistID)")
	}
	return bs, nil
}

// gist_idの情報を追加
func setGistIDFromBufferID(v *nvim.Nvim, buf nvim.Buffer, gistID string) error {
	var bufferMap map[string]string
	if err := v.Var("gilbert#buffer_and_gist_id_info", &bufferMap); err != nil {
		return err
	}
	// bufferMap[int(buf)] = gistID
	b := strconv.Itoa(int(buf))
	// return v.SetVar("gilbert#buffer_and_gist_id_info", bufferMap)
	return v.Command("let g:gilbert#buffer_and_gist_id_info['" + b + "']='" + gistID + "'")
}

// Loadしたものかどうかを入れておく
func setGistIDToGiLoaded(v *nvim.Nvim, gistID string) error {
	var flagMap map[string]int
	if err := v.Var("gilbert#is_loaded_by_giload", &flagMap); err != nil {
		return err
	}
	return v.Command("let g:gilbert#is_loaded_by_giload['" + gistID + "']=1")
}

// mapにして返す
func getGistIDByGiLoaded(v *nvim.Nvim) (map[string]int, error) {
	var flagMap map[string]int
	if err := v.Var("gilbert#is_loaded_by_giload", &flagMap); err != nil {
		return map[string]int{}, err
	}
	return flagMap, nil
}

// delete history
func clearUndo(v *nvim.Nvim) error {
	return v.Command("call gilbert#clear_undo()")
}

// save
func saveFileInCurrentBuffer(v *nvim.Nvim) error {
	return v.Command("w!") // NOTE: !!!overwrite
}

// open
func openFileInCurrentBuffer(v *nvim.Nvim, file string) error {
	return v.Command("e " + file)
}
