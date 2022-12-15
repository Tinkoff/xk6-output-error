package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.k6.io/k6/metrics"
	"go.uber.org/zap"
)

func TestCreateString(t *testing.T) {
	logger := zap.NewNop().Sugar()
	r := metrics.NewRegistry().RootTagSet()

	type inputStruct struct {
		sample   metrics.Sample
		mapField map[string]string
	}

	type testPair struct {
		input  inputStruct
		output string
	}

	testMapFields := map[string]string{
		"myTag_1": "myTag_1",
		"myTag_2": "myTag_2",
	}

	sampleFields := metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: &metrics.Metric{
				Name:       "checks",
				Type:       0,
				Contains:   0,
				Thresholds: metrics.Thresholds{},
				Submetrics: nil,
				Sub:        &metrics.Submetric{},
				Sink:       nil,
			},
			Tags: metrics.NewRegistry().RootTagSet().WithTagsFromMap(map[string]string{
				"name":       "http://httpbin.org/1delete1?verb=delete",
				"scenario":   "default",
				"status":     "404",
				"myTag":      "myTag2",
				"url":        "http://tets.ru",
				"check":      "is verb correct",
				"error_code": "1404",
				"group":      "::DELETE-2",
				"method":     "DELETE",
			}),
		},
		Time:  time.Time{},
		Value: 0,
	}

	var tests = []testPair{
		{
			input:  inputStruct{metrics.Sample{TimeSeries: metrics.TimeSeries{Tags: r}}, map[string]string{}},
			output: "time=\"0001-01-01T00:00:00Z\" level=error  source=\"xk6-output-error\"\n",
		},
		{
			input:  inputStruct{metrics.Sample{TimeSeries: metrics.TimeSeries{Tags: r}}, testMapFields},
			output: "time=\"0001-01-01T00:00:00Z\" level=error myTag_1=\"myTag_1\" myTag_2=\"myTag_2\" source=\"xk6-output-error\"\n",
		},

		{
			input:  inputStruct{sampleFields, testMapFields},
			output: "time=\"0001-01-01T00:00:00Z\" level=error check=\"is verb correct\" error_code=\"1404\" group=\"::DELETE-2\" method=\"DELETE\" myTag=\"myTag2\" myTag_1=\"myTag_1\" myTag_2=\"myTag_2\" name=\"http://httpbin.org/1delete1?verb=delete\" scenario=\"default\" status=\"404\" url=\"http://tets.ru\" source=\"xk6-output-error\"\n",
		},
		{
			input:  inputStruct{sampleFields, map[string]string{}},
			output: "time=\"0001-01-01T00:00:00Z\" level=error check=\"is verb correct\" error_code=\"1404\" group=\"::DELETE-2\" method=\"DELETE\" myTag=\"myTag2\" name=\"http://httpbin.org/1delete1?verb=delete\" scenario=\"default\" status=\"404\" url=\"http://tets.ru\" source=\"xk6-output-error\"\n",
		},
	}

	for _, pair := range tests {
		result := createString(pair.input.sample.Time, pair.input.mapField, pair.input.sample.Tags, logger)
		assert.Equal(t, pair.output, result)
	}
}

func TestFlushStdErr(t *testing.T) {
	logger := zap.NewNop().Sugar()

	type testPair struct {
		sample   metrics.Sample
		mapField map[string]string
	}

	testMapFields := map[string]string{
		"myTag_1": "myTag_1",
		"myTag_2": "myTag_2",
	}

	sampleFields := metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: &metrics.Metric{
				Name:       "checks",
				Type:       0,
				Contains:   0,
				Thresholds: metrics.Thresholds{},
				Submetrics: nil,
				Sub:        &metrics.Submetric{},
				Sink:       nil,
			},
			Tags: metrics.NewRegistry().RootTagSet().WithTagsFromMap(map[string]string{
				"name":       "http://httpbin.org/1delete1?verb=delete",
				"scenario":   "default",
				"status":     "404",
				"myTag":      "myTag2",
				"url":        "http://tets.ru",
				"check":      "is verb correct",
				"error_code": "1404",
				"group":      "::DELETE-2",
				"method":     "DELETE",
			}),
		},
		Time:  time.Time{},
		Value: 0,
	}

	sampleFieldsErr := metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: &metrics.Metric{
				Name:       "myV",
				Type:       0,
				Contains:   0,
				Thresholds: metrics.Thresholds{},
				Submetrics: nil,
				Sub:        &metrics.Submetric{},
				Sink:       nil,
			},
			Tags: metrics.NewRegistry().RootTagSet().WithTagsFromMap(map[string]string{}),
		},
		Time:  time.Time{},
		Value: 0,
	}

	var tests = []testPair{
		{
			sampleFields, testMapFields,
		},
		{
			sampleFields, map[string]string{},
		},
		{
			sampleFieldsErr, map[string]string{},
		},
	}

	for _, pair := range tests {
		err := flushStdErr(pair.sample, pair.mapField, logger)
		assert.NoError(t, err)
	}
}
