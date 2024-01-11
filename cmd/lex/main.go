package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lex",
		Short: "A professional CLI text editor with a file explorer",
		Run:   runEditor,
	}

	app          *tview.Application
	editor       *tview.TextView
	fileExplorer *tview.TreeView
	selectedFile string
)

func runEditor(cmd *cobra.Command, args []string) {
	app = tview.NewApplication()

	// create a file explorer
	fileExplorer = tview.NewTreeView().
		SetRoot(buildFileTree("/", 0), "").
		SetCurrentNode(0).
		SetTopLevel(1).
		SetSelectedFunc(onFileSelected)

	editor = tview.NewTextView().
		SetText("Welcome to lex").
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetWordWrap(true)

	// Create a flex layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(fileExplorer, 0, 1, false).
		AddItem(editor, 0, 2, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildFileTree(root string, level int) *tview.TreeNode {
	node := tview.NewTreeNode(filepath.Base(root)).
		SetReference(root).
		SetSelectable(true)

	files, err := os.ReadDir(root)
	if err != nil {
		return node
	}

	for _, file := range files {
		if file.IsDir() {
			node.AddChild(buildFileTree(filepath.Join(root, file.Name()), level+1))
		}
	}

	return node
}

func onFileSelected(node *tview.TreeNode) {
	if node == nil {
		return
	}

	filePath, ok := node.GetReference().(string)
	if !ok {
		return
	}

	selectedFile = filePath
	loadFileContent(filePath)
}

func loadFileContent(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		editor.SetText(fmt.Sprintf("Error loading file: %v", err))
		return
	}

	editor.SetText(string(content))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
