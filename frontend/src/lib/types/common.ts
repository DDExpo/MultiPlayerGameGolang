type Username = string

type MovementState = [x: number, y: number, angle: number];
type CombatState   = { dead: boolean };

type PlayerState = {
    movement: MovementState;
    combat: CombatState;
};

type PlayersState = Record<Username, PlayerState>

type ProjectileSpawn = [username: string, x: number, y: number, angle: number, ww: number, wr: number]


