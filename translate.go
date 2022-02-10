package main

import (
	"bufio"
	"bytes"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/lxn/walk"
	"io"
	"strings"
)

var translateDict = &Trie{}
var regexpSpace = regexp.MustCompile(`(?<![()\[\]{}%'"A-Za-z]) (?![()\[\]{}%'"A-Za-z])`, regexp.None)

func init() {
	reader := bufio.NewReader(bytes.NewReader(transLateData))
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
			panic(err)
		}
		if len(line) > 0 {
			arr := strings.Split(string(line), "\t")
			var key, val string
			key = arr[0]
			if len(arr) >= 2 {
				val = arr[1]
			}
			if !translateDict.PutIfAbsent(key, val) {
				walk.MsgBox(nil, "警告", fmt.Sprint("出现重复数据：", string(line)), walk.MsgBoxIconWarning)
			}
		}
		if err == io.EOF {
			break
		}
	}
}

func translate(s string) string {
	s = translateDict.ReplaceAll(s)
	s, err := regexpSpace.Replace(s, "", -1, -1)
	if err != nil {
		panic(err)
	}
	return s
}

type trieNode struct {
	child map[rune]*trieNode
	value string
}

type Trie struct {
	root trieNode
}

func (t *Trie) PutIfAbsent(key, value string) bool {
	node := &t.root
	for _, c := range strings.ToLower(key) {
		var n *trieNode
		if node.child == nil {
			node.child = make(map[rune]*trieNode)
		}
		n = node.child[c]
		if n != nil {
			node = n
		} else {
			newNode := &trieNode{}
			node.child[c] = newNode
			node = newNode
		}
	}
	if len(node.value) > 0 {
		return false
	}
	node.value = value
	return true
}

func (t *Trie) getLongest(s string) (string, string) {
	var node, node2 *trieNode
	var key, key2 string
	node = &t.root
	r := []rune(strings.ToLower(s))
	for idx, c := range r {
		if node.child != nil {
			if n, ok := node.child[c]; ok {
				key += string(c)
				node = n
				if len(node.value) > 0 && (idx+1 >= len(s) || symbols[r[idx+1]]) {
					node2 = node
					key2 = key
				}
				continue
			}
		}
		break
	}
	if node2 != nil {
		return key2, node2.value
	}
	return "", ""
}

func (t *Trie) ReplaceAll(str string) string {
	s := []rune(str)
	var s2 []rune
	for len(s) > 0 {
		if !(len(s2) == 0 || symbols[s2[len(s2)-1]]) {
			s2 = append(s2, s[0])
			s = s[1:]
			continue
		}
		key, value := t.getLongest(string(s))
		if len(key) > 0 {
			s2 = append(s2, []rune(value)...)
			s = s[len([]rune(key)):]
		} else {
			s2 = append(s2, s[0])
			s = s[1:]
		}
	}
	return string(s2)
}

var symbols = map[rune]bool{
	' ':  true,
	'(':  true,
	')':  true,
	'[':  true,
	']':  true,
	'-':  true,
	'{':  true,
	'}':  true,
	'%':  true,
	'\'': true,
	'"':  true,
}
