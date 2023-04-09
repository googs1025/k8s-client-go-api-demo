package main

import "k8s-api-practice/restart_pod/pkg"

func main() {
	// 需要更新的副本数
	replicasNum := 4
	// dep name, namespace
	depName := "my-deployment"
	ns := "default"

	pods := pkg.GetPodsByDeployment(depName, ns)
	// 超过，直接更新全部
	if len(pods) < replicasNum {
		replicasNum = len(pods)
	}
	for i := 0; i < replicasNum; i++ {

		// pod原地升级
		// 一定要按照原本的顺序
		pkg.UpgradePodByImage(&pods[i], "nginx:1.19-alpine", "busybox")

		// pod原地重启
		//pkg.RestartPodByImage(&pods[i])
	}

}
