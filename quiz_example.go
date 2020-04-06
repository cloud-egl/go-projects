package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "csv file for quiz setup")
	timeLimit := flag.Int("limit", 2, "timer limit for quiz")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Something went wrong while opening the file %s", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to load file content")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, prob := range problems {
		fmt.Printf("\n Problem #%d: %s = ", i+1, prob.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You answered %d questions correctly \n", correct)
			return
		case answer := <-answerCh:
			if answer == prob.a {
				fmt.Println("Correct !")
				correct += 1
			}
		}
	}
	fmt.Printf("You answered %d questions correctly", correct)

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
