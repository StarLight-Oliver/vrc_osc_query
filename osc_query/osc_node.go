package osc_query

import (
	"path/filepath"
)

type OscNode struct {
	Description string             `json:"DESCRIPTION"`
	FullPath    string             `json:"FULL_PATH"`
	Contents    map[string]OscNode `json:"CONTENTS"`
	Type        string             `json:"TYPE"`
	HostInfo    HostInfo           `json:"HOST_INFO"`
	Access      int                `json:"ACCESS"`
}

type HostInfo struct {
	Name       string          `json:"NAME"`
	Extensions map[string]bool `json:"EXTENSIONS"`
	Ip         string          `json:"OSC_IP"`
	Port       int             `json:"OSC_PORT"`
	Transport  string          `json:"TRANSPORT"`
}

func NewOscNodeTree(applicationName, host string, port int) *OscNode {

	node := OscNode{
		Description: "root node",
		FullPath:    "/",
		Contents:    make(map[string]OscNode),
		HostInfo: HostInfo{
			Name:       applicationName,
			Extensions: make(map[string]bool),
			Ip:         host,
			Port:       port,
			Transport:  "UDP",
		},
		Access: 0,
	}

	node.HostInfo.Extensions["ACCESS"] = true

	return &node
}

const (
	OscTypeInt   = 1
	OscTypeFloat = 2
	OscTypeBool  = 3
)

func (node *OscNode) AddChild(fullpath string, value int, desc string) *OscNode {
	folderPath, fileName := filepath.Split(fullpath)

	if folderPath[0:len(folderPath)-1] == node.FullPath {

		osc_node := OscNode{
			Description: desc,
			FullPath:    fullpath,
			Contents:    make(map[string]OscNode),
		}

		if value == OscTypeInt {
			osc_node.Type = "i"
		} else if value == OscTypeFloat {
			osc_node.Type = "f"
		} else if value == OscTypeBool {
			osc_node.Type = "T"
		}
		osc_node.Access = 3

		node.Contents[fileName] = osc_node
		return &osc_node
	} else {
		folderName := fullpath[len(node.FullPath):]

		for i, c := range folderName {
			if c == '/' {
				if i == 0 {
					continue
				}
				folderName = folderName[0:i]
				break
			}
		}

		if _, ok := node.Contents[folderName]; !ok {
			aaa := OscNode{
				Description: "folder",
				FullPath:    node.FullPath + folderName,
				Contents:    make(map[string]OscNode),
				Access:      0,
			}

			node.Contents[folderName] = aaa
		}
		sub_node := node.Contents[folderName]

		return sub_node.AddChild(fullpath, value, desc)
	}
}
