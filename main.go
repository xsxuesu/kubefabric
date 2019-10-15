package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"kubefabric/kubeutils"
)

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "action",
			Value: "show",
			Usage: "action for kubernetes",
		},
	}

	app.Action = func(c *cli.Context) error {
		name := "show"
		if c.NArg() > 0 {
			name = c.Args().Get(0)
		}
		switch name {
		case "show":
			ShowKube()
			break
		case "create":
			CreateKube()
			ListKube()
			break
		case "update":
			UpdateKube()
			ListKube()
			break
		case "delete":
			DeleteKube()
			ListKube()
			break
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowKube(){
	kubeClient := kubeutils.InitClient()
	podList := kubeClient.GetPod("")

	for _,pod := range podList{
		fmt.Printf("pod name :%s pod type:%s namespace:%s\n",pod.Name,pod.Kind,pod.Namespace)
	}

	nodeList := kubeClient.GetNode()

	fmt.Println("===========================================================")
	for _,node := range nodeList{
		fmt.Printf("node name :%s node type:%s \n",node.Name,node.Kind)
	}
}
func CreateKube(){
	kubeClient := kubeutils.InitClient()
	deployment,err := kubeClient.CreateDeployment()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(deployment.Namespace,deployment.Name)
}
func UpdateKube(){
	kubeClient := kubeutils.InitClient()
	err := kubeClient.UpdateDeployment("demo-deployment")
	if err != nil {
		fmt.Println(err.Error())
	}
}
func ListKube(){
	kubeClient := kubeutils.InitClient()
	listDeployment,err := kubeClient.ListDeployment()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _,deploy:=range listDeployment{
		fmt.Println(deploy.Namespace,deploy.Name)
	}
}
func DeleteKube(){
	kubeClient := kubeutils.InitClient()
	err := kubeClient.DeleteDeployment("demo-deployment")
	if err != nil {
		fmt.Println(err.Error())
	}
}