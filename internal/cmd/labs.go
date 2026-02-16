package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/dataforseo-cli/internal/outfmt"
)

type LabsCmd struct {
	Related     LabsRelatedCmd     `cmd:"" help:"Get related keywords"`
	Suggestions LabsSuggestionsCmd `cmd:"" help:"Get keyword suggestions"`
}

type LabsRelatedCmd struct {
	Keyword  string `help:"Seed keyword" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
	Limit    int    `help:"Max results" default:"10"`
}

func (cmd *LabsRelatedCmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	data, err := client.RelatedKeywords(ctx, cmd.Keyword, cmd.Location, cmd.Language, cmd.Limit)
	if err != nil {
		return fmt.Errorf("related keywords: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"KEYWORD", "SEARCH_VOLUME", "COMPETITION", "CPC"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				d.Keyword,
				fmt.Sprintf("%d", d.SearchVolume),
				fmt.Sprintf("%.2f", d.Competition),
				fmt.Sprintf("%.2f", d.CPC),
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range data {
		fmt.Fprintf(os.Stdout, "%-30s  vol: %-8d  comp: %-6.2f  cpc: $%.2f\n",
			d.Keyword, d.SearchVolume, d.Competition, d.CPC)
	}

	return nil
}

type LabsSuggestionsCmd struct {
	Keyword  string `help:"Seed keyword" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
	Limit    int    `help:"Max results" default:"10"`
}

func (cmd *LabsSuggestionsCmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	data, err := client.KeywordSuggestions(ctx, cmd.Keyword, cmd.Location, cmd.Language, cmd.Limit)
	if err != nil {
		return fmt.Errorf("keyword suggestions: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"KEYWORD", "SEARCH_VOLUME", "COMPETITION", "CPC"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				d.Keyword,
				fmt.Sprintf("%d", d.SearchVolume),
				fmt.Sprintf("%.2f", d.Competition),
				fmt.Sprintf("%.2f", d.CPC),
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range data {
		fmt.Fprintf(os.Stdout, "%-30s  vol: %-8d  comp: %-6.2f  cpc: $%.2f\n",
			d.Keyword, d.SearchVolume, d.Competition, d.CPC)
	}

	return nil
}
