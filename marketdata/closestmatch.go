package marketdata

import (
	"log"
	"strings"

	"github.com/sajari/fuzzy"
)

var model = fuzzy.NewModel()
var searchItems = make(map[string]Item, len(Items))

func init() {
	wordlist := make([]string, 0, len(Items))
	for _, item := range Items {
		// ignore items with @ (e.g. T6.{1,2,3} armor)
		if item.FriendlyName != "" && !strings.Contains(item.UniqueName, "@") {
			searchItems[preprocess(item.FriendlyName)] = item
		}
	}
	for k := range searchItems {
		wordlist = append(wordlist, preprocess(k))
	}

	// Train model from valid wordlist
	model.SetThreshold(1)
	model.SetDepth(3)
	log.Println("Training...")
	model.Train(wordlist)
	log.Println("Finished training...")
}

// preprocesses a string so that we can give a more accurate estimate
func preprocess(s string) string {
	s = strings.ToLower(s)

	// remove punctuation
	punctuation := []string{",", ".", ";", ":", "'s", "'", "\"", "(", ")", "[", "]"}
	for _, p := range punctuation {
		s = strings.ReplaceAll(s, p, "")
	}

	s = tierReplace(s)

	return s
}

// tierReplace replaces tokens matching t[1-8] with their corresponding word,
// such as novice, journeyman, etc
func tierReplace(s string) string {
	tokens := strings.Split(s, " ")
	tierMap := map[string]string{
		"t1": "beginner",
		"t2": "novice",
		"t3": "journeyman",
		"t4": "adept",
		"t5": "expert",
		"t6": "master",
		"t7": "grandmaster",
		"t8": "elder",
	}

	for i := range tokens {
		if val, ok := tierMap[tokens[i]]; ok {
			tokens[i] = val
		}
	}

	return strings.Join(tokens, " ")
}

func Closest(s string) Item {
	return searchItems[model.SpellCheck(preprocess(s))]
}
