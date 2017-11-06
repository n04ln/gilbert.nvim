package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

type Scoop int

const (
	GlobalScoop Scoop = iota + 1
	ScriptScoop
	LocalScoop
	BufferScoop
	WindowScoop
	BuiltinFuncScoop
)

func (s Scoop) String() string {
	switch s {
	case GlobalScoop:
		return "g:"
	case ScriptScoop:
		return "s:"
	case LocalScoop:
		return "l:"
	case BufferScoop:
		return "b:"
	case WindowScoop:
		return "w:"
	case BuiltinFuncScoop:
		return "v:"
	}
	return ""
}

const (
	noName = "NoName"
	// vimscript variable
	allowOpenByBrowser       = "gilbert#allow_open_by_browser"
	shouldCopyURLToClipBoard = "gilbert#should_copy_url_to_clipboard"
	bufferAndGistIDInfo      = "gilbert#buffer_and_gist_id_info"
	isLoadedByGiLoad         = "gilbert#is_loaded_by_giload"
	// the format of func, "****()"
	funcClearUndo = "gilbert#clear_undo()"
	//Error
)

var ErrUnexpectedType = func(t string) error {
	return fmt.Errorf("unexpected type : %s", t)
}

type Gilbert struct{}

// g:gilbert#allow_open_by_browser をチェックし、1ならブラウザで開く
func checkAndOpenGist(v *nvim.Nvim, url string) error {
	var isAllow int
	if v.Var(allowOpenByBrowser, &isAllow); isAllow == 1 {
		if err := util.Exec(v, "open "+url); err != nil {
			return err
		}
	}
	return nil
}

// g:gilbert#should_copy_url_to_clipboard をチェック
func checkAndCopyGistURL(v *nvim.Nvim, url string) error {
	var isAllow int
	if v.Var(shouldCopyURLToClipBoard, &isAllow); isAllow == 1 {
		if err := util.Exec(v, "echo '"+url+"' | pbcopy"); err != nil {
			return err
		}
	}
	return nil
}

// gist_idの情報を得る
func getGistIDFromBufferID(v *nvim.Nvim, buf nvim.Buffer) (string, error) {
	var bufferMap map[string]string
	v.Var(bufferAndGistIDInfo, &bufferMap)
	res, ok := bufferMap[strconv.Itoa(int(buf))]
	if !ok {
		return "", errors.New("this buffer is not gist")
	}
	return res, nil
}

// 指定gist_idのバッファと、vimscript変数からすべて消す
func deleteBuffersOfGistID(v *nvim.Nvim, gistID string) error {
	var bufferMap map[string]string
	v.Var(bufferAndGistIDInfo, &bufferMap)
	b := v.NewBatch()
	for key, value := range bufferMap {
		if value == gistID {
			b.Command("bw! " + key)
			removeKeyOfMapByBatch(b, GlobalScoop, bufferAndGistIDInfo, key)
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
	v.Var(bufferAndGistIDInfo, &bufferMap)
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
	if err := v.Var(bufferAndGistIDInfo, &bufferMap); err != nil {
		return err
	}
	// bufferMap[int(buf)] = gistID
	b := strconv.Itoa(int(buf))
	// return v.SetVar("gilbert#buffer_and_gist_id_info", bufferMap)
	return setValueOfMap(v, GlobalScoop, bufferAndGistIDInfo, b, gistID)
}

// Loadしたものかどうかを入れておく
func setGistIDToGiLoaded(v *nvim.Nvim, gistID string) error {
	var flagMap map[string]int
	if err := v.Var(isLoadedByGiLoad, &flagMap); err != nil {
		return err
	}
	return setValueOfMap(v, GlobalScoop, isLoadedByGiLoad, gistID, 1)
}

// mapにして返す
func getGistIDByGiLoaded(v *nvim.Nvim) (map[string]int, error) {
	var flagMap map[string]int
	if err := v.Var(isLoadedByGiLoad, &flagMap); err != nil {
		return map[string]int{}, err
	}
	return flagMap, nil
}

// setValueOfMap
// scoop => g, s, l, and so on
// TODO: error handling of unexpected scoop
func setValueOfMap(v *nvim.Nvim, s Scoop, variable, index string, value interface{}) error {
	var val string
	switch value.(type) {
	case int:
		val = fmt.Sprintf("%d", value.(int))
	case string:
		val = fmt.Sprintf("'%s'", value.(string))
	case bool:
		if value.(bool) {
			val = "1"
		} else {
			val = "0"
		}
	default:
		valType := fmt.Sprintf("%T", value)
		return ErrUnexpectedType(valType) // TODO: should return TYPE;;
	}
	return v.Command("let " + s.String() + variable + "['" + index + "']=" + val)
}

// removeKeyOfMap
// scoop => g, s, l, and so on
func removeKeyOfMapByBatch(b *nvim.Batch, s Scoop, variable, key string) {
	b.Command("call remove(" + s.String() + variable + ", '" + key + "')")
}

// delete history
func clearUndo(v *nvim.Nvim) error {
	return callFunc(v, funcClearUndo)
}

// call function
func callFunc(v *nvim.Nvim, f string) error {
	return v.Command("call " + f)
}

// save
func saveFileInCurrentBuffer(v *nvim.Nvim) error {
	return v.Command("w!") // NOTE: !!!overwrite
}

// save with filename
func saveFileInCurrentBufferWithName(v *nvim.Nvim, path string) error {
	return v.Command("w " + path) // NOTE: !!!overwrite
}

// open
func openFileInCurrentBuffer(v *nvim.Nvim, file string) error {
	return v.Command("e " + file)
}
