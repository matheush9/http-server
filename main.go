package main

import (
	"fmt"
	"os"
)

func main() {
	serverConfig := NewServerConfig(60, "127.0.0.1", 34093)
	resourceMap := make(map[ResourcePath]ResourceResponse)
	resourceMap["/"] = ResourceResponse{Body: []byte("hi"), ResponseBodyType: TextHtml}
	resourceMap["/image"] = ResourceResponse{Body: LoadTestImage(), ResponseBodyType: ImagePng}

	Start(resourceMap, &serverConfig)
}

func LoadTestImage() []byte {
	img, err := os.ReadFile("./img/gopher.png")
	if err != nil {
		fmt.Println("error reading file: ", err)
		return nil
	}
	return img
}
