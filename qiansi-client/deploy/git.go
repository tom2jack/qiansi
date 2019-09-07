package deploy

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"qiansi/common/dto"
)

func Git(deploy *dto.DeployDTO) error {
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
			LogPush(deploy, "正在尝试初始化..%s", deploy.LocalPath)
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
	head, err := repository.Head()
	if err != nil {
		return err
	}
	LogPush(deploy, "成功抓取仓库分支图，正在执行当前工作区清理...")
	err = tree.Clean(&git.CleanOptions{Dir: true})
	if err != nil {
		return err
	}
	err = tree.Checkout(&git.CheckoutOptions{
		Force:  true,
		Branch: head.Name(),
	})
	if err != nil {
		return err
	}
	// 获取远程分支信息
	remote, err := repository.Storer.Reference(plumbing.NewRemoteReferenceName(git.DefaultRemoteName, deploy.Branch))
	LogPush(deploy, "当前分支: %s(%s), 部署分支: %s, 部署版本：%s", head.Name().Short(), head.Hash(), deploy.Branch, remote.Hash())
	if head.Name().Short() != deploy.Branch {
		LogPush(deploy, "分支数据不相同，正在尝试切换分支...")
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
			LogPush(deploy, "本地未找到部署分支，正在创建并切换...")
			err = repository.CreateBranch(&config.Branch{
				Name:   deploy.Branch,
				Remote: git.DefaultRemoteName,
				Merge:  plumbing.NewBranchReferenceName(deploy.Branch),
			})
			if err != nil {
				LogPush(deploy, "创建分支报错: %v", err)
				return err
			}
			err = tree.Checkout(&git.CheckoutOptions{
				Create: true,
				Force:  true,
				Hash:   remote.Hash(),
				Branch: plumbing.NewBranchReferenceName(deploy.Branch),
			})
			if err != nil {
				return err
			}
		} else {
			LogPush(deploy, "本地寻址成功，正在切换分支到: %s, Hash:%s", deploy.Branch, remote.Hash())
			err = tree.Checkout(&git.CheckoutOptions{
				Create: false,
				Force:  true,
				Branch: plumbing.NewBranchReferenceName(deploy.Branch),
			})
			if err != nil {
				return err
			}
		}
		LogPush(deploy, "分支切换完毕，正在刷新工作区信息数据...")
		head, err = repository.Head()
		if err != nil {
			return err
		}
	}
	LogPush(deploy, "正在拉取远程最新数据...")
	err = tree.Pull(&git.PullOptions{
		// Name of the remote to be pulled. If empty, uses the default.
		RemoteName: git.DefaultRemoteName,
		// Remote branch to clone. If empty, uses HEAD.
		ReferenceName: plumbing.NewBranchReferenceName(deploy.Branch),
		// Fetch only ReferenceName if true.
		SingleBranch: false,
		Auth:         auth,
	})
	if err != nil {
		if err.Error() != git.NoErrAlreadyUpToDate.Error() {
			return err
		}
		LogPush(deploy, "没啥可更新的，当前版本:%s, 当前分支:%s", head.Hash().String(), head.Name().Short())
	} else {
		LogPush(deploy, "数据更新完毕，当前版本:%s, 当前分支:%s", head.Hash().String(), head.Name().Short())
	}
	return nil
}
