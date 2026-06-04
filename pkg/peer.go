package pkg

import (
	"sync"

	"github.com/gotd/td/constant"
	tgpeers "github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

type peerStore struct {
	peers map[ChatID]tg.InputPeerClass
	mu    sync.Mutex
}

func newPeerStore() *peerStore {
	return &peerStore{
		peers: make(map[ChatID]tg.InputPeerClass),
	}
}

func (ps *peerStore) Put(chatID ChatID, peer tg.InputPeerClass) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.peers[chatID] = peer
}

func (ps *peerStore) Get(chatID ChatID) (tg.InputPeerClass, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	peer, ok := ps.peers[chatID]
	return peer, ok
}

func chatIDFromPeer(peer tg.PeerClass) (ChatID, bool) {
	var id constant.TDLibPeerID

	switch p := peer.(type) {
	case *tg.PeerUser:
		id.User(p.UserID)
	case *tg.PeerChat:
		id.Chat(p.ChatID)
	case *tg.PeerChannel:
		id.Channel(p.ChannelID)
	default:
		return 0, false
	}

	return ChatID(id), true
}

func inputPeerFromPeer(peer tg.PeerClass, e tg.Entities) (tg.InputPeerClass, bool) {
	switch p := peer.(type) {
	case *tg.PeerChat:
		return &tg.InputPeerChat{ChatID: p.ChatID}, true

	case *tg.PeerUser:
		u := e.Users[p.UserID]
		if u == nil {
			return nil, false
		}
		return &tg.InputPeerUser{
			UserID:     u.ID,
			AccessHash: u.AccessHash,
		}, true

	case *tg.PeerChannel:
		ch := e.Channels[p.ChannelID]
		if ch == nil {
			return nil, false
		}
		return &tg.InputPeerChannel{
			ChannelID:  ch.ID,
			AccessHash: ch.AccessHash,
		}, true

	default:
		return nil, false
	}
}

func chatTypeFromPeer(peer tgpeers.Peer) ChatType {
	switch p := peer.(type) {
	case tgpeers.User:
		return ChatTypeUser
	case tgpeers.Chat:
		return ChatTypeGroup
	case tgpeers.Channel:
		if p.IsSupergroup() {
			return ChatTypeSupergroup
		}
		return ChatTypeChannel
	default:
		return ChatTypeUnknown
	}
}
