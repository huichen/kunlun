package walker

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (dw *IndexWalker) getPubKeys() *ssh.PublicKeys {
	if dw.pubKeys != nil {
		return dw.pubKeys
	}

	// 获取 ssh public keys
	homeDir, err := os.UserHomeDir()
	if err == nil {
		absDirPath := homeDir + "/.ssh/id_rsa"
		dw.pubKeys, _ = ssh.NewPublicKeysFromFile("git", absDirPath, "")
	}

	return dw.pubKeys
}
