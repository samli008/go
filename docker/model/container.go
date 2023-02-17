package model

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type docker struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Image   string `json:"image"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
	Id      string `json:"id"`
}

type image struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Size string `json:"size"`
}

type net struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func ContainerList(cli *client.Client) []docker {
	lists := make([]docker, 0)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	list := docker{}
	for _, container := range containers {
		list.Name = container.Names[0][1:]
		list.Status = container.Status
		list.Image = container.Image
		list.Id = container.ID[:12]

		for i, v := range container.NetworkSettings.Networks {
			//fmt.Printf("%#v\n", v)
			if i != "network" {
				list.Network = i
				list.Ip = v.IPAddress
			}
		}
		lists = append(lists, list)
	}

	return lists
}

func RemoveContainer(containerID string, cli *client.Client) (string, error) {
	err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
	return containerID, err
}

func CreateContainer(cli *client.Client, containName string, net string, ip string, image string, hostname string) (string, error) {
	config := &container.Config{
		Image:    image,
		Hostname: hostname,
	}

	hostConfig := &container.HostConfig{
		Privileged: true,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:   mount.TypeBind,
		// 		Source: "/opt",
		// 		Target: "/opt",
		// 	},
		// },
	}

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"network": {
				NetworkID: net,
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: ip,
				},
			},
		},
	}
	resp, err1 := cli.ContainerCreate(context.Background(), config, hostConfig, netConfig, nil, containName)
	if err1 != nil {
		return resp.ID, err1
	}
	cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	return "created ok", err1
}

func ListImage(cli *client.Client) []image {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	lists := make([]image, 0)
	list := image{}
	for _, image := range images {
		list.Name = image.RepoTags[0]
		list.ID = image.ID[7:19]
		list.Size = strconv.Itoa(int(image.Size))[:3] + "MB"
		lists = append(lists, list)
	}
	return lists
}

func NetList(cli *client.Client) []net {
	nets, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}
	lists := make([]net, 0)
	list := net{}
	for _, net := range nets {
		if net.Driver == "macvlan" || net.Driver == "bridge" {
			list.Name = net.Name
			list.ID = net.ID[:12]
			lists = append(lists, list)
		}
	}
	return lists
}
