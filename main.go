package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("set json file")
	}

	jsonFile, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var chart Chart
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(byteValue, &chart); err != nil {
		panic(err)
	}

	mermaid, err := chart.ToMermaid()
	if err != nil {
		panic(err)
	}
	fmt.Println(mermaid)
}

type Node struct {
	ID     string `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Color  string `json:"color,omitempty"`
}

type Edge struct {
	ID       string `json:"id"`
	FromNode string `json:"fromNode"`
	FromSide string `json:"fromSide"`
	ToNode   string `json:"toNode"`
	ToSide   string `json:"toSide"`
	Label    string `json:"label,omitempty"`
}

type Chart struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

const mermaidTmp = `
flowchart TD
`

func (c *Chart) ToMermaid() (string, error) {
	result := mermaidTmp
	for _, edge := range c.Edges {
		fromNode, err := c.FindNode(edge.FromNode)
		if err != nil {
			return "", err
		}
		toNode, err := c.FindNode(edge.ToNode)
		if err != nil {
			return "", err
		}
		result += genMermaidLine(
			fromNode.ID,
			conv(fromNode.Text),
			conv(edge.Label),
			toNode.ID,
			conv(toNode.Text))
	}
	return result, nil
}

func (c *Chart) FindNode(nodeId string) (*Node, error) {
	for _, node := range c.Nodes {
		if node.ID == nodeId {
			return &node, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func genMermaidLine(fromNodeID, fromNodeText, edgeLabel, toNodeID, toNodeText string) string {
	noEdgeLabelTmp := "%s[%s] ---> %s[%s]\n"
	edgeLabelTmp := "%s[%s] ---> |%s| %s[%s]\n"
	if conv(edgeLabel) != "" {
		return fmt.Sprintf(
			edgeLabelTmp,
			fromNodeID,
			conv(fromNodeText),
			conv(edgeLabel),
			toNodeID,
			conv(toNodeText),
		)
	} else {
		return fmt.Sprintf(
			noEdgeLabelTmp,
			fromNodeID,
			conv(fromNodeText),
			toNodeID,
			conv(toNodeText),
		)
	}
}

func conv(str string) string {
	return strings.NewReplacer(
		"\r\n", "<br/>",
		"\r", "<br/>",
		"\n", "<br/>",
		"(", "",
		")", "",
		"、", ",",
	).Replace(str)
}
