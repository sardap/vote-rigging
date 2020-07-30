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
	voteUp   = voteType("👍")
	reactB   = voteType("🅱️")
)

var (
	voteInfav *regexp.Regexp
	voteInOp  *regexp.Regexp
)

func init() {
	voteInfav = createWordRegex(
		"734509708618235926", "735402018956378152", "734052391954939904",
		"734731825993482312", "736822440793210890", "736621603206725723",
		"735563866536280134", "738279974926024734",
	)
	voteInOp = createWordRegex(
		"richard", "mitchell", "158496062103879681",
		"278856465895129088", "348349690661699594", "302032635964948482",
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
	if uID == "219332237424984064" {
		reactions = append(reactions, voteUp)
	} else if voteInOp.Match([]byte(uID)) {
		reactions = append(reactions, voteUp)
		reactions = append(reactions, reactB)
	} else if voteInfav.Match([]byte(uID)) {
		reactions = append(reactions, voteDown)
	}

	if len(reactions) == 0 {
		if voteInOp.Match(message) {
			reactions = append(reactions, voteUp)
			reactions = append(reactions, reactB)
		} else if voteInfav.Match(message) {
			reactions = append(reactions, voteDown)
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
	handleMessage(s, m.ID, m.ChannelID, m.Author.ID, []byte(strings.ToLower(m.Content)))
}

//VoteReactUpdateMessage Attach this to your message update and message edit
func VoteReactUpdateMessage(s *discordgo.Session, m *discordgo.MessageUpdate) {
	handleMessage(s, m.ID, m.ChannelID, m.Author.ID, []byte(strings.ToLower(m.Content)))
}
