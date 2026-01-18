
export enum UserStateDelta {
  POS    = 1 << 0, // X, Y, speed, angle
  STATS  = 1 << 1,
  WEAPON = 1 << 2
}

export enum HttpStatus {
    OK                    = 200,
    CREATED               = 201,
    BAD_REQUEST           = 400,
    UNAUTHORIZED          = 401,
    FORBIDDEN             = 403,
    NOT_FOUND             = 404,
    CONFLICT              = 409,
    INTERNAL_SERVER_ERROR = 500,
}
