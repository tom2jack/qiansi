package deploy

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os"
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
		if err != nil {
			LogPush(err.Error())
			return
		}
		// ... retrieving the branch being pointed by HEAD
		ref, err := r.Head()
		if err != nil {
			LogPush(err.Error())
			return
		}
		LogPush(ref.Hash().String())
		LogPush(ref.Name().String())
	} else {
		// We instance\iate a new repository targeting the given path (the .git folder)
		r, err := git.PlainOpen(deploy.LocalPath)
		if err != nil {
			LogPush(err.Error())
			return
		}
		w, _ := r.Worktree()
		if err != nil {
			LogPush(err.Error())
			return
		}

		//清理目录
		err = w.Clean(&git.CleanOptions{
			Dir: true,
		})
		if err != nil {
			LogPush(err.Error())
			return
		}
		err = w.Checkout(&git.CheckoutOptions{
			Force: true,
		})
		if err != nil {
			LogPush(err.Error())
			return
		}

		// Pull the latest changes from the origin remote and merge into the current branch
		LogPush("git pull origin")
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
		if err != nil {
			LogPush(err.Error())
			return
		}
		ref, err := r.Head()
		if err != nil {
			LogPush(err.Error())
			return
		}
		LogPush(ref.Hash().String())
		LogPush(ref.Name().String())
	}
}
