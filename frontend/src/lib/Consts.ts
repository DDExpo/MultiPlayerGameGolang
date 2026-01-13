export const MAX_USERNAME_LENGTH = 100
export const MAX_MESSAGE_LENGTH  = 1024
export const MESSAGE_COOLDOWN_MS = 15_000

export const MsgType = {
  USER: 1,
  CHAT: 2,
  USER_STATE: 3,
} as const
