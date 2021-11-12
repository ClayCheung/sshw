package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	user "os/user"
	"path"

	"gopkg.in/yaml.v2"
)

var (
	localKubeConfPath = path.Join(homePath(), ".kube/config")
)

func homePath() string {
	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return curUser.HomeDir
}

type KubeConfig struct {
	ApiVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Clusters   []ClusterItem `yaml:"clusters"`
	Contexts   []ContextItem `yaml:"contexts"`
	Users      []UserItem    `yaml:"users"`
}

type ClusterItem struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}

type Cluster struct {
	Server                string `yaml:"server"`
	InsecureSkipTlsVerify bool   `yaml:"insecure-skip-tls-verify"`
}

type ContextItem struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type UserItem struct {
	User User   `yaml:"user"`
	Name string `yaml:"name"`
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

func loadKubeConf(yamlFile string) (*KubeConfig, error) {

	fd, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()

	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	data := KubeConfig{}
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func loadKubeConfFromString(kcStr []byte) (*KubeConfig, error) {
	data := KubeConfig{}
	err := yaml.Unmarshal(kcStr, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func SetNewCluster(kcStr []byte, server, contextName string) error {
	newCluster, err := loadKubeConfFromString(kcStr)
	if err != nil {
		return err
	}
	newCluster.Contexts = []ContextItem{
		{
			Context: Context{
				Cluster: contextName,
				User:    contextName,
			},
			Name: contextName,
		},
	}
	newCluster.Clusters = []ClusterItem{
		{
			Cluster: Cluster{
				Server:                fmt.Sprintf("https://%s:6443", server),
				InsecureSkipTlsVerify: true,
			},
			Name: contextName,
		},
	}
	newCluster.Users = newCluster.Users[:1]
	newCluster.Users[0].Name = contextName
	kc, err := loadKubeConf(localKubeConfPath)
	if err != nil {
		return err
	}
	err = setNewCluster(newCluster, kc)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(*kc)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(localKubeConfPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func setNewCluster(newCluster, kc *KubeConfig) error {
	var (
		needRemoveOld                bool
		rmContext, rmCluster, rmUser int
	)
	for i, context := range kc.Contexts {
		if context.Name == newCluster.Contexts[0].Name {
			needRemoveOld = true
			rmContext = i
			for i2, cluster := range kc.Clusters {
				if cluster.Name == context.Context.Cluster {
					rmCluster = i2
				}
			}
			for i2, user := range kc.Users {
				if user.Name == context.Context.User {
					rmUser = i2
				}
			}
		}
	}
	if needRemoveOld {
		if len(kc.Users) > rmUser+1 {
			kc.Users = append(kc.Users[:rmUser], kc.Users[rmUser+1:]...)
		} else {
			kc.Users = kc.Users[:rmUser]
		}

		if len(kc.Clusters) > rmCluster+1 {
			kc.Clusters = append(kc.Clusters[:rmCluster], kc.Clusters[rmCluster+1:]...)
		} else {
			kc.Clusters = kc.Clusters[:rmCluster]
		}

		if len(kc.Contexts) > rmContext+1 {
			kc.Contexts = append(kc.Contexts[:rmContext], kc.Contexts[rmContext+1:]...)
		} else {
			kc.Contexts = kc.Contexts[:rmContext]
		}
	}
	kc.Users = append(kc.Users, newCluster.Users[0])
	kc.Clusters = append(kc.Clusters, newCluster.Clusters[0])
	kc.Contexts = append(kc.Contexts, newCluster.Contexts[0])
	return nil
}
