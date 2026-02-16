package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/builtbyrobben/dataforseo-cli/internal/outfmt"
)

type KeywordsCmd struct {
	SearchVolume KeywordsSearchVolumeCmd `cmd:"" name:"search-volume" help:"Get search volume for keywords"`
	Difficulty   KeywordsDifficultyCmd   `cmd:"" help:"Get keyword difficulty scores"`
}

type KeywordsSearchVolumeCmd struct {
	Keywords string `help:"Comma-separated keywords" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
}

func (cmd *KeywordsSearchVolumeCmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	keywords := splitKeywords(cmd.Keywords)

	data, err := client.SearchVolume(ctx, keywords, cmd.Location, cmd.Language)
	if err != nil {
		return fmt.Errorf("search volume: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"KEYWORD", "SEARCH_VOLUME", "COMPETITION", "COMPETITION_LEVEL", "CPC"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				d.Keyword,
				fmt.Sprintf("%d", d.SearchVolume),
				fmt.Sprintf("%.2f", d.Competition),
				d.CompetitionLevel,
				fmt.Sprintf("%.2f", d.CPC),
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range data {
		fmt.Fprintf(os.Stdout, "%-30s  vol: %-8d  comp: %-6.2f  level: %-10s  cpc: $%.2f\n",
			d.Keyword, d.SearchVolume, d.Competition, d.CompetitionLevel, d.CPC)
	}

	return nil
}

type KeywordsDifficultyCmd struct {
	Keywords string `help:"Comma-separated keywords" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
}

func (cmd *KeywordsDifficultyCmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	keywords := splitKeywords(cmd.Keywords)

	data, err := client.KeywordDifficulty(ctx, keywords, cmd.Location, cmd.Language)
	if err != nil {
		return fmt.Errorf("keyword difficulty: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"KEYWORD", "DIFFICULTY"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				d.Keyword,
				fmt.Sprintf("%d", d.KeywordDifficulty),
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range data {
		fmt.Fprintf(os.Stdout, "%-30s  difficulty: %d\n", d.Keyword, d.KeywordDifficulty)
	}

	return nil
}

func splitKeywords(s string) []string {
	parts := strings.Split(s, ",")
	keywords := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			keywords = append(keywords, p)
		}
	}

	return keywords
}
