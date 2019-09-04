package directedmessage

import "github.com/DCNT-Hammer/dcnt/electionsCore/imessage"

type DirectedMessage struct {
	LeaderIdx int
	Msg       imessage.IMessage
}
