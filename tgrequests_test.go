package main

import (
	"fmt"
	"testing"
)

func TestRequestFileLink(t *testing.T) {
	fileId := "AgACAgIAAxkBAAMIY-ZU_yOyqZTBIPNdVag2eYIafOQAAl_DMRusxzlLvNNzQ4DNBKEBAAMCAAN5AAMuBA"
	link, err := getLinkToPhoto(fileId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(link)
}
