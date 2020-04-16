// package main

// import (
// 	"bufio"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path"
// 	"strconv"
// 	"strings"
// 	"time"
// 	"tldr/fetch"
// 	"tldr/storage"

// 	log "github.com/sirupsen/logrus"
// 	yaml "gopkg.in/yaml.v2"
// )

// var debugLogging bool = false
// var verboseLogging bool = false

// var configDir string = ""
// var sourceDir string = ""
// var sourceFile string = ""

// // var tldr []record
// var source storage.Source

// func main() {
// 	handleFlags()
// 	setLogLevel()

// 	if err := mainWithError(flag.Args()); err != nil {
// 		log.Fatalf("Error: %v", err.Error())
// 	}
// }

// func handleFlags() {
// 	flag.StringVar(&configDir, "c", "", "Where to find the configuration")
// 	flag.StringVar(&sourceDir, "d", "", "Where to find the links database file")
// 	flag.StringVar(&sourceFile, "f", "tldr.yaml", "The links database file")
// 	flag.BoolVar(&verboseLogging, "v", false, "Show verbose logging")
// 	flag.BoolVar(&debugLogging, "vv", false, "Show debug logging")
// 	flag.Parse()
// }

// func setLogLevel() {
// 	if debugLogging {
// 		log.SetLevel(log.DebugLevel)
// 	} else if verboseLogging {
// 		log.SetLevel(log.InfoLevel)
// 	} else {
// 		log.SetLevel(log.ErrorLevel)
// 	}
// }

// func mainWithError(args []string) error {
// 	var err error

// 	if len(args) == 0 {
// 		return fmt.Errorf("No arguments given, nothing to do")
// 	}

// 	configDir, err = getConfigDir()
// 	if err != nil {
// 		return fmt.Errorf("Error getting config dir: %w", err)
// 	}
// 	sourceDir, err = getSourceDir()
// 	if err != nil {
// 		return fmt.Errorf("Error getting data dir: %w", err)
// 	}

// 	log.Debugf("Config dir: %s", configDir)
// 	log.Debugf("Source dir: %s", sourceDir)
// 	log.Debugf("Source file: %s", sourceFile)

// 	var source []byte
// 	if source, err = readSource(); err != nil {
// 		return fmt.Errorf("Error reading data file (%s): %w", sourceFile, err)
// 	}
// 	log.Debugf("Read %d bytes of yaml", len(source))

// 	if err = yaml.Unmarshal(source, &tldr); err != nil {
// 		return fmt.Errorf("Error parsing data file: %w", err)
// 	}
// 	log.Debugf("Unmarshalled %d records", len(tldr))

// 	firstArg := args[0]
// 	if strings.HasPrefix(firstArg, "http") {
// 		err = addEntry(firstArg)
// 	}

// 	return nil
// }

// func addEntry(url string) error {
// 	log.Debugf("Fetching %v", url)
// 	var res *fetch.Details
// 	var err error
// 	if res, err = fetch.Preview(url); err != nil {
// 		return fmt.Errorf("Error fetching (%s): %w", url, err)
// 	}
// 	log.Debugf("Fetch result: %+v", res)

// 	var newEntry = entry{
// 		URL:    res.URL,
// 		Title:  res.Titles[0],
// 		Unread: true,
// 	}

// 	reader := bufio.NewReader(os.Stdin)

// 	for true {
// 		fmt.Println("=== Adding ===")
// 		fmt.Printf("Title: %s\n", newEntry.Title)
// 		fmt.Printf("URL: %s\n", newEntry.URL)
// 		fmt.Printf("Unread: %v\n", newEntry.Unread)
// 		if len(newEntry.SourceURL) > 0 {
// 			fmt.Printf("Source: %s\n", newEntry.SourceURL)
// 		}
// 		if len(newEntry.RelatedURLs) > 0 {
// 			fmt.Println("Related:")
// 			for _, u := range newEntry.RelatedURLs {
// 				fmt.Printf("- %s\n", u)
// 			}
// 		}
// 		fmt.Println("---")
// 		fmt.Println("Press Enter to accept entry and exit")
// 		fmt.Println("Enter L to see all titles")
// 		fmt.Println("Enter R to toggle unread status")
// 		fmt.Println("Enter T to enter custom title")
// 		fmt.Println("Enter S to enter source URL")
// 		fmt.Println("Enter U to enter related URL")
// 		fmt.Println("Enter Q to discard entry")
// 		fmt.Print("Your selection: ")

// 		var selection string
// 		selection, err = reader.ReadString('\n')
// 		selection = strings.ToUpper(strings.TrimSpace(selection))

// 		if len(selection) == 0 {
// 			break
// 		}

