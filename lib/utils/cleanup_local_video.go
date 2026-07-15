package utils

import (
	"log"
	"os"
)

func CleanUpLocalVideo(fileName string) {
	err := os.Remove(fileName + ".mp4")
	if err != nil {
		log.Print("Failed to cleanup local video (original): ", err)
	}

	err = os.Remove(fileName + "_720.mp4")
	if err != nil {
		log.Print("Failed to cleanup local video (re-encoded): ", err)
	}
}
