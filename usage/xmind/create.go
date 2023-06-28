package xmind

import "fmt"

const (
	SecondNode = 1
	ThirdNode  = 2
)

type XMindTrans struct {
	Root     string
	Children []*ItemNode
}

func NewXMindTrans(root string) *XMindTrans {
	return &XMindTrans{Root: root, Children: make([]*ItemNode, 0)}
}

func (x *XMindTrans) AppendChild(child *ItemNode) {
	x.Children = append(x.Children, child)
}

type ItemNode struct {
	Name     string
	Children []*ItemNode
}

func (i *ItemNode) AppendChild(child *ItemNode) {
	i.Children = append(i.Children, child)
}

func (x *XMindTrans) Output() (result string) {
	result = fmt.Sprintf("#%s", x.Root)

	for _, child := range x.Children {
		itemStr := child.output(SecondNode)
		result += "\n" + itemStr
	}
	return
}

func (i *ItemNode) output(level int) (result string) {
	defer func() {
		if i.Children != nil {
			for _, child := range i.Children {
				str := child.output(level + 1)
				if str != "" {
					result += "\n" + str
				}
			}
		}
	}()
	if level > ThirdNode {
		content := ""
		tableNum := level - ThirdNode - 1
		for i := 0; i < tableNum; i++ {
			content += "\t"
		}
		return content + "- " + i.Name
	} else if level == SecondNode {
		return fmt.Sprintf("##%s", i.Name)
	} else if level == ThirdNode {
		return fmt.Sprintf("###%s", i.Name)
	} else {
		return ""
	}
}

func CreateXMindFromMap(rootName string, data map[string]interface{}) *XMindTrans {
	root := NewXMindTrans(rootName)
	for name, itemData := range data {
		newNode := &ItemNode{
			Name:     name,
			Children: make([]*ItemNode, 0),
		}

		switch itemData.(type) {
		case struct{}:
		case map[string]interface{}:
			recursionNode(itemData.(map[string]interface{}), newNode)
		}
		root.AppendChild(newNode)
	}
	return root
}

func recursionNode(data map[string]interface{}, parent *ItemNode) {
	for name, itemData := range data {
		newNode := &ItemNode{
			Name:     name,
			Children: make([]*ItemNode, 0),
		}

		switch itemData.(type) {
		case struct{}:
		case map[string]interface{}:
			recursionNode(itemData.(map[string]interface{}), newNode)
		}
		parent.AppendChild(newNode)
	}
}
