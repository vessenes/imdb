package main

//go:generate msgp

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/peterh/liner"
	"github.com/tinylib/msgp/msgp"
)

type Trie struct {
	Value    string
	Children []Child
	End      bool
}

type Child struct {
	Value rune
	Link  *Trie
}

func NewTrie() *Trie {
	return &Trie{}
}

func NewChild(value rune) Child {
	child := Child{}
	child.Value = value
	return child
}

func (t *Trie) Add(value rune) *Trie {
	link := NewTrie()
	link.Value = t.Value + string(value)
	child := NewChild(value)
	child.Link = link

	t.Children = append(t.Children, child)

	return link
}

func (t *Trie) Insert(value string) {
	trie := t

	for i, val := range value {
		link, err := trie.Find(val)

		if i < len(value) {
			if err != nil {
				link = trie.Add(val)
			}

			trie = link
		}

	}

	trie.End = true
}

func (t *Trie) Find(value rune) (*Trie, error) {
	if t != nil {
		for _, child := range t.Children {
			if child.Value == value {
				return child.Link, nil
			}
		}
	}
	return nil, errors.New("Value not found")
}

func (t *Trie) AllPrefixes() ([]string, error) {
	results := []string{}
	if t != nil {
		if t.End == true {
			results = append(results, t.Value)
		}

		for _, val := range t.Children {
			prefixes, err := val.Link.AllPrefixes()
			if err != nil {
				break
			}
			results = append(results, prefixes...)
		}
	} else {
		return results, errors.New("No end found, so no completions available")
	}
	return results, nil
}

func (t *Trie) AutoComplete(prefix string) ([]string, error) {
	trie := t

	for _, val := range prefix {
		link, _ := trie.Find(val)
		trie = link
	}

	results, err := trie.AllPrefixes()

	return results, err
}

type IMDB struct {
	A2aid map[string]string // Actor name to id
	Aid2a map[string]string // Actor ID to name

	F2fids map[string]map[string]bool // Film name to film ids
	Fid2f  map[string]string          // Film id to name

	Aid2Fids map[string]map[string]bool
	Fid2Aids map[string]map[string]bool

	F2Fs map[string]uint8 // Film to Film, an edge is if any film shares actors. key is film1.film2 val is film1 score
	F2F  map[string]map[string]bool

	Frating       map[string]uint8
	FBroadlyRated map[string]bool

	Actors Trie
	Films  Trie
	AIds   Trie
	FIds   Trie

	AVor    map[string]float64
	ARating map[string]float64
}

func (i *IMDB) filmRatingMinusA(aid string, fid string) float64 {
	if len(i.Fid2Aids[fid]) == 1 { // this is the only 'important' actor, just return the films rating
		return float64(i.Frating[fid]) / 10.0
	}
	actors := make(map[string]bool)
	for aid2 := range i.Fid2Aids[fid] {
		actors[aid2] = true
	}
	delete(actors, aid)

	var ratingTotal float64
	for aid2 := range actors {
		ratingTotal += i.ARating[aid2]
	}
	return ratingTotal / float64(len(actors))
}

func (i *IMDB) calcVor() {
	log.Println("Calculating Actor Average Rating")
	i.ARating = make(map[string]float64, 1000000)
	i.AVor = make(map[string]float64, 1000000)

	ct := 0
	for aid := range i.Aid2a {
		var ratings float64
		var skips int
		for fid := range i.Aid2Fids[aid] {

			if i.Frating[fid] == 0 || !i.FBroadlyRated[fid] {
				skips++
				continue
			}
			ratings += float64(i.Frating[fid])
		}
		if len(i.Aid2Fids[aid]) == skips { // only zero rated films
			continue
		}
		i.ARating[aid] = ratings / (10.0 * float64(len(i.Aid2Fids[aid])-skips))
		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
		ct++
	}

	log.Println("Calculating AVOR")

	for aid := range i.Aid2a {
		var vor float64
		var skips int
		for fid := range i.Aid2Fids[aid] {
			if i.Frating[fid] == 0 || !i.FBroadlyRated[fid] {
				skips++
				continue
			}
			//vor is film rating - film rating without the actor
			vor += float64(i.Frating[fid])/10.0 - i.filmRatingMinusA(aid, fid)
		}
		if len(i.Aid2Fids[aid]) == skips { // only zero rated films
			continue
		}

		i.AVor[aid] = vor / (float64(len(i.Aid2Fids[aid]) - skips))

		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
		ct++
	}
}

