package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	siren "github.com/xogeny/go-siren"
)

func initRoot() *cobra.Command {
	var rootCmd = cobra.Command{
		Use:   "lessonplan",
		Short: "Create Modelica Playground lesson plans",
		Long:  `Create lesson plans for the Modelica Playground application found at http://playground.modelica.university`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}
			out, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			lessons, err := ParseLessons(dir)
			if err != nil {
				return err
			}
			return OutputLessons(lessons, out)
		},
	}
	rootCmd.PersistentFlags().StringP("dir", "d", ".", "Directory that contains the lesson plan files")
	rootCmd.PersistentFlags().StringP("output", "o", "lessonplan.json", "Name of file to generate")
	return &rootCmd
}

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
	rootCmd := initRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func OutputLessons(lessons LessonPlan, outputFile string) error {
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
	return ioutil.WriteFile(outputFile, contents, fs.FileMode(0777))
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
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", expFile, err.Error())
	} else {
		ret.Explanation = explanation
	}

	modelFile := path.Join(dirname, "model.mo")
	model, err := fileAsString(modelFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", modelFile, err.Error())
	} else {
		ret.Model = model
	}

	repFile := path.Join(dirname, "report.md")
	report, err := fileAsString(repFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", repFile, err.Error())
	} else {
		ret.Report = report
	}

	preFile := path.Join(dirname, "preamble.md")
	preamble, err := fileAsString(preFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%s': %s", preFile, err.Error())
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
