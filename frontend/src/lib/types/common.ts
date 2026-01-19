type Username = string
type State = [x: number, y: number, angle: number]
type Projectile = [x: number, y: number, angle: number, damage: number, weaponWidth: number, weaponRange: number]

type PlayersState = Record<Username, State>

type OtherProjectiles = Record<Username, Projectile>
