package main

import (
	"multiplayerGame/game"
	"time"
)

func (h *Hub) RunGameLoop() {
	ticker := time.NewTicker(game.TickDuration)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		hits, deadIDs := game.TickProjectiles(now)
		for _, id := range deadIDs {
			msg := SerializeUserShootStatus(false, id)
			h.broadcast <- msg
		}

		for targetId, targetHits := range hits {
			target, targetOk := h.players[targetId]
			if !targetOk {
				continue
			}

			for _, hit := range targetHits {
				attacker, attackerOk := h.players[hit.OwnerId]
				if attackerOk {
					game.ApplyDamage(target, attacker)
				}
			}
		}

		for _, p := range h.players {
			game.ApplyInput(p)
			deltaMask := computeDeltaMask(p)

			if deltaMask != 0 {
				if p.Combat.HP <= 0 {
					h.broadcast <- SerializeUserDead(p.Meta.Username)
				}
				msg := SerializeUserCurrentState(deltaMask, p)
				h.broadcast <- msg
				updateLastSent(p, deltaMask)
			}
		}
	}
}
