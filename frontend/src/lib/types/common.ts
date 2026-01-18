type Username = string
type State = [x: number, y: number, speed: number, angle: number]

type PlayersState = Record<Username, State>
