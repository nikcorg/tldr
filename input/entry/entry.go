package entry

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tldr/storage"
)

// EditContext represents additional data useful for the editing context
type EditContext struct {
	Titles []string
}

// Edit takes an entry and presents it for editing via user input
func Edit(newEntry *storage.Entry, ctx *EditContext) error {
	return edit(newEntry, ctx, "Editing")
}

// Create creates a new entry
func Create(newEntry *storage.Entry, ctx *EditContext) error {
	return edit(newEntry, ctx, "Adding")
}

func edit(newEntry *storage.Entry, ctx *EditContext, mode string) error {
	var err error

	reader := bufio.NewReader(os.Stdin)

	for true {
		fmt.Printf("=== %s ===\n", mode)
		fmt.Printf("Title: %s\n", newEntry.Title)
		fmt.Printf("URL: %s\n", newEntry.URL)
		fmt.Printf("Unread: %v\n", newEntry.Unread)
		if len(newEntry.SourceURL) > 0 {
			fmt.Printf("Source: %s\n", newEntry.SourceURL)
		}
		if len(newEntry.RelatedURLs) > 0 {
			fmt.Println("Related:")
			for _, u := range newEntry.RelatedURLs {
				fmt.Printf("- %s\n", u)
			}
		}
		fmt.Println("---")
		fmt.Println("Press Enter to accept entry and exit")
		fmt.Println("L) see all titles")
		fmt.Println("R) toggle unread status")
		fmt.Println("T) enter custom title")
		fmt.Println("S) enter source URL")
		fmt.Println("U) enter related URL")
		fmt.Println("Q) quit without saving")
		fmt.Print("Your selection: ")

		var selection string
		selection, err = reader.ReadString('\n')
		if err != nil {
			return err
		}

		selection = strings.ToUpper(strings.TrimSpace(selection))

		if len(selection) == 0 {
			break
		}

		switch selection {
		case "L":
			for n, t := range ctx.Titles {
				fmt.Printf("%d) %s\n", n, t)
			}
			fmt.Print("Select title or press enter to keep current title: ")
			selection, _ = reader.ReadString('\n')
			selection = strings.TrimSpace(selection)
			if len(selection) > 0 {
				n, _ := strconv.Atoi(selection)
				newEntry.Title = ctx.Titles[n]
			}

		case "R":
			newEntry.Unread = !newEntry.Unread

		case "T":
			fmt.Printf("Enter title: ")
			selection, _ = reader.ReadString('\n')
			newEntry.Title = strings.TrimSpace(selection)

		case "S":
			fmt.Printf("Enter source: ")
			selection, _ = reader.ReadString('\n')
			newEntry.SourceURL = strings.ToLower(strings.TrimSpace(selection))

		case "U":
			fmt.Printf("Enter related: ")
			selection, _ = reader.ReadString('\n')
			newEntry.RelatedURLs = append(newEntry.RelatedURLs, strings.ToLower(strings.TrimSpace(selection)))

		case "Q":
			fmt.Println("Ok, quitting without saving.")
			os.Exit(0)

		default:
			fmt.Printf("I'm sorry, I don't understand '%s'. Please try again.\n", selection)
		}
	}

	return nil
}
