package deploy

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os"
	"qiansi/models"
)

func Git(deploy *models.Deploy) error {
	signer, _ := ssh.ParsePrivateKey([]byte(deploy.DeployKeys))
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
	var (
		repository *git.Repository
		err        error
	)
	// 判断项目是否存在
	_, err = os.Stat(deploy.LocalPath)
	if err != nil {
		LogPush("正在克隆..")
		//克隆项目
		repository, err = git.PlainClone(deploy.LocalPath, false, &git.CloneOptions{
			Auth: auth,
			URL:  deploy.RemoteUrl,
		})
		if err != nil {
			return err
		}
	} else {
		// We instance\iate a new repository targeting the given path (the .git folder)
		repository, err = git.PlainOpen(deploy.LocalPath)
		if err != nil {
			return err
		}
		w, err := repository.Worktree()
		if err != nil {
			return err
		}
		//清理目录
		err = w.Clean(&git.CleanOptions{
			Dir: true,
		})
		if err != nil {
			return err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Force: true,
		})
		if err != nil {
			return err
		}
		err = w.Pull(&git.PullOptions{
			Auth: auth,
		})
		fmt.Printf("dellllll::::: %v", err)
		if err != nil {
			return err
		}
	}
	ref, err := repository.Head()
	if err != nil {
		return err
	}
	LogPush(ref.Hash().String())
	LogPush(ref.Name().String())
	return nil
}
