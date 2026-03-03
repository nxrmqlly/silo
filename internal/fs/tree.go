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
	Expanded bool // ui state
}

// convert raw data to `FileNode`s
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

// refresh the node (dir) when tree / files update
func RefreshNode(node *FileNode) {
	expanded := map[string]bool{}
	collectExpanded(node, expanded)
	node.Children = nil
	buildChildren(node)
	applyExpanded(node, expanded)
}

// recursively build a child tree for a dir
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

// recursively find and collect expanded nodes
func collectExpanded(node *FileNode, out map[string]bool) {
	if node.Expanded {
		out[node.Path] = true
	}
	for _, c := range node.Children {
		collectExpanded(c, out)
	}
}

// recursively apply expanded to all nodes
func applyExpanded(node *FileNode, expanded map[string]bool) {
	if expanded[node.Path] {
		node.Expanded = true
	}
	for _, c := range node.Children {
		applyExpanded(c, expanded)
	}
}