// 		switch selection {
// 		case "L":
// 			for n, t := range res.Titles {
// 				fmt.Printf("%d) %s\n", n, t)
// 			}
// 			fmt.Print("Select title or press enter to keep current title: ")
// 			selection, _ = reader.ReadString('\n')
// 			selection = strings.TrimSpace(selection)
// 			if len(selection) > 0 {
// 				n, _ := strconv.Atoi(selection)
// 				newEntry.Title = res.Titles[n]
// 			}
// 			break
// 		case "R":
// 			newEntry.Unread = !newEntry.Unread
// 			break
// 		case "T":
// 			fmt.Printf("Enter title: ")
// 			selection, _ = reader.ReadString('\n')
// 			newEntry.Title = strings.TrimSpace(selection)
// 			break
// 		case "S":
// 			fmt.Printf("Enter source: ")
// 			selection, _ = reader.ReadString('\n')
// 			newEntry.SourceURL = strings.TrimSpace(selection)
// 			break
// 		case "U":
// 			fmt.Printf("Enter related: ")
// 			selection, _ = reader.ReadString('\n')
// 			newEntry.RelatedURLs = append(newEntry.RelatedURLs, strings.TrimSpace(selection))
// 			break
// 		case "Q":
// 			fmt.Println("Ok, quitting without saving.")
// 			os.Exit(0)
// 			break
// 		default:
// 			fmt.Printf("I'm sorry, I don't understand '%s'. Please try again.\n", selection)
// 			break
// 		}
// 	}

// 	addEntryToTLDR(newEntry)

// 	var yamlString []byte
// 	yamlString, err = yaml.Marshal(tldr)
// 	if err != nil {
// 		return fmt.Errorf("Error serialising yaml: %w", err)
// 	}
// 	log.Debugf("Marshalled %d entries into %d bytes", len(tldr), len(yamlString))

// 	err = writeSource(yamlString)
// 	if err != nil {
// 		log.Debugf("Error writing yaml: %s", err.Error())
// 		return err
// 	}

// 	fmt.Println("Done")

// 	return nil
// }

// func addEntryToTLDR(newEntry entry) {
// 	y1, m1, d1 := time.Now().Date()

// 	if len(tldr) == 0 {
// 		tldr = []record{
// 			{
// 				Date:    time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC),
// 				Entries: []entry{newEntry},
// 			},
// 		}
// 		return
// 	}

// 	y2, m2, d2 := tldr[0].Date.Date()

// 	if y1 == y2 && m1 == m2 && d1 == d2 {
// 		log.Debug("Entry for today already exists, appending")
// 		tldr[0].Entries = append(tldr[0].Entries, newEntry)
// 	} else {
// 		log.Debug("Entry for today doesn't exist, creating")
// 		tldr = append([]record{
// 			{
// 				Date:    time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC),
// 				Entries: []entry{newEntry},
// 			},
// 		}, tldr...)
// 	}
// }

// func getConfigDir() (string, error) {
// 	var configDir string
// 	var err error

// 	// If a value was passed on the command line, don't overrule it
// 	if configDir != "" {
// 		return configDir, nil
// 	}

// 	if configDir, err = os.UserConfigDir(); err != nil {
// 		return "", err
// 	}

// 	return path.Join(configDir, "tldr"), nil
// }

// func getSourceDir() (string, error) {
// 	var configDir string
// 	var err error

// 	if sourceDir != "" {
// 		return sourceDir, nil
// 	}

// 	if configDir, err = os.UserHomeDir(); err != nil {
// 		return "", err
// 	}

// 	return path.Join(configDir, "tldr"), nil
// }

// func readSource() ([]byte, error) {
// 	var err error
// 	var source []byte

// 	fullSourcePath := path.Join(sourceDir, sourceFile)

// 	_, err = os.Stat(fullSourcePath)
// 	if err != nil && os.IsNotExist(err) {
// 		log.Debugf("Source file does not exist: %s", fullSourcePath)
// 		return []byte{}, nil
// 	} else if err != nil {
// 		return nil, fmt.Errorf("Error reading source: %w", err)
// 	}

// 	source, err = ioutil.ReadFile(fullSourcePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error reading %s: %w", fullSourcePath, err)
// 	}

// 	return source, nil
// }

// func writeSource(b []byte) error {
// 	var err error

// 	err = os.MkdirAll(sourceDir, 0755)
// 	if err != nil && !os.IsExist(err) {
// 		return fmt.Errorf("Error creating data dir: %s %w", sourceDir, err)
// 	}

// 	fullSourcePath := path.Join(sourceDir, sourceFile)
// 	err = ioutil.WriteFile(fullSourcePath, b, 0644)
// 	if err != nil && !os.IsNotExist(err) {
// 		log.Debugf("Error writing %s: %s", fullSourcePath)
// 		return fmt.Errorf("Error writing %s: %w", fullSourcePath, err)
// 	}

// 	log.Debugf("Wrote %d bytes into %s", len(b), fullSourcePath)

// 	return nil
// }
