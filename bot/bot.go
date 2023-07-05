package bot

import (
	"encoding/json"
	"fmt"
	"go-bot/config"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

var BotID string

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)
	goBot.User("radgurlgamer420#5684")

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running...")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content[:1] == config.BotPrefix {
		cmd := strings.Split(m.Content, " ")[0][1:]
		switch cmd {
		case "ping":
			handlePing(s, m.ChannelID)
		case "pyth":
			handlePyth(s, m.ChannelID, cmd, m.Content)
		case "rps":
			handleRPS(s, m.ChannelID, cmd, m.Content)
		case "insult":
			handleInsult(s, m.ChannelID)
		case "dog":
			handleDog(s, m.ChannelID)
		case "mineMe":
			handleMineMe(s, m.ChannelID, cmd, m.Content)
		case "help":
			handleHelp(s, m.ChannelID)
		default:
			handleDefault(s, m.ChannelID, cmd)
		}
	}
}

func handleHelp(s *discordgo.Session, channelID string) {
	rawCmds := []string{"ping", "pyth", "rps", "insult", "dog", "mineMe"}
	fmtdCmds := make([]string, len(rawCmds))
	for i := 0; i < len(rawCmds); i++ {
		fmtdCmds[i] = "`!" + rawCmds[i] + "`"
	}
	_, err := s.ChannelMessageSend(channelID, "here's a list of commands:\n"+strings.Join(fmtdCmds, "\n"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func handlePing(s *discordgo.Session, channelID string) {
	_, err := s.ChannelMessageSend(channelID, "pong")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func handleDefault(s *discordgo.Session, channelID string, cmd string) {
	_, err := s.ChannelMessageSend(channelID, "**!"+cmd+" error:** \n!"+cmd+" is not a command")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func handlePyth(s *discordgo.Session, channelID string, cmd string, content string) {
	params := strings.Split(content, " ")[1:] // everything after the command
	if len(params) != 2 {
		returnError(s, channelID, cmd)
		return
	}
	aString, bString := params[0], params[1]
	aFloat, err := strconv.ParseFloat(aString, 64)
	if err != nil {
		returnError(s, channelID, cmd)
		return
	}
	bFloat, err := strconv.ParseFloat(bString, 64)
	if err != nil {
		returnError(s, channelID, cmd)
		return
	}

	cFloat := math.Sqrt(math.Pow(aFloat, 2) + math.Pow(bFloat, 2))

	_, err = s.ChannelMessageSend(channelID, fmt.Sprintf("%.2f", cFloat))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func handleRPS(s *discordgo.Session, channelID string, cmd string, content string) {
	params := strings.Split(content, " ")[1:] // everything after the command
	if len(params) != 1 {
		returnError(s, channelID, cmd)
		return
	}
	rpsChoice := strings.ToLower(params[0])
	validChoice := []string{"r", "p", "s", "rock", "paper", "scissors"}
	if slices.Contains(validChoice, rpsChoice) {
		humanChoice := rpsChoice[:1]
		botChoice := []string{"r", "p", "s"}[rand.Intn(3)]
		letterToString := map[string]string{
			"r": "Rock",
			"p": "Paper",
			"s": "Scissors",
		}
		s.ChannelMessageSend(channelID, "I chose... \n**"+letterToString[botChoice]+"!**")
		if humanChoice == botChoice {
			s.ChannelMessageSend(channelID, "*It's a tie!*")
		}
		if (humanChoice == "r" && botChoice == "s") || (humanChoice == "p" && botChoice == "r") || (humanChoice == "s" && botChoice == "p") {
			s.ChannelMessageSend(channelID, "*You win!*")
		}
		if (botChoice == "r" && humanChoice == "s") || (botChoice == "p" && humanChoice == "r") || (botChoice == "s" && humanChoice == "p") {
			s.ChannelMessageSend(channelID, "*I win!*")
		}
	} else {
		returnError(s, channelID, cmd)
	}

}

func handleInsult(s *discordgo.Session, channelID string) {
	resp, err := http.Get("https://evilinsult.com/generate_insult.php")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	s.ChannelMessageSend(channelID, string(body))
}

func handleDog(s *discordgo.Session, channelID string) {
	resp, err := http.Get("https://random.dog/woof.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	dogBody := struct {
		URL string `json:"url"`
	}{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal(body, &dogBody)

	dogEmbed := &discordgo.MessageEmbed{Image: &discordgo.MessageEmbedImage{URL: dogBody.URL}}

	s.ChannelMessageSendEmbed(channelID, dogEmbed)

}

func handleMineMe(s *discordgo.Session, channelID string, cmd string, content string) {
	params := strings.Split(content, " ")[1:] // everything after the command
	if len(params) != 1 {
		returnError(s, channelID, cmd)
		return
	}
	userName := params[0]

	resp, err := http.Get("https://playerdb.co/api/player/minecraft/" + userName)
	if err != nil {
		returnError(s, channelID, cmd)
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var mineId struct {
		Data struct {
			Player struct {
				ID string `json:"raw_id"`
			} `json:"player"`
		} `json:"data"`
	}

	json.Unmarshal(body, &mineId)

	mineMeEmbed := &discordgo.MessageEmbed{Image: &discordgo.MessageEmbedImage{URL: "https://crafatar.com/renders/body/" + mineId.Data.Player.ID}}
	s.ChannelMessageSendEmbed(channelID, mineMeEmbed)

}

func returnError(s *discordgo.Session, channelID string, cmd string) {
	switch cmd {
	case "pyth":
		s.ChannelMessageSend(channelID, "**!pyth error:** \npls provide 2 numbers: `!pyth <a> <b>`")
	case "rps":
		s.ChannelMessageSend(channelID, "**!rps error:** \npls choose rock, paper, or scissors")
	case "mineMe":
		s.ChannelMessageSend(channelID, "**!mineMe error:** \nplease enter your mindcraft username: `!mineMe <userName>`")
	default:
		s.ChannelMessageSend(channelID, "error sending message")
	}
}
