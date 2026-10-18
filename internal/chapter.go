package internal

import (
	"strings"
)

type ChapterSpec struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}

type Chapter struct {
	ChapterSpec
	Start    int64
	Duration int64 // milliseconds
}

func NewChapter(spec *ChapterSpec) (*Chapter, error) {
	duration, err := GetDurationInMs(spec.Source)
	if err != nil {
		return nil, err
	}

	ch := &Chapter{
		Title:    spec.Title,
		Source:   spec.Source,
		Duration: duration,
	}
	return ch, nil
}

func loadChaptersInDir(dir string) ([]*Chapter, error) {
	results := []*Chapter{}

	err := VisitAudioSourcesInDir(dir, func(file string) error {
		title := file[strings.LastIndex(file, "/")+1 : strings.LastIndex(file, ".")]
		ch, err := NewChapter(&ChapterSpec{
			Title:  title,
			Source: file,
		})

		if err != nil {
			return err
		}

		results = append(results, ch)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

type ChapterVisitor func(ch *Chapter) error

func VisitChaptersInDir(dir string, vis ChapterVisitor) error {
	return VisitAudioSourcesInDir(dir, func(file string) error {
		title := file[strings.LastIndex(file, "/")+1 : strings.LastIndex(file, ".")]
		ch, err := NewChapter(&ChapterSpec{
			Title:  title,
			Source: file,
		})
		if err != nil {
			return err
		}

		return vis(ch)
	})
}
