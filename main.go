package main

import (
	"time"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
    "math/rand"
    "os"
)

var botToken = os.Getenv("BOT_TOKEN")
var recipients = []int64{}
func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panicln(err)
	}
    b.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
        log.Printf("handling groupadd to %d", m.Chat.ID)
        if m.FromGroup() && ! alreadyRegistered(m.Chat.ID) {
            recipients = append(recipients, m.Chat.ID)
            log.Printf("Registered target group %d\n", m.Chat.ID)
        }
    })
    go b.Start()
    log.Printf("started")

    var hasRunOn = -1
    next := getRandomHour()
    a := time.After(time.Duration(next) * time.Second)
    log.Printf("[startup] Will run in %d hours", next)
    for {
        select {
        case <- a:
            if hasRunOn != time.Now().Day() {
                for _, g := range recipients {
                    c := tb.Chat {
                        ID: g,
                    }
                    b.Send(&c, "the game")
                }
                hasRunOn = time.Now().Day()
                a = reschedule()
            } else {
                log.Printf("already ran today")
            }
        }
    }
}

func getRandomHour() int {
    n := rand.Intn(24)
    if n == 0 {
        n = 1
    }
    return n
}

func reschedule() (a <-chan time.Time) {
    h := getRandomHour()
    log.Printf("will run at %d", h)
    return time.After(time.Duration(h) * time.Hour)
}

func alreadyRegistered(gid int64) bool {
    for _, g := range recipients {
        if gid == g {
            return true
        }
    }
    return false
}
