package simple

import (
	"bufio"
	"fmt"
	"log"
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
	Log *BufferedLogger
	//
	counterForIDs uint
	cmdTree       *commandTree
	//
	interrupt    chan func(m *Menu)
	resumeSignal chan struct{}
	inputChan    chan string
	errorChan    chan error
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
		Log:              &BufferedLogger{},
		lineCounter:      0,
		counterForIDs:    1,
		cmdTree:          createCommandTree(),
		interrupt:        make(chan func(m *Menu)),
		resumeSignal:     make(chan struct{}),
		inputChan:        make(chan string),
		errorChan:        make(chan error),
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

func (m *Menu) AddExecCommand(key, name string, action CommandAction) (*commandNode, error) {
	if m.cmdTree == nil {
		return nil, CommandTree_notInit
	}
	return m.cmdTree.Root.AddExecCommand(key, name, action), nil
}

func (m *Menu) AddGroupingCommand(key, name string) (*commandNode, error) {
	if m.cmdTree == nil {
		return nil, CommandTree_notInit
	}
	return m.cmdTree.Root.AddGroupingCommand(key, name), nil
}

func (m *Menu) clearLines() {
	for i := 0; i < int(m.lineCounter); i++ {
		fmt.Printf("\033[F\033[K")
	}
	m.lineCounter = 0
}

func (m *Menu) handleInput(input string) {
	var preHandler func()
	var action CommandAction
	var afterHandler func()
	if m.usingEscapeCodes {
		preHandler = func() {
			m.clearLines()
		}
	}
	for _, node := range m.cmdTree.Pointer.Children {
		cmd := node.Val
		if cmd == nil {
			continue
		}
		if cmd.Key == input {
			action = cmd.Action
			if cmd.MoveTo {
				afterHandler = func() {
					m.cmdTree.Pointer = node
				}
			}
			break
		}
	}
	if m.cmdTree.Pointer != m.cmdTree.Root && m.BackKey == input {
		afterHandler = func() {
			m.cmdTree.Pointer = m.cmdTree.Pointer.Parent
		}
	}
	if preHandler != nil {
		preHandler()
	}
	if action != nil {
		action(m)
	}
	if afterHandler != nil {
		afterHandler()
	}
}

func (m *Menu) printMenu() {
	fmt.Printf("%s\n", m.Title)
	m.lineCounter++
	path := ""
	for node := m.cmdTree.Pointer; node != nil; node = node.Parent {
		if node.Parent == nil {
			path += node.Val.Name
		} else {
			path += fmt.Sprintf("%s -> ", node.Val.Name)
		}
	}
	fmt.Printf("At %s\n", path)
	m.lineCounter++
	for _, node := range m.cmdTree.Pointer.Children {
		cmd := node.Val
		if cmd == nil {
			// TODO: Check it
			continue
		}
		cmd.Print()
		m.lineCounter++
	}
	if m.cmdTree.Pointer != m.cmdTree.Root {
		fmt.Printf("%s: Go back one layer\n", m.BackKey)
		m.lineCounter++
	}
}

func (m *Menu) readInput() {
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		msg, err := stdinReader.ReadString('\n')
		if err != nil {
			// fmt.Printf("Error: ReadString: %v\n", err)
			m.errorChan <- err
			return
		}
		m.lineCounter++
		m.inputChan <- msg
	}
}

func (m *Menu) GetInput(prompt string) string {
	nlCount := strings.Count(prompt, "\n")
	m.lineCounter += uint(nlCount)
	stdinReader := bufio.NewReader(os.Stdin)
	rawMsg, err := stdinReader.ReadString('\n')
	if err != nil {
		// fmt.Printf("Error: ReadString: %v\n", err)
		m.errorChan <- err
		return ""
	}
	m.lineCounter++
	msg := strings.TrimRight(rawMsg, "\n")
	return msg
}

func (m *Menu) iteration() {
	m.Log.Flush()
	m.printMenu()

	// for {
	select {
	case msg := <-m.inputChan:
		msg = strings.TrimRight(msg, "\n")
		m.handleInput(msg)
	case err := <-m.errorChan:
		m.Log.Logf("err: %v", err)
	case f := <-m.interrupt:
		if m.usingEscapeCodes {
			m.clearLines()
		}
		log.Printf("Interrupt")
		f(m)
	}
	// }
}

func (m *Menu) SendInterrupt(f func(m *Menu)) {
	m.interrupt <- f
}

func (m *Menu) Start() {
	go m.readInput()
	for {
		m.iteration()
	}
}
