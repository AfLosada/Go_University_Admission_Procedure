package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	Chemistry   = "Chemistry"
	Physics     = "Physics"
	Mathematics = "Mathematics"
	Biotech     = "Biotech"
	Engineering = "Engineering"
)

func (a Applicant) chooseExamForDepartment(department string) float32 {
	var normalExam float32
	switch department {
	case Physics:
		normalExam = (float32(a.physics) + float32(a.math)) / 2
	case Chemistry:
		normalExam = float32(a.chemistry)
	case Mathematics:
		normalExam = float32(a.math)
	case Biotech:
		normalExam = (float32(a.chemistry) + float32(a.physics)) / 2
	case Engineering:
		normalExam = (float32(a.computerScience) + float32(a.math)) / 2
	}
	return float32(math.Max(float64(normalExam), float64(a.special)))
}

type Student struct {
	firstName string
	lastName  string
	exam      float32
}

type Applicant struct {
	firstName        string
	lastName         string
	physics          int
	chemistry        int
	math             int
	computerScience  int
	special          int
	firstPreference  string
	secondPreference string
	thirdPreference  string
}

func readApplicantFile() ([]Applicant, error) {
	file, err := os.Open("applicants.txt")
	if err != nil {
		fmt.Println("There was an error reading the file")
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	applicants := make([]Applicant, 0)
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		wordsOfLine := strings.Split(line, " ")
		firstName, lastName := wordsOfLine[0], wordsOfLine[1]
		examScores := wordsOfLine[2:7]
		preferences := wordsOfLine[7:]
		physics, _ := strconv.Atoi(examScores[0])
		chemistry, _ := strconv.Atoi(examScores[1])
		math, _ := strconv.Atoi(examScores[2])
		computerScience, _ := strconv.Atoi(examScores[3])
		special, _ := strconv.Atoi(examScores[4])
		applicant := Applicant{
			firstName:        firstName,
			lastName:         lastName,
			physics:          physics,
			chemistry:        chemistry,
			math:             math,
			computerScience:  computerScience,
			special:          special,
			firstPreference:  preferences[0],
			secondPreference: preferences[1],
			thirdPreference:  preferences[2],
		}
		applicants = append(applicants, applicant)
	}
	return applicants, err
}
