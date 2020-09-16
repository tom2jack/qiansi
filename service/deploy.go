package service

type deployService struct{}

func GetDeployService() *deployService {
	return &deployService{}
}
