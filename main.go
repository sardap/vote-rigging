package voterigging

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type voteType string

const (
	voteDown voteType = voteType("👎")
	voteUp            = voteType("👍")
	reactB            = voteType("🅱️")
)

var (
	voteInfav *regexp.Regexp
	voteInOp  *regexp.Regexp
)

func init() {
	voteInfav = createWordRegex(
		"734509708618235926", "735402018956378152", "734052391954939904",
		"734731825993482312", "736822440793210890", "736621603206725723",
		"735563866536280134", "738279974926024734", "219332237424984064",
		"734205649180688594",
	)
	voteInOp = createWordRegex(
	// "richard", "mitchell", "158496062103879681", "732423316081737839",
	// "278856465895129088", "348349690661699594", "302032635964948482",
	)
}

func createWordRegex(input ...string) *regexp.Regexp {
	var builder strings.Builder
	for i, word := range input {
		if i == 0 {
			fmt.Fprintf(&builder, "(%s)", regexp.QuoteMeta(word))
		}
		fmt.Fprintf(&builder, "|(%s)", regexp.QuoteMeta(word))
	}

	return regexp.MustCompile(strings.Replace(builder.String(), "\r", "", -1))
}

func handleMessage(s *discordgo.Session, mID, cID, uID string, message []byte) {
	reactions := make([]voteType, 0)
	if voteInOp.Match([]byte(uID)) {
		reactions = append(reactions, voteDown)
	} else if voteInfav.Match([]byte(uID)) {
		reactions = append(reactions, voteUp)
		reactions = append(reactions, reactB)
	}

	if len(reactions) == 0 {
		if voteInOp.Match(message) {
			reactions = append(reactions, voteDown)
		} else if voteInfav.Match(message) {
			reactions = append(reactions, voteUp)
			reactions = append(reactions, reactB)
		}
	}

	for _, vote := range reactions {
		s.MessageReactionAdd(
			cID, mID, string(vote),
		)
	}
}

//VoteReactCreateMessage Attach this to your message create and message edit
func VoteReactCreateMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == "219332237424984064" && regexp.MustCompile("vote good \\d+").Match([]byte(m.Content)) {
		mID := strings.Split(m.Content, " ")[2]
		reactions := make([]voteType, 0)
		reactions = append(reactions, voteUp)
		reactions = append(reactions, reactB)
		for _, vote := range reactions {
			s.MessageReactionAdd(
				m.ChannelID, mID, string(vote),
			)
		}
	} else if m.Author.ID == "219332237424984064" && regexp.MustCompile("vote bad \\d+").Match([]byte(m.Content)) {
		mID := strings.Split(m.Content, " ")[2]
		reactions := make([]voteType, 0)
		reactions = append(reactions, voteDown)
		for _, vote := range reactions {
			s.MessageReactionAdd(
				m.ChannelID, mID, string(vote),
			)
		}
	}
}

//VoteReactUpdateMessage Attach this to your message update and message edit
func VoteReactUpdateMessage(s *discordgo.Session, m *discordgo.MessageUpdate) {
}
