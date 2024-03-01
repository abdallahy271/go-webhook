package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func writeToTempFile(data interface{}) (string, error) {
	file, err := ioutil.TempFile("", "spec*.yml")
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func deleteTempFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func GetSpecDiff(content string) {

	// Write API response to a temporary file
	tempFilePath, err := writeToTempFile(content)
	if err != nil {
		log.Fatalf("Error writing API response to temp file: %v", err)
	}
	defer func() {
		err := deleteTempFile(tempFilePath)
		if err != nil {
			log.Fatalf("Error deleting temp file: %v", err)
		}
	}()

	// Get the difference between temp file and original file
	cmd := exec.Command("diff", "original.json", tempFilePath)
	diffOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running diff command: %v", err)
	}

	fmt.Println("Difference between original and API response:")
	fmt.Println(string(diffOutput))
}
