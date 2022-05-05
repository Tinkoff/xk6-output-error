package pkg

import (
	"fmt"
	"go.k6.io/k6/metrics"
	"os"
	"sort"
	"strings"
	"time"

	"go.k6.io/k6/output"
	"go.uber.org/zap"
)

const NameOutput = "xk6-output-error"

// Output implements the lib.Output interface for saving to CSV files.
type Output struct {
	output.SampleBuffer

	params          output.Params
	periodicFlusher *output.PeriodicFlusher

	logger *zap.SugaredLogger
	config Config
}

// New Creates new instance of CSV output
func New(params output.Params) (output.Output, error) {
	return newOutput(params)
}

func newOutput(params output.Params) (*Output, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return &Output{}, fmt.Errorf("could not create logger: %w", err)
	}
	config, er := GetConsolidatedConfig(params.JSONConfig, params.ConfigArgument)
	if er != nil {
		return &Output{}, fmt.Errorf("could not parse configs: %w", er)
	}

	return &Output{
		logger: logger.Sugar(),
		config: config,
	}, nil
}

// Description returns a human-readable description of the output.
func (o *Output) Description() string {
	return NameOutput
}

// Start writes the csv header and starts a new output.PeriodicFlusher
func (o *Output) Start() error {
	o.logger.Debug("Starting...")
	pf, err := output.NewPeriodicFlusher(time.Duration(o.config.FTime.Int64), o.flushMetrics)
	if err != nil {
		return err
	}
	o.logger.Debug("Started!")
	o.periodicFlusher = pf
	// add custom plugin fields
	o.config.Fields = append([]string{"url", "status", "name", "method", "scenario", "group", "error", "error_code"}, o.config.Fields...)

	return nil
}

// Stop flushes any remaining metrics and stops the goroutine.
func (o *Output) Stop() error {
	o.logger.Debug("Stopping...")
	defer o.logger.Debug("Stopped!")
	o.periodicFlusher.Stop()
	return nil
}

// flushMetrics Writes samples to the csv file
func (o *Output) flushMetrics() {
	samples := o.GetBufferedSamples()
	if len(samples) == 0 {
		return
	}

	mapFields := map[string]string{}

	for _, sc := range samples {
		for _, sample := range sc.GetSamples() {
			for _, value := range o.config.Fields {
				buffer, _ := sample.GetTags().Get(value)
				if buffer != "" {
					mapFields[value] = buffer
				}
			}
			if err := flushStdErr(sample, mapFields, o.logger); err != nil {
				o.logger.Error("could not flush logs", err)
			}
		}
	}
}

func flushStdErr(sample metrics.Sample, mapFields map[string]string, log *zap.SugaredLogger) error {
	if (sample.Metric.Name == "checks" && sample.Value != 1) ||
		(sample.Metric.Name == "http_reqs" && sample.Metric.Tainted.Bool) {
		log.Debugw("StdErr logs", "time", sample.Time, "params", mapFields, "tags", sample.Tags)
		_, err := os.Stderr.WriteString(createString(sample.Time, mapFields, sample.Tags, log))
		return err
	}
	return nil
}

func createString(t time.Time, mapFields map[string]string, tags *metrics.SampleTags, log *zap.SugaredLogger) string {
	// 	add	check and extra tags
	for key, value := range tags.CloneTags() {
		mapFields[key] = value
	}
	// create sort list
	keys := make([]string, 0, len(mapFields))
	for k := range mapFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Debugw("Create keys list", "keys", keys)
	// create msg
	var msg string
	for _, k := range keys {
		msg += fmt.Sprintf(" %s=%q", k, mapFields[k])
	}
	out := fmt.Sprintf("time=%q level=error %s source=%q\n", t.Format(time.RFC3339), strings.TrimSpace(msg), NameOutput)
	log.Debugw("Create stdError string", "outString", out)
	return out
}
