package deploy

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os"
	"qiansi/common"
	"qiansi/models"
)

func Git(deploy *models.Deploy) {
	signer, _ := ssh.ParsePrivateKey([]byte(deploy.DeployKeys))
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}

	_, err := os.Stat(deploy.LocalPath)
	if err != nil {
		//克隆项目
		r, err := git.PlainClone(deploy.LocalPath, false, &git.CloneOptions{
			Auth: auth,
			URL:  deploy.RemoteUrl,
		})
		common.CheckIfError(err)
		// ... retrieving the branch being pointed by HEAD
		ref, err := r.Head()
		common.CheckIfError(err)
		common.Info(ref.Hash().String())
		common.Info(ref.Name().String())
	} else {
		// We instance\iate a new repository targeting the given path (the .git folder)
		r, err := git.PlainOpen(deploy.LocalPath)
		common.CheckIfError(err)
		w, _ := r.Worktree()
		common.CheckIfError(err)

		//清理目录
		err = w.Clean(&git.CleanOptions{
			Dir: true,
		})
		common.CheckIfError(err)
		err = w.Checkout(&git.CheckoutOptions{
			Force: true,
		})
		common.CheckIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		common.Info("git pull origin")
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
		common.CheckIfError(err)
		ref, err := r.Head()
		common.CheckIfError(err)
		common.Info(ref.Hash().String())
		common.Info(ref.Name().String())
	}
}
