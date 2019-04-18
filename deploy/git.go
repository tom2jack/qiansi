package deploy

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
)

func Git(cfg *DeployConfig) {
	var (
		remote_url = "ssh://git@gitee.com/273000727/go-test.git"
		path       = "D:\\go\\tools-client\\dist"
		//branch = "master"
		git_rsa = ""
	)
	Info(remote_url)
	//这里先读取文件，后置于远程
	pemBytes, _ := ioutil.ReadFile("git_rsa.pem")
	git_rsa = string(pemBytes)
	Info(git_rsa)
	signer, _ := ssh.ParsePrivateKey([]byte(git_rsa))
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}

	_, err := os.Stat(path)
	if err != nil {
		//克隆项目
		r, err := git.PlainClone(path, false, &git.CloneOptions{
			Auth: auth,
			URL:  remote_url,
		})
		CheckIfError(err)
		// ... retrieving the branch being pointed by HEAD
		ref, err := r.Head()
		CheckIfError(err)
		Info(ref.Hash().String())
		Info(ref.Name().String())
	} else {
		// We instance\iate a new repository targeting the given path (the .git folder)
		r, err := git.PlainOpen(path)
		CheckIfError(err)
		w, _ := r.Worktree()
		CheckIfError(err)

		//清理目录
		err = w.Clean(&git.CleanOptions{
			Dir: true,
		})
		CheckIfError(err)
		err = w.Checkout(&git.CheckoutOptions{
			Force: true,
		})
		CheckIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		Info("git pull origin")
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
		CheckIfError(err)
		ref, err := r.Head()
		CheckIfError(err)
		Info(ref.Hash().String())
		Info(ref.Name().String())
	}
}
