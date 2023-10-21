package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	var maxApplicantsPerDepartment int
	_, err := fmt.Scan(&maxApplicantsPerDepartment)
	if err != nil {
		return
	}
	applicantList, err := readApplicantFile()
	if err != nil {
		return
	}
	departments := admission(applicantList, maxApplicantsPerDepartment)
	printAdmissionResults(departments)
}

func saveDepartmentResults(department string, results []string) error {
	file, err := os.Create(department + ".txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, result := range results {
		_, err = fmt.Fprint(file, result)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func printAdmissionResults(departments map[string][]Student) {
	departmentNames := make([]string, 0)
	for key := range departments {
		departmentNames = append(departmentNames, key)
	}
	sort.SliceStable(departmentNames, func(i, j int) bool {
		return departmentNames[i] < departmentNames[j]
	})
	for _, department := range departmentNames {
		studentList := departments[department]
		fmt.Printf("%s\n", department)
		formattedStudents := sprintStudentList(studentList)
		fmt.Println()
		err := saveDepartmentResults(strings.ToLower(department), formattedStudents)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func sprintStudentList(students []Student) []string {
	formattedStudents := make([]string, 0)
	for _, student := range students {
		formattedStudent := fmt.Sprintf("%s %s %.1f\n", student.firstName, student.lastName, student.exam)
		formattedStudents = append(formattedStudents, formattedStudent)
	}
	return formattedStudents
}

func buildDepartmentStudentList() map[string][]Student {
	return map[string][]Student{
		"Mathematics": make([]Student, 0),
		"Physics":     make([]Student, 0),
		"Biotech":     make([]Student, 0),
		"Chemistry":   make([]Student, 0),
		"Engineering": make([]Student, 0),
	}
}

func admission(
	applicantList []Applicant,
	maxStudentsPerDepartment int) map[string][]Student {
	studentsPerDepartment := buildDepartmentStudentList()
	for priority := 0; priority < 3; priority++ {
		for department, students := range studentsPerDepartment {
			takenSlots := len(students)
			availableSlots := maxStudentsPerDepartment - takenSlots
			if availableSlots == 0 {
				continue
			}
			currentApplicants := findApplicantsGivenPriorityAndDepartment(applicantList, priority, department)
			currentApplicants = sortApplicantListByGPA(currentApplicants, department)
			if len(currentApplicants) == 0 {
				continue
			}
			var acceptedApplicants []Applicant
			if availableSlots > len(currentApplicants) {
				acceptedApplicants = currentApplicants
			} else {
				acceptedApplicants = currentApplicants[0:availableSlots]
			}
			if len(acceptedApplicants) == 0 {
				continue
			}
			studentList := make([]Student, 0)
			for _, applicant := range acceptedApplicants {
				studentList = append(studentList, applicant.parseApplicantToStudent(department))
			}
			studentsPerDepartment[department] = append(studentsPerDepartment[department], studentList...)
			applicantList = deleteApplicantsFromList(applicantList, acceptedApplicants)
		}
	}
	return orderDepartments(studentsPerDepartment)
}

func (a Applicant) parseApplicantToStudent(department string) Student {
	return Student{
		firstName: a.firstName,
		lastName:  a.lastName,
		exam:      a.chooseExamForDepartment(department),
	}
}

func orderDepartments(departments map[string][]Student) map[string][]Student {
	orderedDepartmentMap := make(map[string][]Student)
	for department := range departments {
		orderedDepartmentMap[department] = sortStudentList(departments[department])
	}
	return orderedDepartmentMap
}

func findPositionOfApplicant(applicantList []Applicant, applicant Applicant) int {
	for position, value := range applicantList {
		if value == applicant {
			return position
		}
	}
	return -1
}

func deleteApplicant(applicantList []Applicant, applicant Applicant) []Applicant {
	position := findPositionOfApplicant(applicantList, applicant)
	if position == -1 {
		return applicantList
	}
	applicantList[position] = applicantList[len(applicantList)-1]
	return applicantList[:len(applicantList)-1]
}

func deleteApplicantsFromList(applicantList []Applicant, applicantsToDelete []Applicant) []Applicant {
	for _, applicant := range applicantsToDelete {
		applicantList = deleteApplicant(applicantList, applicant)
	}
	return applicantList
}

func getPreferenceByPriority(applicant Applicant, priority int) string {
	switch priority {
	case 0:
		return applicant.firstPreference
	case 1:
		return applicant.secondPreference
	case 2:
		return applicant.thirdPreference
	default:
		return ""
	}
}

func findApplicantsGivenPriorityAndDepartment(applicantList []Applicant, priority int, department string) []Applicant {
	filteredApplicants := make([]Applicant, 0)
	for _, applicant := range applicantList {
		if getPreferenceByPriority(applicant, priority) == department {
			filteredApplicants = append(filteredApplicants, applicant)
		}
	}
	return filteredApplicants
}

func sortStudentList(students []Student) []Student {
	sort.SliceStable(students, func(i, j int) bool {
		if students[i].exam != students[j].exam {
			return students[i].exam > students[j].exam
		}
		return (students[i].firstName + students[i].lastName) < (students[j].firstName + students[j].lastName)
	})
	return students
}

func sortApplicantListByGPA(applicants []Applicant, department string) []Applicant {
	sort.SliceStable(applicants, func(i, j int) bool {
		if applicants[i].chooseExamForDeparment(department) != applicants[j].chooseExamForDeparment(department) {
			return applicants[i].chooseExamForDeparment(department) > applicants[j].chooseExamForDeparment(department)
		}
		return (applicants[i].firstName + applicants[i].lastName) < (applicants[j].firstName + applicants[j].lastName)
	})
	return applicants
}
