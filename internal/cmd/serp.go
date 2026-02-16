package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/dataforseo-cli/internal/outfmt"
)

type SerpCmd struct {
	Google SerpGoogleCmd `cmd:"" help:"Get Google organic SERP results"`
	PAA    SerpPAACmd    `cmd:"" name:"paa" help:"Get People Also Ask results"`
}

type SerpGoogleCmd struct {
	Keyword  string `help:"Search keyword" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
	Device   string `help:"Device type" default:"desktop" enum:"desktop,mobile"`
}

func (cmd *SerpGoogleCmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	data, err := client.SerpGoogle(ctx, cmd.Keyword, cmd.Location, cmd.Language, cmd.Device)
	if err != nil {
		return fmt.Errorf("serp google: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"POSITION", "TITLE", "URL", "DOMAIN"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				fmt.Sprintf("%d", d.RankGroup),
				d.Title,
				d.URL,
				d.Domain,
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range data {
		fmt.Fprintf(os.Stdout, "%2d. %s\n    %s\n    %s\n\n", d.RankGroup, d.Title, d.URL, d.Domain)
	}

	return nil
}

type SerpPAACmd struct {
	Keyword  string `help:"Search keyword" required:""`
	Location int    `help:"Location code" default:"2840"`
	Language string `help:"Language code" default:"en"`
}

func (cmd *SerpPAACmd) Run(ctx context.Context) error {
	client, err := getDataForSEOClient()
	if err != nil {
		return err
	}

	data, err := client.SerpPAA(ctx, cmd.Keyword, cmd.Location, cmd.Language)
	if err != nil {
		return fmt.Errorf("serp paa: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, data)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"TITLE", "URL", "DOMAIN"}
		rows := make([][]string, len(data))

		for i, d := range data {
			rows[i] = []string{
				d.Title,
				d.URL,
				d.Domain,
			}
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "No People Also Ask results found")
		return nil
	}

	for i, d := range data {
		fmt.Fprintf(os.Stdout, "%d. %s\n", i+1, d.Title)

		if d.URL != "" {
			fmt.Fprintf(os.Stdout, "   %s\n", d.URL)
		}

		fmt.Fprintln(os.Stdout)
	}

	return nil
}
