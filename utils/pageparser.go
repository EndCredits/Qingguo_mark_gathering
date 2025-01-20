package utils

import (
	"encoding/json"
	"log"
	"qingguo/defs"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isCourses(text string) bool {
	courseRegex := regexp.MustCompile(`\[\w+\]\p{Han}*[a-zA-Z0-9]*`)
	return courseRegex.MatchString(text)
}

func ParsePage(content string) ([]byte, error) {
	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	// Extract necessary personal info
	var college, className, name, studentID string
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "院(系)/部") {
			college = strings.Split(text, "：")[1]
		} else if strings.Contains(text, "行政班级") {
			className = strings.Split(text, "：")[1]
		} else if strings.Contains(text, "学号") {
			studentID = strings.Split(text, "：")[1]
		} else if strings.Contains(text, "姓名") {
			name = strings.Split(text, "：")[1]
		}
	})

	// New student
	var newStudent defs.Student

	// Gather personal info
	newStudent.ID = studentID
	newStudent.Name = name
	newStudent.Institution = college
	newStudent.ClassName = className

	// Extract and gather scores
	newScore := make(defs.TotalScore)

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, tr *goquery.Selection) {
			var course, score string
			tr.Find("td").Each(func(k int, td *goquery.Selection) {
				// Course name is the first row
				// Score is the 7rd row
				if k == 1 {
					course = td.Text()
				} else if k == 7 {
					score = td.Text()
				}
			})
			if course != "" && score != "" && isCourses(course) {
				newScore[course] = score
			}
		})
	})

	newStudent.Scores = newScore

	jsonData, err := json.MarshalIndent(newStudent, "", "	")
	if err != nil {
		log.Fatal("Unable to convert student data to json: ", err)
		return nil, err
	}
	return jsonData, nil
}
