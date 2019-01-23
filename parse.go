package gathering

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Event\.MatchCreated[\r\n\s\w\W]*matchId":\s"39130126-411d-4071-a4c7-e5a1e9596442"
const (
	logSplitRegex           = `(\[UnityCrossThreadLogger\]|\[Client GRE\])`   // rm
	isCollectionRegex       = `<==\sPlayerInventory\.GetPlayerCardsV3\(\d*\)` //rm
	isDeckListRegex         = `<==\sDeck\.GetDeckLists\(\d*\)`
	isPlayerInventoryRegex  = `<==\sPlayerInventory\.GetPlayerInventory\(\d*\)` //rm
	isPlayerConnectionRegex = `ClientToMatchServiceMessageType_AuthenticateRequest`
	isRankInfoRegex         = `<==\sEvent\.GetCombinedRankInfo\(\d*\)` //rm
	isMatchPlayerCourse     = `<==\sEvent\.GetPlayerCourse\(\d*\)`
	isMatchStartRegex       = `Incoming\sEvent\.MatchCreated`
	isMatchEndRegex         = `DuelScene\.GameStop`
)

// ParseCollection looks for a MTGA Collection JSON object in a given input
func ParseCollection(raw string) (map[string]int, error) {
	isCollection := regexp.MustCompile(isCollectionRegex)
	texts := strings.Split(raw, "[UnityCrossThreadLogger]")
	match := getLastRegex(texts, isCollection, 2)
	if match == "" {
		return nil, errors.New("collection not found")
	}
	var collection map[string]int
	err := json.Unmarshal([]byte(match), &collection)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

// ParsePlayerInventory gets a players details
func ParsePlayerInventory(raw string) (*ArenaPlayerInventory, error) {
	isPlayerInventory := regexp.MustCompile(isPlayerInventoryRegex)
	texts := strings.Split(raw, "[UnityCrossThreadLogger]")
	match := getLastRegex(texts, isPlayerInventory, 2)
	if match == "" {
		return nil, errors.New("inventory not found")
	}
	var inventory ArenaPlayerInventory
	err := json.Unmarshal([]byte(match), &inventory)
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// ParseRankInfo finds a users rank info
func ParseRankInfo(raw string) (*ArenaRankInfo, error) {
	isRankInfo := regexp.MustCompile(isRankInfoRegex)
	texts := splitLogText(raw, logSplitRegex)
	rankJSON := getLastRegex(texts, isRankInfo, 2)
	if rankJSON == "" {
		return nil, errors.New("rank info not found")
	}
	var rank ArenaRankInfo
	err := json.Unmarshal([]byte(rankJSON), &rank)
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

// ParseAuthRequest parses the auth request (for username)
func ParseAuthRequest(raw string) (*ArenaAuthRequest, error) {
	isPlayerConnection := regexp.MustCompile(isPlayerConnectionRegex)
	texts := splitLogText(raw, logSplitRegex)
	match := getLastRegex(texts, isPlayerConnection, 1)
	if match == "" {
		return nil, errors.New("auth request not found")
	}
	var auth ArenaAuthRequest
	err := json.Unmarshal([]byte(match), &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func splitLogText(raw, split string) []string {
	r := regexp.MustCompile(split)
	return r.Split(raw, -1)
}

// Split
// Split the file into log objects, which has a timestamp for the log (if present),
// the type of log statement (unity logger or client gre), and the body of the log,
// which is the content until the next log statement. Ignores everything until the
// first log statement.
// has methods for trying to smart parse the body as JSON.
// Ignores common issues when doing so.
func splitLogText2(raw, split string) {
	runes := []rune(raw)
	r := regexp.MustCompile(`\[UnityCrossThreadLogger\].*|\[Client GRE\].*`)
	lines := r.FindAllStringIndex(raw, -1)
	fmt.Println(string(runes[lines[0][0]:lines[0][1]]))
}

func getLastRegex(texts []string, regex *regexp.Regexp, i int) string {
	var match string
	for _, rawLogText := range texts {
		if regex.MatchString(rawLogText) {
			split := strings.SplitN(rawLogText, "\n", i+1)
			if len(split) != i+1 {
				continue
			}
			match = split[i]
		}
	}
	return match
}

func parseJSONBackoff(s string, res interface{}) error {
	if s == "" {
		log.Println("json backoff failed with empty string")
		return errors.New("unable to parse")
	}
	if err := json.Unmarshal([]byte(s), &res); err != nil {
		split := strings.Split(s, "\n")
		split = split[:len(split)-1]
		return parseJSONBackoff(strings.Join(split, "\n"), res)
	}
	return nil
}

func stripNonJSON(text string) string {
	text = regexp.MustCompile(`<<<<<<<<<<.*`).ReplaceAllString(text, `$1.$2`)
	text = regexp.MustCompile(`\[\w.*`).ReplaceAllString(text, `$1.$2`)
	text = regexp.MustCompile(`\dx[\d\w]+.*`).ReplaceAllString(text, `$1.$2`)
	text = strings.TrimLeftFunc(text, func(r rune) bool {
		return r != '{' && r != '['
	})
	return strings.TrimRightFunc(text, func(r rune) bool {
		return r != '}' && r != ']'
	})
}
