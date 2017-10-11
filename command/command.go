package command

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/NoahOrberg/gilbert/gist"
	"github.com/NoahOrberg/nvim-go-util/util"
	"github.com/neovim/go-client/nvim"
)

const noName = "NoName"
const loadingDuration = 50
const loadingCnt = 10

type Gilbert struct {
}

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

// gist_idの情報を追加
func setGistIDFromBufferID(v *nvim.Nvim, buf nvim.Buffer, gistID string) error {
	var bufferMap map[int]string
	if err := v.Var("gilbert#buffer_and_gist_id_info", &bufferMap); err != nil {
		return err
	}
	// bufferMap[int(buf)] = gistID
	b := strconv.Itoa(int(buf))
	// return v.SetVar("gilbert#buffer_and_gist_id_info", bufferMap)
	return v.Command("let g:gilbert#buffer_and_gist_id_info[" + b + "]='" + gistID + "'")
}

func (g *Gilbert) GilbertPatch(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	filename, err := v.BufferName(buf)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	temp := strings.Split(filename, "/")
	filename = temp[len(temp)-1]

	id, err := getGistIDFromBufferID(v, buf)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	gi := gist.Gist{
		Files: map[string]gist.File{
			filename: gist.File{
				Content: content,
			},
		},
	}

	res, err := gist.PatchGist(id, gi)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	util.Echom(v, res.HTMLURL)

	if err := checkAndCopyGistURL(v, res.HTMLURL); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	if err := checkAndOpenGist(v, res.HTMLURL); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	return nil
}

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	if err := util.NewBuffer(v); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	buf, err := v.CurrentBuffer()
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	kill := make(chan bool)

	go loadAnimation(v, buf, kill)

	// args[0] は、id,URLのどちらかを想定しているが、スラッシュで区切って最後の要素なのは変わらない
	temp := strings.Split(args[0], "/")
	id := temp[len(temp)-1]

	gi, err := gist.GetGist(id)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	if len(gi.Files) != 1 {
		err := errors.New("can open only one file gist")
		util.Echom(v, err.Error())
		return err
	}

	var filename string
	var strLines []string
	for key, value := range gi.Files {
		filename = key
		strLines = strings.Split(value.Content, "\n")
	}

	lines := make([][]byte, 0, len(strLines))
	for _, line := range strLines {
		lines = append(lines, []byte(line))
	}

	if err := setGistIDFromBufferID(v, buf, id); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	if err := v.SetBufferName(buf, filename); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	kill <- true

	if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	return nil
}

func loadAnimation(v *nvim.Nvim, buf nvim.Buffer, kill chan bool) {
	var cnt uint64
	d := byte('.')
	for {
		loading := []byte("Loading")
		waiting := []byte("Please wait")
		var i uint64
		for i = 0; i < cnt%loadingCnt; i++ {
			loading = append(loading, d)
			waiting = append(waiting, d)
		}
		if err := v.SetBufferLines(buf, 0, -1, true, [][]byte{
			loading,
			waiting,
		}); err != nil {
			util.Echom(v, err.Error())
			return
		}
		cnt++

		time.Sleep(loadingDuration * time.Millisecond)

		select {
		case <-kill:
			return
		default:
		}
	}
}

func (g *Gilbert) GilbertUpload(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	filepath, err := v.BufferName(buf)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	var url string
	var filename string
	url = "Missing"
	if filepath == "" {
		if len(args) > 0 {
			filename = args[0]
			filepath = args[0]
		} else {
			filepath = noName
			filename = noName
		}
	} else {
		temp := strings.Split(filepath, "/")
		filename = temp[len(temp)-1]
	}

	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	url, err = gist.PostToGistByContent("", filename, content)
	if err != nil {
		util.Echom(v, err.Error())
		return err
	}

	util.Echom(v, url)

	if err := checkAndCopyGistURL(v, url); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	if err := checkAndOpenGist(v, url); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	splittedURL := strings.Split(url, "/")
	id := splittedURL[len(splittedURL)-1]
	if err := setGistIDFromBufferID(v, buf, id); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	if err := v.SetBufferName(buf, filepath); err != nil {
		util.Echom(v, err.Error())
		return err
	}

	return nil
}
