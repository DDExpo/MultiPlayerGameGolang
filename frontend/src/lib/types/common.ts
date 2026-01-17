type Username = string
type State = [username: string, x: number, y: number, speed: number, angle: number]

type PlayersState = Record<Username, State>
