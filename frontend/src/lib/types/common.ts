type WSMessage<T = any> = {
  type: string
  data: T
}

type PlayerMoveMessage = {
  username: string
  x:        number
  y:        number
  angle:    number
}


type Username = string
type Position = [x: number, y: number, angle: number]

type PlayersState = Record<Username, Position>
