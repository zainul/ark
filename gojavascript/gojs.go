package gojavascript

import (
	"strings"

	"golang.org/x/net/html"
)

// Get Eelement by ementType, attributeType and attributeName
// Ex: elementType = input
//     attributeType = class
//     attributeNames = ["input", "text"]
func GetElement(n *html.Node, elemType, attrType string, attrNames []string) *html.Node {
	return traverse(n, elemType, attrType, attrNames)
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkHtmlElement(n *html.Node, elemType, attrType string, attrNames []string) bool {
	if n.Type == html.ElementNode {
		if n.Data == elemType {
			s, ok := getAttribute(n, attrType)
			if ok {
				attrList := strings.Fields(s)
				//to handle multiple attribute name ex: <span class="place holder">
				expectedMatch := len(attrNames)
				matchedAttr := 0
				for i := range attrList {
					for j := range attrNames {
						if attrList[i] == attrNames[j] {
							matchedAttr++
						}
						if expectedMatch == matchedAttr {
							return true
						}
					}
				}

				//If no attribute matched
				return false
			}
		}
	}
	return false
}

func traverse(n *html.Node, elemType, attrType string, attrNames []string) *html.Node {
	if checkHtmlElement(n, elemType, attrType, attrNames) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, elemType, attrType, attrNames)
		if result != nil {
			return result
		}
	}

	return nil
}
