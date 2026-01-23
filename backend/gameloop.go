package main

import (
	"multiplayerGame/game"
	"time"
)

func (h *Hub) RunGameLoop() {
	ticker := time.NewTicker(game.TickDuration)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-h.gameCmd:
			switch cmd := cmd.(type) {

			case UserResumeSession:
				{
					h.broadcast <- SerializeUserReg(cmd.Player)
				}
			case UserRegistrationCmd:
				{
					h.broadcast <- SerializeUserReg(cmd.Player)
				}
			case UserResumedDeathCmd:
				{
					game.ResetStats(cmd.Player)
					h.broadcast <- SerializeUserReg(cmd.Player)
				}
			case UserInputCmd:
				{
					cmd.Player.Input = cmd.Input
				}
			case SpawnProjectileCmd:
				{
					projectile := game.CreateProjectile(cmd.Player)
					game.AddProjectile(projectile)
					msg := SerializeUserPressedShoot(cmd.Player, projectile.ProjectileId)
					h.broadcast <- msg
				}
			}

		case now := <-ticker.C:
			hits, deadIDs := game.TickProjectiles(now)

			for _, id := range deadIDs {
				h.broadcast <- SerializeUserShootStatus(false, id)
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
}
