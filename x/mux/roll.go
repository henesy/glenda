package mux

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"fmt"
	"strconv"
	"math/big"
	"crypto/rand"
)

// Roll dice in the form XdY
func (m *Mux) Roll(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ".\n"
	input := dm.Content
	usage := func(s string) {
		ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Malformed input. %s. \n\nExample usage: \n`g!roll 2d6 3d4` to roll two six-sided die and three four-sided die.", s))
	}

	explode := false
	infExplode := false
	if strings.HasSuffix(input, "!!") {
		infExplode = true
		input = input[:len(input)-2]
	} else if strings.HasSuffix(input, "!") {
		input = input[:len(input)-1]
		explode = true
	}

	fields := strings.Fields(input)
	if len(fields) < 2 {
		usage("No arguments provided")
		return
	}
	// Strip g!roll
	fields = fields[1:]

	for _, entry := range fields {
		resp += entry + ":\n"
		resp += "```\n"

		// {2, 6}
		parts := strings.Split(entry, "d")
		//fmt.Println(parts)
		if len(parts) != 2 {
			usage("Die entries must be in the form XdY where X and Y are valid, positive, integers")
			return
		}

		count, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			usage("Invalid numeric for count: " + err.Error())
			return
		}
		if count <= 0 {
			usage("Count must be > 0")
			return
		}

		sides, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			usage("Invalid numeric for sides: " + err.Error())
			return
		}
		if sides <= 0 {
			usage("Sides must be > 0")
			return
		}

		// Do the rolls
		for i := int64(0); i < count; i++ {
			bSides := big.NewInt(sides)

			roll := func() (*big.Int, error) {
					n, err := rand.Int(rand.Reader, bSides)
					if err != nil {
						usage("Could not perform rng")
						fmt.Println("rng fail →", err)
						return n, err
					}
					n.Add(n, big.NewInt(int64(1)))
					return n, err
			}

			insert := func(r string, n *big.Int) string {
				return r + fmt.Sprintf("%s	", n.String())
			}


			n, err := roll()
			if err != nil {
				return
			}
			resp = insert(resp, n)

			if explode && n.Cmp(bSides) == 0 {
				n, err := roll()
				if err != nil {
					return
				}
				resp = insert(resp, n)
			} else if infExplode && n.Cmp(bSides) == 0 {
				inf:
				for {
					n, err := roll()
					if err != nil {
						return
					}
					resp = insert(resp, n)
					
					if n.Cmp(bSides) == 0 {
						continue inf	
					} else {
						break inf
					}
				}
			}

		}
		resp += "\n```\n"
	}

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