/*
func (i *IMDB) fillEdges() {
	// for each film, get actors then their films then make edges
	ct := 0
	for fid := range i.Fid2Aids { // for each film
		for aid := range i.Fid2Aids[fid] { // get actors
			for fid2 := range i.Aid2Fids[aid] { // then their films
				ct++
				if fid == fid2 {
					continue
				}
				i.F2Fs[fid+"."+fid2] = 1

				if _, ok := i.F2F[fid]; !ok {
					i.F2F[fid] = make(map[string]bool, 500)
				}

				i.F2F[fid][fid2] = true

				if ct%100000 == 0 {
					os.Stdout.Write([]byte("."))
				}
				if ct%1000000 == 0 {
					fmt.Println("Parsing", fid, i.Fid2f[fid])
				}
			}
		}
	}
	fmt.Println()
}
*/
var imdb IMDB

func In(check string, all []string) bool {
	for _, item := range all {
		if check == item {
			return true
		}
	}
	return false
}

func dbFill() {
	imdb.A2aid = make(map[string]string, 15000000)
	imdb.Aid2a = make(map[string]string, 15000000)

	var aid2fid = make(map[string]map[string]bool, 15000000)
	var fid2aid = make(map[string]map[string]bool, 15000000)

	log.Println("Parsing Title Basics")

	fid2f := make(map[string]string, 1000000)
	f2fid := make(map[string]map[string]bool, 1000000)

	pfile, err := os.Open("./title.basics.tsv")
	if err != nil {
		log.Println(err)
	}
	preader := csv.NewReader(pfile)
	preader.Comma = '\t'
	rec, _ := preader.Read()

	for ct := 0; ; ct++ {
		rec, _ = preader.Read()
		if rec == nil {
			break
		}
		if len(rec) < 3 {
			continue
		}

		if rec[1] != "movie" {
			continue
		}

		fid2f[rec[0]] = rec[2]

		if _, ok := f2fid[rec[2]]; !ok {
			f2fid[rec[2]] = make(map[string]bool, 10)
		}
		f2fid[rec[2]][rec[0]] = true

		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
	}

	imdb.F2fids = f2fid
	imdb.Fid2f = fid2f

	log.Println("Do we have Kid Galahad / tt0056138?")
	log.Println("ID2f:", imdb.Fid2f["tt0056138"])
	log.Println("F2id:", imdb.F2fids["Kid Galahad"])

	namefile, err := os.Open("./name.basics.tsv")

	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(namefile)

	reader.Comma = '\t'

	rec, _ = reader.Read()

	log.Println("Reading Actors, 1 . = 100k actors")

	for ct := 0; ; ct++ {
		rec, _ = reader.Read()
		if rec == nil {
			break
		}
		if len(rec) < 2 {
			continue
		}

		imdb.A2aid[rec[1]] = rec[0]

		imdb.Aid2a[rec[0]] = rec[1]

		// Fill "Known fors" - they have more movies some times. Only if an actor.

		roleTypes := strings.Split(rec[4], ",")
		if roleTypes[0] == "actor" || roleTypes[0] == "actress" {

			knowns := strings.Split(rec[5], ",")

			for _, fid := range knowns {
				if _, ok := imdb.Fid2f[fid]; !ok { // we aren't interested in this one (TV?)
					continue
				}

				if _, ok := aid2fid[rec[0]]; !ok {
					aid2fid[rec[0]] = make(map[string]bool, 10)
				}

				aid2fid[rec[0]][fid] = true

				if _, ok := fid2aid[fid]; !ok {
					fid2aid[fid] = make(map[string]bool, 10)
				}

				fid2aid[fid][rec[0]] = true

			}

		}
		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
	}
	fmt.Println()
	namefile.Close()

	log.Println("Matching Actors Up With Their Films..")

	pfile, err = os.Open("./title.principals.tsv")
	if err != nil {
		log.Println(err)
	}
	preader = csv.NewReader(pfile)
	preader.Comma = '\t'
	rec, _ = preader.Read()

	for ct := 0; ; ct++ {
		rec, _ = preader.Read()
		if rec == nil {
			break
		}
		if len(rec) < 3 {
			continue
		}

		if rec[3] != "actor" && rec[3] != "actress" {
			continue
		}

		if _, ok := imdb.Fid2f[rec[0]]; !ok { // skip any info about something not in our film list
			continue
		}

		if _, ok := aid2fid[rec[2]]; !ok {
			aid2fid[rec[2]] = make(map[string]bool, 10)
		}

		aid2fid[rec[2]][rec[0]] = true

		if _, ok := fid2aid[rec[0]]; !ok {
			fid2aid[rec[0]] = make(map[string]bool, 10)
		}

		fid2aid[rec[0]][rec[2]] = true

		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
	}

	imdb.Aid2Fids = aid2fid
	imdb.Fid2Aids = fid2aid

	pfile, err = os.Open("./title.ratings.tsv")
	if err != nil {
		log.Println(err)
	}
	preader = csv.NewReader(pfile)
	preader.Comma = '\t'
	rec, _ = preader.Read()

	imdb.Frating = make(map[string]uint8, 1000000)
	imdb.FBroadlyRated = make(map[string]bool, 1000000)

	var rating string
	var film string
	for ct := 0; ; ct++ {
		rec, _ = preader.Read()
		if rec == nil {
			break
		}
		if len(rec) < 3 {
			continue
		}
		rating = rec[1]
		film = rec[0]
		fRating, err := strconv.ParseFloat(rating, 32)
		if err != nil {
			continue
		}
		imdb.Frating[film] = uint8(fRating * 10.0)
		numrating, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			continue
		}

		if numrating >= 1000 {
			imdb.FBroadlyRated[film] = true
		}

		if ct%100000 == 0 {
			os.Stdout.Write([]byte("."))
		}
	}

	imdb.F2Fs = make(map[string]uint8, 100000000)
	imdb.F2F = make(map[string]map[string]bool, 100000000)

	log.Println("About to fill film edges..")

	//imdb.fillEdges()
	imdb.calcVor()

	log.Println("Making Trie of Actors..")
	imdb.Actors = *(NewTrie())

	for actor, aid := range imdb.A2aid {
		imdb.Actors.Insert(actor)
		imdb.AIds.Insert(aid)
	}

	log.Println("Making Trie of Films..")
	imdb.Films = *(NewTrie())

	for film, fids := range imdb.F2fids {
		imdb.Films.Insert(film)
		for fid := range fids {
			imdb.FIds.Insert(fid)
		}
	}

	log.Println("About to write the db")
	imdbFile, _ := os.Create("./imdb.msgp")
	w := msgp.NewWriter(imdbFile)

	imdb.EncodeMsg(w)
	w.Flush()
	log.Println("Done")
}

