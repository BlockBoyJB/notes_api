package validator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Если честно, то это отвратительный сервис. Непонятно каким образом он работает.
// Проблемы начались уже на уровне http запросов. Как выяснилось, application/json он не поддерживает
// Если слов >300 в строке он начинает багать и делать вид, что все отлично и ошибок нет.
// Из всех описанных ошибок работает только ERROR_UNKNOWN_WORD. Остальные вообще не отображаются.
// Флаги (options) не работают.
// 0 из 10 оценка сервису

const (
	defaultUrl = "https://speller.yandex.net/services/spellservice.json"
)

const (
	ErrUnknownWord = iota + 1
	ErrRepeatWord
	ErrCapitalization
	ErrTooManyErrors
)

var Errors = map[int]string{
	ErrUnknownWord:    "Unknown word",
	ErrRepeatWord:     "Repeating word",
	ErrCapitalization: "Incorrect uppercase use",
	ErrTooManyErrors:  "Too many errors in text",
}

type SpellResult struct {
	Code int      `json:"code"` // one of 4 codes
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func CheckText(text string) ([]SpellResult, error) {
	body := []byte("text=" + text) // Только в таком формате

	resp, err := http.Post(defaultUrl+"/checkText", "application/x-www-form-urlencoded", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []SpellResult

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
