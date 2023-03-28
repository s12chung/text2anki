package text

import (
	"strings"
	"time"

	"github.com/asticode/go-astisub"
)

// ParseSubtitles parses the subtitles to return an array of Text
func ParseSubtitles(sourceFile, translationFile string) ([]Text, error) {
	sourceSub, err := astisub.OpenFile(sourceFile)
	if err != nil {
		return nil, err
	}
	translationSub, err := astisub.OpenFile(translationFile)
	if err != nil {
		return nil, err
	}

	texts := make([]Text, 0, len(sourceSub.Items))
	translationIndex := 0
	for i, item := range sourceSub.Items {
		var nextItem *astisub.Item
		if i+1 < len(sourceSub.Items) {
			nextItem = sourceSub.Items[i+1]
		}

		var translation []string
		for ; translationIndex < len(translationSub.Items); translationIndex++ {
			translationItem := translationSub.Items[translationIndex]
			if !isRelatedToItem(item, nextItem, translationItem) {
				break
			}
			translation = append(translation, itemString(translationItem))
		}

		texts = append(texts, Text{
			Text:        itemString(item),
			Translation: strings.Join(translation, " "),
		})
	}
	return texts, nil
}

func isRelatedToItem(item, nextItem, translationItem *astisub.Item) bool {
	if nextItem == nil {
		return true
	}
	iOverlap, nextOverlap := itemOverlap(item, translationItem), itemOverlap(nextItem, translationItem)
	if iOverlap != nextOverlap {
		return iOverlap > nextOverlap
	}
	// equal non-zero overlap
	if iOverlap != 0 {
		return true
	}
	// no overlap, so calculate distance
	return translationItem.StartAt-item.EndAt < nextItem.StartAt-translationItem.EndAt
}

func itemOverlap(a, b *astisub.Item) time.Duration {
	if b.StartAt < a.StartAt {
		b, a = a, b
	}
	if a.EndAt < b.StartAt {
		return 0
	}
	return a.EndAt - b.StartAt
}

func itemString(i *astisub.Item) string {
	os := make([]string, len(i.Lines))
	for i, l := range i.Lines {
		os[i] = strings.TrimSpace(l.String())
	}
	return strings.Join(os, " ")
}
