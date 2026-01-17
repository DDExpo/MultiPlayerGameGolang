
export const WORLD_WIDTH  = 8000
export const WORLD_HEIGHT = 8000

export const MAX_USERNAME_LENGTH = 100
export const MAX_MESSAGE_LENGTH  = 1024
export const MESSAGE_COOLDOWN_MS = 15_000


export const MsgType = {
  USER_REG: 1,
  USER_CHAT: 2,
  USER_STATE: 3,
  USER_RESUME: 4,
  USER_INPUT: 5
} as const

export enum UserStateDelta {
  POS    = 1 << 0, // X, Y, speed, angle
  STATS  = 1 << 1,
  WEAPON = 1 << 2
}