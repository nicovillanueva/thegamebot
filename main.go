package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var botToken = os.Getenv("BOT_TOKEN")
var forceGroups = os.Getenv("PREADD_GROUPS")
var recipients = []int64{}

var maxDelta = 48
var timeScale = time.Hour
var games = []string{
	"the game",
	"you lost the game",
	"lost the gaem",
	"the game, you lost",
	"game the lost you",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := forceGroups
	if g != "" {
		log.Printf("[startup] pre-adding groups: %s", g)
		for _, e := range strings.Split(g, ",") {
			i, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				log.Printf("[startup] skipping invalid group ID: %s", e)
				continue
			}
			registerGroup(i)
		}
	}
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panicln(err)
	}
	b.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		log.Printf("[debog] handling groupadd to %d", m.Chat.ID)
		if m.FromGroup() {
			registerGroup(m.Chat.ID)
		}
	})
	go b.Start()
	log.Printf("[startup] started")
	greet(b)

	next := getRandomHour()
	a := time.After(time.Duration(next) * timeScale)
	log.Printf("[startup] will run in %d hours", next)
	for {
		select {
		case <-a:
			for _, g := range recipients {
				log.Printf("[debog] spamming %d", g)
				c := tb.Chat{
					ID: g,
				}
				b.Send(&c, games[rand.Intn(len(games))])
			}
			a = reschedule()
		}
	}
}

func getRandomHour() int {
	n := rand.Intn(maxDelta)
	if n == 0 {
		n = 2
	}
	return n
}

func reschedule() (a <-chan time.Time) {
	h := getRandomHour()
	log.Printf("[reschedule] will run in %d", h)
	return time.After(time.Duration(h) * timeScale)
}

func alreadyRegistered(gid int64) bool {
	for _, g := range recipients {
		if gid == g {
			return true
		}
	}
	return false
}

func registerGroup(gid int64) bool {
	if !alreadyRegistered(gid) {
		recipients = append(recipients, gid)
		log.Printf("[debog] registered target group %d\n", gid)
		return true
	}
	log.Printf("[debog] not re-registering group %d\n", gid)
	return false
}

func greet(b *tb.Bot) {
	for _, g := range recipients {
		c := tb.Chat{ID: g}
		b.Send(&c, "hello motherfuckers")
	}
}
