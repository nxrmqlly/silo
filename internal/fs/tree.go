package fs

import (
	"os"
	"path/filepath"
)

type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*FileNode
	Parent   *FileNode

	// UI state
	Expanded bool
}

func NewFileNode(path string, parent *FileNode) *FileNode {
	info, _ := os.Stat(path)
	return &FileNode{
		Name:   filepath.Base(path),
		Path:   path,
		IsDir:  info.IsDir(),
		Parent: parent,
	}
}

func BuildFileTree(rootPath string) *FileNode {
	root := NewFileNode(rootPath, nil)
	buildChildren(root)
	return root
}

func buildChildren(node *FileNode) {
	if !node.IsDir {
		return
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return
	}

	for _, e := range entries {
		child := NewFileNode(filepath.Join(node.Path, e.Name()), node)
		node.Children = append(node.Children, child)
		buildChildren(child)
	}
}
