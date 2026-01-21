
export const WORLD_WIDTH  = 8000
export const WORLD_HEIGHT = 8000

export const PROJECTILE_SPEED = 358
export const PROJECTILE_LIFE  = 5.0 // seconds

export const PADDING_USERNAME = 36

export const MAX_USERNAME_LENGTH = 128
export const MAX_MESSAGE_LENGTH  = 512
export const MESSAGE_COOLDOWN_MS = 15_000


export const MsgType = {
  USER_CHAT:           1,
  USER_STATE:          2,
  USER_SHOOT_STATUS:   3,
  USER_RESUMED_DEATH:  4,
}

export const StateType = {
  USER_REG:           5,
  USER_DEAD:          6,
  USER_INPUT:         7,
  USER_RESUME:        8,
  USER_CURRENT_STATE: 9,
  USER_PRESSED_SHOOT: 10
}