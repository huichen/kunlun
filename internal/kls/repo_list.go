package kls

func (kls *KLS) redrawRepoList() {
	if kls.response == nil {
		return
	}

	kls.repoList.Clear()
	if len(kls.response.Repos) == 0 {
		kls.repoList.AddItem("没有", "", 0, nil)
		kls.app.Draw()
		return
	}

	kls.repoList.AddItem("[green]全部", "", 0, func() {
		kls.currentRepoID = -1
		kls.redrawFileList()
	})

	for id, repo := range kls.response.Repos {
		repoID := id
		name := repo.String()

		kls.repoList.AddItem("[green]"+name, "", 0, func() {
			kls.currentRepoID = repoID
			kls.redrawFileList()
		})
	}
}
