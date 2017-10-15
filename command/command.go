package command

import (
	"errors"
	"strconv"
	"strings"

	"github.com/NoahOrberg/gilbert/gist"
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

func (g *Gilbert) GilbertPatch(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	gistID, err := getGistIDFromBufferID(v, buf)
	if err != nil {
		return err
	}

	bufs, err := getBufferIDsFromGistID(v, gistID)
	if err != nil {
		return err
	}

	files := map[string]gist.File{}

	for _, buf := range bufs {
		filename, err := v.BufferName(buf)
		if err != nil {
			return err
		}

		if filename == "" {
			filename = noName
		}

		temp := strings.Split(filename, "/")
		filename = temp[len(temp)-1]

		lines, err := v.BufferLines(buf, 0, -1, true)
		if err != nil {
			return err
		}

		var content string
		for i, c := range lines {
			content += string(c)
			if i < len(lines)-1 {
				content += "\n"
			}
		}

		files[filename] = gist.File{
			Content: content,
		}
	}

	gi := gist.Gist{
		Files: files,
	}

	res, err := gist.PatchGist(gistID, gi)
	if err != nil {
		return err
	}

	util.Echom(v, res.HTMLURL)

	if err := checkAndOpenGist(v, res.HTMLURL); err != nil {
		return err
	}

	if err := checkAndCopyGistURL(v, res.HTMLURL); err != nil {
		return err
	}

	if err := deleteBuffersOfGistID(v, gistID); err != nil {
		return err
	}

	return nil
}

func (g *Gilbert) GilbertLoad(v *nvim.Nvim, args []string) error {
	temp := strings.Split(args[0], "/")
	id := temp[len(temp)-1]

	gi, err := gist.GetGist(id)
	if err != nil {
		return err
	}

	for filename, file := range gi.Files {
		if err := util.NewBuffer(v); err != nil {
			return err
		}

		buf, err := v.CurrentBuffer()
		if err != nil {
			return err
		}

		var strLines []string
		strLines = strings.Split(file.Content, "\n")

		lines := make([][]byte, 0, len(strLines))
		for _, line := range strLines {
			lines = append(lines, []byte(line))
		}

		if err := setGistIDFromBufferID(v, buf, id); err != nil {
			return err
		}

		if err := v.SetBufferName(buf, id+"/"+filename); err != nil {
			return err
		}

		if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gilbert) GilbertUpload(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	filepath, err := v.BufferName(buf)
	if err != nil {
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
		return err
	}

	util.Echom(v, url)

	if err := checkAndOpenGist(v, url); err != nil {
		return err
	}

	if err := checkAndCopyGistURL(v, url); err != nil {
		return err
	}

	splittedURL := strings.Split(url, "/")
	id := splittedURL[len(splittedURL)-1]
	if err := setGistIDFromBufferID(v, buf, id); err != nil {
		return err
	}

	return nil
}
