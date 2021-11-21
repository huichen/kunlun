package kls

import "github.com/gdamore/tcell/v2"

func (kls *KLS) SetInputFieldKeys() {
	kls.inputField.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyUp:
				return nil
			case tcell.KeyDown:
				kls.app.SetFocus(kls.repoList)
				return nil
			case tcell.KeyTab:
				kls.app.SetFocus(kls.repoList)
				return nil
			}
			return event
		})
}

func (kls *KLS) SetRepoListKeys() {
	kls.repoList.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case '/':
				kls.app.SetFocus(kls.inputField)
				return nil
			}
			switch event.Key() {
			case tcell.KeyTab:
				kls.app.SetFocus(kls.fileList)
				return nil
			case tcell.KeyLeft:
				return nil
			case tcell.KeyRight:
				kls.app.SetFocus(kls.fileContent)
				return nil
			case tcell.KeyDown:
				newItemID := (kls.repoList.GetCurrentItem() + 1) % kls.repoList.GetItemCount()
				kls.repoList.SetCurrentItem(newItemID)
				kls.currentRepoID = newItemID - 1
				kls.redrawFileList()
				kls.redrawFileContent()
				return nil
			case tcell.KeyUp:
				newItemID := (kls.repoList.GetCurrentItem() + kls.repoList.GetItemCount() - 1) % kls.repoList.GetItemCount()
				kls.repoList.SetCurrentItem(newItemID)
				kls.currentRepoID = newItemID - 1
				kls.redrawFileList()
				kls.redrawFileContent()
				return nil
			}
			return event
		})
}

func (kls *KLS) SetFileListKeys() {
	kls.fileList.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case '/':
				kls.app.SetFocus(kls.inputField)
				return nil
			}
			switch event.Key() {
			case tcell.KeyTab:
				kls.app.SetFocus(kls.fileContent)
				return nil
			case tcell.KeyLeft:
				return nil
			case tcell.KeyRight:
				kls.app.SetFocus(kls.fileContent)
				return nil
			case tcell.KeyDown:
				newItemID := (kls.fileList.GetCurrentItem() + 1) % kls.fileList.GetItemCount()
				kls.fileList.SetCurrentItem(newItemID)
				kls.currentFileID = newItemID - 1
				kls.redrawFileContent()
				return nil
			case tcell.KeyUp:
				newItemID := (kls.fileList.GetCurrentItem() + kls.fileList.GetItemCount() - 1) % kls.fileList.GetItemCount()
				kls.fileList.SetCurrentItem(newItemID)
				kls.currentFileID = newItemID - 1
				kls.redrawFileContent()
				return nil
			}
			return event
		})
}

func (kls *KLS) SetFileContentKey() {
	kls.fileContent.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case '/':
				kls.app.SetFocus(kls.inputField)
				return nil
			}
			switch event.Key() {
			case tcell.KeyTab:
				kls.app.SetFocus(kls.inputField)
				return nil
			case tcell.KeyLeft:
				kls.app.SetFocus(kls.fileList)
				return nil
			}
			return event
		})
}
