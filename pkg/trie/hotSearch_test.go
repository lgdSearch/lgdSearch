package trie

import (
	"io/ioutil"
	"lgdSearch/pkg/utils"
	"log"
	"testing"
	"time"
)

func TestHotSearch(t *testing.T) {
	filepath := "../data/HotSearch.txt"
	InitHotSearch(filepath)

	data := make([]queueNode, 0)
	utils.Read(&data, filepath)
	raw, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Println("error")
	}
	log.Println(len(data), len(raw))

	time.Sleep(time.Second * 5)

	GetHotSearch().showMapElement()
	GetHotSearch().showQueueElement()
	GetHotSearch().showArrayElement()

	SendText("你好")
	SendText("什么鬼")
	SendText("是的")
	SendText("你好")
	SendText("你好")
	SendText("世界")

	time.Sleep(time.Second)

	GetHotSearch().showMapElement()

	GetHotSearch().ReGetArray()
	GetHotSearch().showArrayElement()

	GetHotSearch().showQueueElement()

	time.Sleep(time.Second * 10)

	log.Println("Test Time Sub: ")
	for head := GetHotSearch().Queue().head; head != nil; head = head.Next {
		log.Println(time.Now().Sub(head.TimeMessage))
		log.Println(time.Now().Sub(head.TimeMessage).Hours())
		log.Println(time.Now().Sub(head.TimeMessage).Hours() > 24.)
	}

	node := GetHotSearch().Queue().Pop()
	GetHotSearch().searchMessage[node.Text]--
	log.Println(node.Text, node.TimeMessage)
	if GetHotSearch().searchMessage[node.Text] == 0 {
		delete(GetHotSearch().searchMessage, node.Text)
	}

	GetHotSearch().showQueueElement()
	GetHotSearch().showMapElement()
	GetHotSearch().ReGetArray()
	GetHotSearch().showArrayElement()

	SendText("中国")
	SendText("中国")
	SendText("外国")
	GetHotSearch().showQueueElement()
	GetHotSearch().showMapElement()
	GetHotSearch().ReGetArray()
	GetHotSearch().showArrayElement()

	node = GetHotSearch().Queue().Pop()
	GetHotSearch().searchMessage[node.Text]--
	if GetHotSearch().searchMessage[node.Text] == 0 {
		delete(GetHotSearch().searchMessage, node.Text)
	}
	log.Println(node.Text, node.TimeMessage)
	node = GetHotSearch().Queue().Pop()
	GetHotSearch().searchMessage[node.Text]--
	log.Println(node.Text, node.TimeMessage)
	if GetHotSearch().searchMessage[node.Text] == 0 {
		delete(GetHotSearch().searchMessage, node.Text)
	}

	GetHotSearch().showQueueElement()
	GetHotSearch().showMapElement()
	GetHotSearch().ReGetArray()
	GetHotSearch().showArrayElement()

	GetHotSearch().Flush(filepath)
}
