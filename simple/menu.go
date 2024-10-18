package simple

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	id     uint
	Key    string
	Name   string
	Descr  string
	Action CommandAction
	MoveTo bool
}

type CommandAction func(m *Menu) error

func EmptyAction(m *Menu) error {
	return nil
}

func (c *Command) Print() {
	if c.Descr != "" {
		fmt.Printf("%s: %s (%s)\n", c.Key, c.Name, c.Descr)
	} else {
		fmt.Printf("%s: %s\n", c.Key, c.Name)
	}
}

type MenuConfig struct {
	Title            string
	BackKey          string
	UsingEscapeCodes bool
}

func getDefConfig() *MenuConfig {
	return &MenuConfig{
		Title:            "Default Menu Title",
		BackKey:          "<",
		UsingEscapeCodes: false,
	}
}

type commandNode struct {
	Val      *Command
	Children []*commandNode
	Parent   *commandNode
}

func createCommandNodeWoParent(cmd *Command) *commandNode {
	return &commandNode{
		Val:      cmd,
		Children: []*commandNode{},
	}
}
func createCommandNode(cmd *Command, parent *commandNode) *commandNode {
	return &commandNode{
		Val:      cmd,
		Children: []*commandNode{},
		Parent:   parent,
	}
}

func (cn *commandNode) AddChild(child *commandNode) {
	cn.Children = append(cn.Children, child)
}

type commandTree struct {
	Root    *commandNode
	Pointer *commandNode
}

func createCommandTree() *commandTree {
	root := createCommandNode(&Command{
		Name: "Root",
	}, nil)
	return &commandTree{
		Root:    root,
		Pointer: root,
	}
}

type Menu struct {
	Title   string
	BackKey string
	// Escape Code part
	//
	usingEscapeCodes bool
	lineCounter      uint
	//
	counterForIDs uint
	cmdTree       *commandTree
}

func CreateMenu(options ...SimpleMenuOption) *Menu {
	conf := getDefConfig()

	for _, opt := range options {
		opt(conf)
	}

	return &Menu{
		Title:            conf.Title,
		BackKey:          conf.BackKey,
		usingEscapeCodes: conf.UsingEscapeCodes,
		lineCounter:      0,
		counterForIDs:    1,
		cmdTree:          createCommandTree(),
	}
}

func (node *commandNode) AddCommand(key, name string, action CommandAction) *commandNode {
	cmd := &Command{
		Key:    key,
		Name:   name,
		Action: action,
		MoveTo: true,
	}
	newNode := createCommandNode(cmd, node)
	node.AddChild(newNode)

	return newNode
}

func (node *commandNode) AddExecCommand(key, name string, action CommandAction) *commandNode {
	cmd := &Command{
		Key:    key,
		Name:   name,
		Action: action,
		MoveTo: false,
	}
	newNode := createCommandNode(cmd, node)
	node.AddChild(newNode)

	return newNode
}

func (node *commandNode) AddGroupingCommand(key, name string) *commandNode {
	cmd := &Command{
		Key:    key,
		Name:   name,
		Action: EmptyAction,
		MoveTo: true,
	}
	newNode := createCommandNode(cmd, node)
	node.AddChild(newNode)

	return newNode
}

func (m *Menu) AddCommand(key, name string, action CommandAction) (*commandNode, error) {
	if m.cmdTree == nil {
		return nil, CommandTree_notInit
	}
	return m.cmdTree.Root.AddCommand(key, name, action), nil
}

func (m *Menu) clearLines() {
	for i := 0; i < int(m.lineCounter); i++ {
		fmt.Printf("\033[F\033[K")
	}
	m.lineCounter = 0
}

func (m *Menu) iteration() {
	fmt.Printf("%s\n", m.Title)
	path := ""
	for node := m.cmdTree.Pointer; node != nil; node = node.Parent {
		if node.Parent == nil {
			path += node.Val.Name
		} else {
			path += fmt.Sprintf("%s -> ", node.Val.Name)
		}
	}
	fmt.Printf("At %s\n", path)
	for _, node := range m.cmdTree.Pointer.Children {
		cmd := node.Val
		if cmd == nil {
			// TODO: Check it
			continue
		}
		cmd.Print()
	}
	if m.cmdTree.Pointer != m.cmdTree.Root {
		fmt.Printf("%s: Go back one layer\n", m.BackKey)
	}
	stdinReader := bufio.NewReader(os.Stdin)
	msg, err := stdinReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: ReadString: %v\n", err)
		return
	}
	msg = strings.TrimRight(msg, "\n")
	for _, node := range m.cmdTree.Pointer.Children {
		cmd := node.Val
		if cmd == nil {
			continue
		}
		if cmd.Key == msg {
			cmd.Action(m)
			m.cmdTree.Pointer = node
		}
	}
	if m.BackKey == msg {
		m.cmdTree.Pointer = m.cmdTree.Pointer.Parent
	}
}

func (m *Menu) Start() {
	for {
		m.iteration()
	}
}
