package deploy

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"qiansi/models"
)

func Git(deploy *models.Deploy) error {
	signer, err := ssh.ParsePrivateKey([]byte(deploy.DeployKeys))
	if err != nil {
		return err
	}
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
	repository, err := git.PlainOpen(deploy.LocalPath)
	if err != nil {
		if err.Error() == "repository does not exist" {
			LogPush("正在尝试初始化..%s", deploy.LocalPath)
			//克隆项目
			repository, err = git.PlainClone(deploy.LocalPath, false, &git.CloneOptions{
				Auth: auth,
				URL:  deploy.RemoteUrl,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	tree, err := repository.Worktree()
	if err != nil {
		return err
	}
	LogPush("成功抓取仓库分支图，正在执行当前工作区清理...")
	err = tree.Clean(&git.CleanOptions{Dir: true})
	if err != nil {
		return err
	}
	head, err := repository.Head()
	if err != nil {
		return err
	}
	LogPush("当前分支: %s, 部署分支: %s", head.Name().Short(), deploy.Branch)
	if head.Name().Short() != deploy.Branch {
		LogPush("分支数据不相同，正在尝试切换分支...")
		b, err := repository.Branches()
		if err != nil {
			return nil
		}
		needCreate := true
		_ = b.ForEach(func(reference *plumbing.Reference) error {
			if reference.Name().Short() == deploy.Branch {
				needCreate = false
			}
			return nil
		})
		if needCreate {
			LogPush("本地未找到部署分支，正在创建并切换...")
			err = repository.CreateBranch(&config.Branch{
				Name:   deploy.Branch,
				Remote: git.DefaultRemoteName,
				Merge:  plumbing.NewBranchReferenceName(deploy.Branch),
			})
			if err != nil {
				LogPush("创建分支报错: %v", err)
				return err
			}
		}
		err = tree.Checkout(&git.CheckoutOptions{
			// Create: needCreate,
			Force:  true,
			Branch: plumbing.NewBranchReferenceName(deploy.Branch),
		})
		if err != nil {
			return err
		}
		LogPush("分支切换完毕，正在刷新工作区信息数据...")
		head, err = repository.Head()
		if err != nil {
			return err
		}
	}
	LogPush("正在拉取远程最新数据...")
	err = tree.Pull(&git.PullOptions{
		// TODO:reference not found
		RemoteName:    git.DefaultRemoteName,
		ReferenceName: plumbing.NewRemoteReferenceName("heads", deploy.Branch),
		SingleBranch:  false,
		Auth:          auth,
	})
	if err != nil {
		if err.Error() != git.NoErrAlreadyUpToDate.Error() {
			return err
		}
		LogPush("没啥可更新的，当前hash:%s, 当前ref:%s", head.Hash().String(), head.Name().Short())
	} else {
		LogPush("数据更新完毕，当前hash:%s, 当前ref:%s", head.Hash().String(), head.Name().Short())
	}
	return nil
}
