package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	parsetorrentname "github.com/middelink/go-parse-torrent-name"
	"github.com/middelink/go-rssfeed-filter-download/qnap"
	"github.com/mmcdole/gofeed"
)

var (
	baseuri string
	user    string
	pass    string
	rssfeed = flag.String("rss", "http://horriblesubs.info/rss.php?res=sd", "rss feed to monitor")
	filter  = flag.String("filter", "1080,720", "Which resolutions we do not care for")
	silent  = flag.Bool("silent", false, "silence mode")
)

func init() {
	defaultUrl, _ := os.LookupEnv("QNAP_URL")
	defaultUser, _ := os.LookupEnv("QNAP_USER")
	defaultPass, _ := os.LookupEnv("QNAP_PASS")
	flag.StringVar(&baseuri, "baseuri", defaultUrl, "url for the qnap, e.g. http://192.168.1.5:8080/. Defaults to env QNAP_URL")
	flag.StringVar(&user, "user", defaultUser, "qnap user to log in as, defaults to env QNAP_USER")
	flag.StringVar(&pass, "pass", defaultPass, "qnap pass to log in with, defaults to env QNAP_PASS")
}

func match(title string) bool {
	return false /*strings.Contains(title, "Clockwork Planet") ||
	strings.Contains(title, "Tsugumomo") ||
	strings.Contains(title, "Uchouten") ||
	strings.Contains(title, "Alice to Zouroku") ||
	strings.Contains(title, "Atom - The Beginning") ||
	strings.Contains(title, "Eromanga-sensei") ||
	strings.Contains(title, "Re-Creators") ||
	strings.Contains(title, "Gin no Guardian") ||
	strings.Contains(title, "Detective Conan") ||
	strings.Contains(title, "Boku no Hero Academia") ||
	strings.Contains(title, "Berserk") ||
	strings.Contains(title, "Sagrada Reset")*/
}

func main() {
	flag.Parse()

	if baseuri == "" || user == "" || pass == "" {
		baseurl, user, pass = qnap.GetDefaults("qnap-downloader")
	}
	if baseuri == "" || user == "" || pass == "" {
		fmt.Printf("Either baseuri, user or passwd is not given\n")
		os.Exit(1)
	}

	resolutions := make(map[string]struct{}, 8)
	for _, v := range strings.Split(*filter, ",") {
		resolutions[strings.TrimSpace(v)] = struct{}{}
	}

	qf, err := qnap.filestation.New(baseuri, user, pass)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer qf.Close()
	var items map[string]bool
	if items, err = qf.GetList("/Download"); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	//fmt.Printf("items are %v\n", items)
	interest := make(map[string]struct{}, 128)
	for item := range items {
		//fmt.Printf("%+v\n", item)
		tor, err := parsetorrentname.Parse(item)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		interest[tor.Title] = struct{}{}
	}
	if !*silent {
		fmt.Println("You seem to be interested in:")
		for k := range interest {
			fmt.Printf("  %v\n", k)
		}
	}

	qd, err := qnap.downloadstation.New(baseuri, user, pass)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer qd.Close()
	var tasks map[string]qnap.TaskState
	if tasks, err = qd.TaskQuery(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	//fmt.Printf("tasks are: %v\n", tasks)

	var feed *gofeed.Feed
	fp := gofeed.NewParser()
	if feed, err = fp.ParseURL(*rssfeed); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	//fmt.Printf("feed is %v\n", feed)

	for _, item := range feed.Items {
		if nsTorrent, ok := item.Extensions["torrent"]; ok {
			if name, ok := nsTorrent["magnetURI"]; ok {
				item.Link = name[0].Value
			}
			if name, ok := nsTorrent["fileName"]; ok {
				item.Title = name[0].Value
			}
		}

		if !*silent {
			fmt.Printf("item %q\n", item.Title)
		}
		tor, err := parsetorrentname.Parse(item.Title)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		_, interested := interest[tor.Title]
		if match(item.Title) || interested {
			if _, ok := items[item.Title]; ok {
				if !*silent {
					fmt.Printf("  already have it\n")
				}
			} else if _, ok := tasks[item.Title]; ok {
				if !*silent {
					fmt.Printf("  already have an active task\n")
				}
			} else if _, ok := resolutions[tor.Resolution]; ok {
				if !*silent {
					fmt.Printf("  filtered resolution\n")
				}
			} else {
				if *silent {
					fmt.Printf("downloading %s\n", item.Title)
				} else {
					fmt.Printf("  downloading\n")
				}
				if err = qd.TaskAddUrl(item.Link); err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
			}
		}
	}
}