var history_fn = filepath.Join(os.TempDir(), "~/.imdb_history")

func completeInput(input string) []string {
	var resp []string

	fmt.Println()

	res, _ := imdb.Actors.AutoComplete(input)
	if len(res) > 0 {
		fmt.Println("Actors", res)
		resp = append(resp, res...)
	}

	res, _ = imdb.AIds.AutoComplete(input)
	if len(res) > 0 {
		fmt.Println("Actor IDs", res)
		resp = append(resp, res...)
	}

	res, _ = imdb.Films.AutoComplete(input)
	if len(res) > 0 {
		fmt.Println("Films", res)
		resp = append(resp, res...)
	}

	res, _ = imdb.FIds.AutoComplete(input)
	if len(res) > 0 {
		fmt.Println("Film Ids", res)
		resp = append(resp, res...)
	}

	return resp
}

type Actor struct {
	Name string
	ID   string
	Vor  float64
}

type vorList []Actor

type ByVor []Actor

func (a ByVor) Len() int           { return len(a) }
func (a ByVor) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVor) Less(i, j int) bool { return a[i].Vor > a[j].Vor }

var actors vorList

func (i *IMDB) vorRankings(start int, stop int) []Actor {
	// Sort vors and return start to stop from the top

	if len(actors) == 0 { // fill me
		log.Println("Doing first fill of sorted Vor list")
		for aid, vor := range i.AVor {
			if math.IsNaN(vor) {
				continue
			}
			if len(i.Aid2Fids[aid]) < 3 { // must have 3 movies
				continue
			}
			ct := 0
			for fid := range i.Aid2Fids[aid] {
				if i.Frating[fid] == 0 {
					continue
				}
				if !i.FBroadlyRated[fid] {
					continue
				}

				ct++
			}
			if ct < 5 { // must have 5 broadly rated movies
				continue
			}
			actors = append(actors, Actor{i.Aid2a[aid], aid, vor})
		}

		log.Println("Have", len(actors), "Actors in the Vor")
	}
	sort.Sort(ByVor(actors))

	if start > len(actors)-1 {
		return actors[:len(actors)-1]
	}
	if stop > len(actors)-1 {
		return actors[start : len(actors)-1]
	}
	return actors[start:stop]
}

type FilmLink struct {
	Fid       string
	Film      string
	ActorLink string
	Rating    uint8
}

type ByRating []FilmLink

func (a ByRating) Len() int           { return len(a) }
func (a ByRating) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRating) Less(i, j int) bool { return a[i].Rating < a[j].Rating }

