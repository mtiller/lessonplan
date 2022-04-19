package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	siren "github.com/xogeny/go-siren"
)

type Lesson struct {
	Title       string
	Explanation *string
	Model       *string
	Preamble    *string
	Report      *string
}

type IndexContents struct {
	Title    string   `json:"title"`
	Contents []string `json:"contents"`
}

type LessonPlan struct {
	Title   string
	Lessons []Lesson
}

func main() {
	lessons, err := ParseLessons(".")
	if err != nil {
		log.Fatalf(err.Error())
	}
	OutputLessons(lessons)
}

func OutputLessons(lessons LessonPlan) error {
	output := siren.NewSirenEntity()
	output.Title = lessons.Title
	for _, lesson := range lessons.Lessons {
		lent := siren.SirenEmbed{
			Title: lesson.Title,
		}
		lent.Title = lesson.Title
		lent.Class = []string{"lesson"}
		properties := map[string]interface{}{}
		if lesson.Explanation != nil {
			properties["explanation"] = *lesson.Explanation
		}
		if lesson.Model != nil {
			properties["model"] = *lesson.Model
		}
		if lesson.Preamble != nil {
			properties["preamble"] = *lesson.Preamble
		}
		if lesson.Report != nil {
			properties["report"] = *lesson.Report
		}
		lent.Properties = properties
		output.AddEmbed([]string{"item"}, lent)
	}
	contents, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling lessons: %s", err.Error())
	}
	fmt.Printf("%s\n", string(contents))
	return nil
}
func ParseLessons(dirname string) (LessonPlan, error) {
	ret := LessonPlan{}
	index, err := readIndex(path.Join(dirname, "index.json"))
	if index == nil {
		return LessonPlan{}, fmt.Errorf("Directory %s did not contain a valid index.json file", dirname)
	}
	if err != nil {
		return LessonPlan{}, err
	}
	ret.Title = index.Title
	for _, lessonName := range index.Contents {
		lesson, err := ParseLesson(path.Join(dirname, lessonName))
		if lesson == nil {
			fmt.Fprintf(os.Stderr, "No lesson found in directory %s (missing index.json?), skipping\n", lessonName)
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning, error searching for lesson in '%s', skipping: %s", lessonName, err.Error())
		} else {
			ret.Lessons = append(ret.Lessons, *lesson)
		}
	}
	return ret, nil
}

func ParseLesson(dirname string) (*Lesson, error) {
	ret := Lesson{}
	index, err := readIndex(path.Join(dirname, "index.json"))
	if index == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ret.Title = index.Title

	expFile := path.Join(dirname, "explanation.md")
	explanation, err := fileAsString(expFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", expFile)
	} else {
		ret.Explanation = explanation
	}

	modelFile := path.Join(dirname, "model.mo")
	model, err := fileAsString(modelFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", modelFile)
	} else {
		ret.Model = model
	}

	repFile := path.Join(dirname, "report.md")
	report, err := fileAsString(repFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", repFile)
	} else {
		ret.Report = report
	}

	preFile := path.Join(dirname, "preamble.md")
	preamble, err := fileAsString(preFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", preFile)
	} else {
		ret.Preamble = preamble
	}

	return &ret, nil
}

func readIndex(filename string) (*IndexContents, error) {
	contents, err := fileAsString(filename)
	if contents == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	index := IndexContents{}
	err = json.Unmarshal([]byte(*contents), &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

func fileAsString(filename string) (*string, error) {
	input, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to open file '%s': %s", filename, err.Error())
	}
	defer input.Close()
	content, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("error reading '%s': %s", filename, err.Error())
	}
	ret := string(content)
	return &ret, nil
}
