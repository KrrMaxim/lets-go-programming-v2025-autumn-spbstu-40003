package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"
	"errors"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

var (
	ErrParseXml     = errors.New("error with parsing XML")
	ErrOpenXml      = errors.New("error with opening xml file")
	ErrNumCode      = errors.New("error: invalid num code")
	ErrEmptyValue   = errors.New("error: empty Value")
	ErrInvalidValue = errors.New("error: invalid value")
)

type ExchangeTrade struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type ExchangeData struct {
	Valutes []ExchangeTrade `xml:"Valute"`
}

func (t *ExchangeTrade) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var tmp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := dec.DecodeElement(&tmp, &start); err != nil {
		return err
	}

	t.CharCode = strings.TrimSpace(tmp.CharCode)

	numStr := strings.TrimSpace(tmp.NumCode)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return ErrNumCode
	}
	t.NumCode = num

	val := strings.TrimSpace(tmp.Value)
	val = strings.ReplaceAll(val, ",", ".")

	if val == "" {
		return ErrEmptyValue
	}
	
	valFLoat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return ErrInvalidValue
	}
	t.Value = valFLoat

	return nil
}

func XmlParse(path string) ([]ExchangeTrade, error) {
	var exData ExchangeData 

	file, err := os.Open(path)
	if err != nil {
		return nil, ErrOpenXml
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Errorf("error with closing XML file: %w", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&exData); err != nil {
		return nil, ErrParseXml
	}

	sort.Slice(exData.Valutes, func(i, j int) bool {
		return exData.Valutes[i].Value > exData.Valutes[j].Value
	})

	return exData.Valutes, nil
}