func doCLI() {

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			log.Println("See ya")
			break
		}
		line.AppendHistory(input)

		table := tablewriter.NewWriter(os.Stdout)
		if len(input) > 2 && input[:2] == "v:" { // vor query
			nums := strings.Split(input[2:], " ")
			num1, err := strconv.Atoi(nums[0])
			if err != nil {
				continue
			}
			num2, err := strconv.Atoi(nums[1])
			if err != nil {
				continue
			}

			actors := imdb.vorRankings(num1, num2)

			for _, a := range actors {
				score := fmt.Sprintf("%.3f", a.Vor)
				table.Append([]string{a.ID, a.Name, score})
			}

			table.Render()

		} else if len(input) > 2 && input[:2] == "l:" { //links from this film by actor
			fid := input[2:]
			fmt.Println("Finding all linked films, connected by major actors, from", input[2:])
			// actors in this film:

			log.Println(fid)
			var films = make(map[string]FilmLink, 100) // fid -> film
			for aid := range imdb.Fid2Aids[fid] {
				for fid := range imdb.Aid2Fids[aid] {
					films[fid] = FilmLink{fid, imdb.Fid2f[fid], imdb.Aid2a[aid], imdb.Frating[fid]}
				}
			}

			finalList := make([]FilmLink, len(films))
			ct := 0
			for _, f := range films {
				finalList[ct] = f
				ct++
			}

			sort.Sort(ByRating(finalList))

			table.SetHeader([]string{"Film ID", "Film", "Linker", "Rating", "Broadly Rated?"})
			for _, item := range finalList {
				table.Append([]string{item.Fid, item.Film, item.ActorLink, strconv.Itoa(int(item.Rating)), strconv.FormatBool(imdb.FBroadlyRated[item.Fid])})
			}
			table.Render()
			table.ClearRows()
		} else {
			if a, ok := imdb.A2aid[input]; ok { // actor query
				fmt.Println("Actor:", input)
				fmt.Println("ID: ", a)
				fmt.Println("Value:", imdb.AVor[a])
				fmt.Println("Average Rating:", imdb.ARating[a])
				fmt.Println("Filmography:")
				table.SetHeader([]string{"ID", "Film", "Rating", "Broadly Rated?"})
				for fid := range imdb.Aid2Fids[a] {
					table.Append([]string{fid, imdb.Fid2f[fid], strconv.Itoa(int(imdb.Frating[fid])), strconv.FormatBool(imdb.FBroadlyRated[fid])})
				}
				table.Render()
				table.ClearRows()
			}

			if aname, ok := imdb.Aid2a[input]; ok { // actor ID query
				fmt.Println("Actor:", aname)
				fmt.Println("ID: ", input)
				fmt.Println("Value:", imdb.AVor[input])
				fmt.Println("Average Rating:", imdb.ARating[input])
				fmt.Println("Filmography:")
				for fid := range imdb.Aid2Fids[input] {
					table.Append([]string{fid, imdb.Fid2f[fid], strconv.Itoa(int(imdb.Frating[fid])), strconv.FormatBool(imdb.FBroadlyRated[fid])})
				}
				table.Render()
				table.ClearRows()
			}

			if fids, ok := imdb.F2fids[input]; ok { //film name
				for fid := range fids {
					fmt.Println("film: ", input)
					fmt.Println("id: ", fid)
					fmt.Println("rated:", imdb.Frating[fid])
					fmt.Println("Broadly Rated?", imdb.FBroadlyRated[fid])
					fmt.Println("actors:")

					for aid := range imdb.Fid2Aids[fid] {
						table.Append([]string{aid, imdb.Aid2a[aid], fmt.Sprintf("%.2f", imdb.ARating[aid]), fmt.Sprintf("%.3f", imdb.AVor[aid])})
					}
					table.Render()
					table.ClearRows()
				}
			}

			if film, ok := imdb.Fid2f[input]; ok { // Film ID
				fmt.Println("film: ", film)
				fmt.Println("id: ", input)
				fmt.Println("rated:", imdb.Frating[input])
				fmt.Println("Broadly Rated?", imdb.FBroadlyRated[input])
				fmt.Println("actors:")
				for aid := range imdb.Fid2Aids[input] {
					table.Append([]string{aid, imdb.Aid2a[aid], fmt.Sprintf("%.2f", imdb.ARating[aid]), fmt.Sprintf("%.3f", imdb.AVor[aid])})

				}
				table.Render()
				table.ClearRows()
			}
		}
	}
}

var line = liner.NewLiner()

func main() {
	db, err := os.Open("./imdb.msgp")
	if err != nil {
		log.Println("Couldn't open db, going into dbFill mode")
		dbFill()
		os.Exit(0)
	}

	log.Println("Hello, Reading Database. Please Wait.")

	err = msgp.ReadFile(&imdb, db)
	if err != nil {
		log.Println("Couldn't parse imdb, going into dbFill Mode", err)
		dbFill()
	}

	log.Println("Marlon Brando is", imdb.A2aid["Marlon Brando"])
	a := imdb.A2aid["Marlon Brando"]
	log.Println("In films", imdb.Aid2Fids[a])

	defer line.Close()
	line.SetCtrlCAborts(true)
	line.SetCompleter(completeInput)

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	doCLI()

}
