package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type QuizRecord struct {
	Question string
	Answer   int
}

func createQuiz(data [][]string) []QuizRecord {
	var quiz []QuizRecord
	for _, line := range data {

		var rec QuizRecord
		for j, field := range line {
			if j == 0 {
				rec.Question = field
			} else if j == 1 {
				rec.Answer,_ = strconv.Atoi(field)
			}
		}
		quiz = append(quiz, rec)

	}
	return quiz
}

func countdown(seconds int, currentTime chan int) {

	totalTime := seconds
	i:= 0

	for range time.Tick(1 * time.Second) {
		
		timeRemaining := totalTime - i
		currentTime <- timeRemaining

		if timeRemaining <= 0 {
			break
		}
		i++

		
	}
}



func main() {
	var correctAnswers int = 0
	reader := bufio.NewReader(os.Stdin)
	
	quizPath := flag.String( "path", "problems.csv", "path to csv file containg quiz")
	timeLimit := flag.Int( "limit", 30, "quiz time")
	flag.Parse()


	// open file
	f, err := os.Open(*quizPath)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	quiz := createQuiz(data)
	fmt.Printf("You have %d seconds to complete the quiz\n", *timeLimit)

	time := time.NewTimer(time.Duration(*timeLimit) * time.Second)


	quizLoop:
	for i, rec := range quiz {
		fmt.Printf("question %v / %v: \n", i + 1, len(quiz))
		fmt.Printf("%+v = ", rec.Question)
		answerChan := make(chan int)

		go func() {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		// parse input string into int
		input = strings.Trim(input, "\n")
		var answer int
		_, e := fmt.Sscan(input, &answer)
		
		if e != nil {
			fmt.Println("input must be integer", err)
			return
		}
		answerChan <- answer

	}()
	select {
	case answer := <- answerChan:
		if answer == rec.Answer {
			correctAnswers++
			
		}
	case <-time.C:
		fmt.Printf("\n\nTime's up\n")
		
		break quizLoop
		
	}



		
	}


	// user feedback
	fmt.Printf("correct answers: %v \n", correctAnswers)

}
