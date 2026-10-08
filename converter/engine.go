package converter

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
	"bufio"

	"strings"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) StartWorker() {
	for i := 0; i < config.ConversionWorkers; i++ {
		go e.workerLoop()
	}
}

func (e *Engine) workerLoop() {
	for id := range data.ConversionQueue {
		item, exists := data.Store.GetItem(id)
		if !exists {
			continue
		}

		// item.SetProgress(0)

		cmd, err := processor.Process.CompileConversionCmd(item)
		if err != nil {
			item.SetStatus(status.Failed)
			continue
		}

		conversionState := processor.Process.NewConversionState()

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			item.SetStatus(status.Failed)
			continue
		}
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			item.SetStatus(status.Failed)
			continue
		}

		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			rawLine := scanner.Text()
			cleanLine := strings.TrimSpace(rawLine)
			if cleanLine == "" {
				continue
			}

			pct, customMeta, showInLog := processor.Process.ParseProgressLine(cleanLine, conversionState)

			if pct >= 0 {
				// item.SetProgress(pct)
			}

			if customMeta != "" {
				// item.SetETA(customMeta)
			}

			if showInLog {
				// item.AppendTerminalLog(cleanLine)
			}
		}
		if err := scanner.Err(); err != nil {
			return
		}

		if err := cmd.Wait(); err != nil {
			item.SetStatus(status.Failed)
			continue
		}

		// item.SetProgress(100)
		item.SetStatus(status.Verified)
	}
}
