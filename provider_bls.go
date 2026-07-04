// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tseventfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/thorsphere/tserr"
	"github.com/thorsphere/tstrading"
)

const (
	dateFormat     = "2006-01-02"
	datetimeFormat = "2006-01-02 15:04:05"
)

type BLSProvider struct {
	registrationKey string
	seriesID        string
	httpClient      *http.Client
	endpoint        string
	catalog         bool
	calculations    bool
	annualaverage   bool
	aspects         bool
}

type blsResponse struct {
	Status       string   `json:"status"`
	ResponseTime int      `json:"responseTime"`
	Message      []string `json:"message"`
	Results      struct {
		Series []blsSeries `json:"series"`
	} `json:"Results"`
}

type blsSeries struct {
	SeriesID string     `json:"seriesID"`
	Catalog  blsCatalog `json:"catalog"`
	Data     []blsData  `json:"data"`
}

type blsCatalog struct {
	SeriesTitle         string `json:"series_title"`
	SeriesID            string `json:"series_id"`
	Seasonality         string `json:"seasonality"`
	SurveyName          string `json:"survey_name"`
	SurveyAbbreviation  string `json:"survey_abbreviation"`
	MeasureDataType     string `json:"measure_data_type"`
	CommerceIndustry    string `json:"commerce_industry"`
	Occupation          string `json:"occupation"`
	OccupationWorkClass string `json:"occupation_work_class"`
	Area                string `json:"area"`
}

type blsData struct {
	Year       string        `json:"year"`
	Period     string        `json:"period"`
	PeriodName string        `json:"periodName"`
	Latest     string        `json:"latest"`
	Value      string        `json:"value"`
	Aspects    []blsAspect   `json:"aspects"`
	Footnotes  []blsFootnote `json:"footnotes"`
}

type blsAspect struct {
	Name      string        `json:"name"`
	Value     string        `json:"value"`
	Footnotes []blsFootnote `json:"footnotes"`
}

type blsFootnote struct {
	// empty object in the sample, define fields if needed
}

func NewBLSProvider(apiKey string) *BLSProvider {
	return &BLSProvider{
		registrationKey: apiKey,
		seriesID:        "CUSR0000SA0",
		catalog:         true,
		calculations:    true,
		annualaverage:   false,
		aspects:         false,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		endpoint:        "https://api.bls.gov/publicAPI/v2/timeseries/data/",
	}
}

func (f *BLSProvider) GetEvents(ctx context.Context, period *tstrading.Period) ([]tstrading.Event, error) {

	fromStr := period.From.Format("2006")
	toStr := period.To.Format("2006")
	url := fmt.Sprintf("%s?from=%s&to=%s&apikey=%s", f.endpoint, fromStr, toStr, f.registrationKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "New HTTP Request", Fn: url, Err: err})
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, tserr.Op(&tserr.OpArgs{Op: "HTTP Client Do", Fn: url, Err: err})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, tserr.StatusNotMatching(&tserr.StatusNotMatchingArgs{Expected: http.StatusOK, Actual: resp.StatusCode})
	}

	var blsResp blsResponse

	if err := json.NewDecoder(resp.Body).Decode(&blsResp); err != nil {
		return nil, err
	}

	var events []tstrading.Event
	for _, series := range blsResp.Results.Series {
		for _, d := range series.Data {
			// Parse year + period into a time (approximate).
			// BLS periods like "A01" = annual, "M01" = January, etc.
			t := parseBLSPeriod(d.Year, d.Period)

			value := parseFloat(d.Value)

			events = append(events, tstrading.Event{
				Name:     series.Catalog.SeriesTitle,
				Time:     t,
				Country:  "US",
				Actual:   &value,
				Unit:     strPtr(series.Catalog.MeasureDataType),
				Source:   "BLS",
				// ... populate other fields as needed ...
			})
		}
	}

	return events, nil
}

func fmpImpactLevel(impact string) tstrading.ImpactLevel {
	switch strings.ToLower(impact) {
	case "low":
		return tstrading.ImpactLow
	case "medium":
		return tstrading.ImpactMedium
	case "high":
		return tstrading.ImpactHigh
	default:
		return tstrading.ImpactUnknown // unknown impact level
	}
}
