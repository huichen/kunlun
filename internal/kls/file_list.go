package kls

import (
	"path/filepath"
	"strings"
)

// 重新绘制左边的文件列表窗口
func (kls *KLS) redrawFileList() {
	if kls.response == nil || len(kls.response.Repos) == 0 {
		return
	}

	kls.fileList.Clear()
	if len(kls.response.Repos) == 0 {
		kls.fileList.AddItem("没有找到任何结果", "", 0, nil)
		kls.app.Draw()
		return
	}

	kls.fileList.AddItem("[green]全部", "", 0, func() {
		kls.currentFileID = -1
		kls.redrawFileContent()
	})

	for _, repo := range kls.response.Repos {
		for id, doc := range repo.Documents {
			fileID := id
			name := doc.Filename
			fields := strings.Split(name, "/")
			shortName := ""

			if len(fields) > 4 {
				shortName = strings.Join(fields[:3], "/") + "/~/" + filepath.Base(name)
			} else if len(fields) == 1 {
				shortName = name
			} else {
				shortName = strings.Join(fields[:len(fields)-1], "/") + "/" + filepath.Base(name)
			}

			kls.fileList.AddItem("[green]"+shortName, "", 0, func() {
				kls.currentFileID = fileID
				kls.redrawFileContent()
			})
		}
	}
}
