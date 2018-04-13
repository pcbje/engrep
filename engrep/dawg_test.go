package engrep

import (
	"fmt"
	"testing"
)

func Find(node *Node, str string) []string {
	found := []string{}
	path := ""
	for _, char := range []rune(str) {
		node = node.Transition(char)
		if node == nil {
			break
		}
		path += string(char)

		if node.IsFinal() {
			found = append(found, path)
		}
	}
	return found
}

func Test1(t *testing.T) {
	lazyDawg := CreateDawg(2)

	lazyDawg.AddPattern("abbabbba")
	lazyDawg.AddPattern("peter christian")
	lazyDawg.AddPattern("petter chrisian abdx")

	root := lazyDawg.Iterator()

	fmt.Println(Find(root, "petr cristian"))
}

func Test2(t *testing.T) {
	lazyDawg := CreateDawg(2)

	lazyDawg.AddPattern("aaaaaabbaaa")
	lazyDawg.AddPattern("aabba")

	root := lazyDawg.Iterator()

	fmt.Println(Find(root, "aabba"))
}

func Test3(t *testing.T) {
	lazyDawg := CreateDawg(2)

	lazyDawg.AddPattern("0122")
	lazyDawg.AddPattern("1100")

	root := lazyDawg.Iterator()

	println(root.Transition('1').Transition('2'))

	//fmt.Println(Find(root, "122"))

	// "112020021","212020021"
}

func Test4(t *testing.T) {
	lazyDawg := CreateDawg(2)

	lazyDawg.AddPattern("011")
	lazyDawg.AddPattern("100")

	root := lazyDawg.Iterator()

	println(root.Transition('0').Transition('0'))

	//fmt.Println(Find(root, "122"))

	// "112020021","212020021"
}
