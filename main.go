package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Tree struct {
	word        string
	count       int
	left, right *Tree
}

func ValidCharts(ch rune) bool {
	if unicode.IsPunct(ch) {
		if strings.Contains(string(ch), "-") {
			return true
		} else {
			return false
		}
	} else if unicode.IsLetter(ch) {
		return true
	} else {
		return false
	}
}

func ReadWords(line string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	words := strings.Fields(line)
	for _, word := range words {
		word = strings.Map(func(ch rune) rune {
			if ValidCharts(ch) {
				return unicode.ToLower(ch)
			}
			return -1
		}, word)
		if word != "" {
			out <- word
		}
	}
}

func InsertShort(ptr **Tree, data string) {
	if *ptr == nil {
		*ptr = &Tree{
			word:  data,
			count: 1,
			left:  nil,
			right: nil,
		}
	} else if data == (*ptr).word {
		(*ptr).count++
	} else if data < (*ptr).word {
		InsertShort(&(*ptr).left, data)
	} else {
		InsertShort(&(*ptr).right, data)
	}
}

func SaveStorage(file *os.File, ptr *Tree) {
	if ptr != nil {
		SaveStorage(file, ptr.left)
		fmt.Fprintln(file, ptr.word, ptr.count)
		SaveStorage(file, ptr.right)
	}
}

func CleanStorage(ptr **Tree) {
	if *ptr != nil {
		CleanStorage(&(*ptr).left)
		CleanStorage(&(*ptr).right)
		*ptr = nil
	}
}

func main() {
	var wg sync.WaitGroup
	root := (*Tree)(nil)
	scanner := bufio.NewScanner(os.Stdin)
	startTime := time.Now()
	scanner.Split(bufio.ScanWords)
	wordsChan := make(chan string)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go ReadWords(line, wordsChan, &wg)
	}

	go func() {
		wg.Wait()
		close(wordsChan)
	}()

	for word := range wordsChan {
		InsertShort(&root, word)
	}
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()
	SaveStorage(file, root)
	CleanStorage(&root)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Скрипт выполнился за: %v\n", elapsedTime)
}
