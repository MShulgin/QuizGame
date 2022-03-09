package main

import (
	"bufio"
	"container/list"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type QuizGame struct {
	problems list.List
	timeout  uint32
}

func (game QuizGame) Start() {
	scoreCount := 0
	inputReader := bufio.NewReader(os.Stdin)
	gameCh := make(chan bool)
	go func() {
		for quizElem := game.problems.Front(); quizElem != nil; quizElem = quizElem.Next() {
			quizItem := quizElem.Value.(QuizProblem)
			fmt.Printf("%v = ?\n", quizItem.question)
			text, _ := inputReader.ReadString('\n')
			text = strings.Trim(text, "\n")
			if strings.Compare(text, "quit") == 0 {
				break
			} else if strings.Compare(text, quizItem.answer) == 0 {
				scoreCount += 1
			}
		}
		gameCh <- true
	}()

	select {
	case <-time.After(time.Duration(game.timeout) * time.Second):
		fmt.Println("Time is out")
	case <-gameCh:
		fmt.Println("No more questions")
	}

	fmt.Printf("Your score: %d\n", scoreCount)
}

type QuizProblem struct {
	question string
	answer   string
}

func parseProblems(fileName string) (list.List, error) {
	quizList := list.New()
	if quizFile, err := os.Open("quiz.csv"); err == nil {
		defer quizFile.Close()
		csvReader := csv.NewReader(io.Reader(quizFile))
		for {
			if csvRecord, err := csvReader.Read(); err == nil {
				item := QuizProblem{
					question: csvRecord[0],
					answer:   csvRecord[1],
				}
				quizList.PushBack(item)
			} else if err == io.EOF {
				break
			} else {
				return *quizList, err
			}
		}
	} else {
		return *quizList, err
	}

	return *quizList, nil
}

func main() {
	timeoutFlag := flag.Int("t", 30, "Timeout")
	flag.Parse()
	timeout := uint32(*timeoutFlag)
	if quizList, err := parseProblems("quiz.csv"); err == nil {
		game := QuizGame{problems: quizList, timeout: timeout}
		game.Start()
	} else {
		fmt.Println(err)
	}
}
