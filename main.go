package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("LINK_SHORTENER_DISCORD_TOKEN"))
	err = dg.Open()
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}
	dg.AddHandler(onMessageCreate)
	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stopBot
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	d := amazonLinkVerifier(m.Content)
	if d == nil {
		return
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       d.Title,
		Description: d.Url,
		URL:         d.Url,
		Type:        discordgo.EmbedTypeLink,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Amazonリンクが共有されました",
		},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

type linkInfo struct {
	Title string
	Url   string
}

func amazonLinkVerifier(link string) *linkInfo {
	r := regexp.MustCompile(`^https://www\.amazon\.co\.jp/?([^　 ])*/dp/(?P<id>[A-Z\d]{10})`)
	//second parameter of FindAllStringSubmatch is the maximum size of result
	result := r.FindAllStringSubmatch(link, 1)
	if len(result) <= 0 {
		return nil
	}
	outLink := "https://www.amazon.co.jp/dp/"
	outLink = result[0][1] + outLink

	title := getOGP(outLink)
	return &linkInfo{title, outLink}
}

func getOGP(link string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", link, nil)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var title string
	if t, ok := GetHtmlTitle(resp.Body); ok {
		title = t
	} else {
		println("Fail to get HTML title")
	}
	return title
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", false
	}
	return traverse(doc)
}
