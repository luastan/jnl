package jnl

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	CSV_DATE_FORMAT = "2006-01-02 15:04:05"
)

func InfoInTable(cInfos []*CommandInfo, f *os.File) error {
	w := new(tabwriter.Writer)
	w.Init(f, 0, 8, 2, '\t', tabwriter.TabIndent)
	defer func(w *tabwriter.Writer) {
		_ = w.Flush()
	}(w)

	for _, cInfo := range cInfos {
		end := cInfo.End.Format(time.Stamp)
		elapsed := cInfo.End.Sub(cInfo.Start).String()

		if !cInfo.hasEnded() {
			end = ""
			elapsed = ""
		}
		_, err := fmt.Fprintf(
			w,
			"%s - %s\t%s\t%s\t%d\t%s\t\n",
			cInfo.Start.Format(time.Stamp),
			end,
			elapsed,
			strings.Join(cInfo.Command, " "),
			cInfo.ExitCode,
			cInfo.Dir,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func InfoInCSV(cInfos []*CommandInfo, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()
	err := w.Write([]string{
		"Start",
		"End",
		"Command",
		"ExitCode",
	})
	if err != nil {
		return err
	}
	for _, cInfo := range cInfos {
		end := cInfo.End.Format(CSV_DATE_FORMAT)
		if !cInfo.hasEnded() {
			end = ""
		}
		err := w.Write([]string{
			cInfo.Start.Format(CSV_DATE_FORMAT),
			end,
			strings.Join(cInfo.Command, " "),
			strconv.Itoa(cInfo.ExitCode),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
