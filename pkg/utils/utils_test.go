package utils

import (
	"log"
	"testing"
)

func TestUtils(t *testing.T) {
	wordMap := make(map[uint32]int)
	encoder := Encoder(wordMap)
	log.Println(len(encoder))

	wordMap[0] = 1
	wordMap[2] = 3
	wordMap[5]++
	encoder = Encoder(wordMap)
	log.Println(len(encoder))

	wordMap[4]++
	encoder = Encoder(&wordMap)
	log.Println(len(encoder))

}

func TestIndex(t *testing.T) {
	var data uint32 = 780867
	encoder := Encoder(&data)
	Decoder(encoder, &data)
	log.Println(data)

	Write(&data, "../../../data/dataIndex.txt")
	log.Println(data)
}
