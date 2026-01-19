
export class SpatialHash {
  private cellSize: number
  private grid: Map<string, Map<string, [x: number, y: number]>>
  private playerCells: Map<string, string> 
  
  constructor(cellSize: number = 400) {
    this.cellSize = cellSize
    this.grid = new Map()
    this.playerCells = new Map()
  }

  private hash(x: number, y: number): string {
    const cellX = Math.floor(x / this.cellSize)
    const cellY = Math.floor(y / this.cellSize)
    return `${cellX},${cellY}`
  }

  update(username: string, x: number, y: number ) {
    const newKey = this.hash(x, y)
    const oldKey = this.playerCells.get(username)
    
    if (oldKey && oldKey !== newKey) {
      const oldCell = this.grid.get(oldKey)
      if (oldCell) {
        oldCell.delete(username)
        if (oldCell.size === 0) {
          this.grid.delete(oldKey)
        }
      }
    }
    if (!this.grid.has(newKey)) {
      this.grid.set(newKey, new Map())
    }
    
    this.grid.get(newKey)!.set(username, [x, y])
    this.playerCells.set(username, newKey)
  }

  remove(username: string) {
    const key = this.playerCells.get(username)
    
    if (key) {
      const cell = this.grid.get(key)
      if (cell) {
        cell.delete(username)
        if (cell.size === 0) {
          this.grid.delete(key)
        }
      }
      this.playerCells.delete(username)
    }
  }

  getCell(x: number, y: number): Map<string, [x: number, y: number]> {
    const key = this.hash(x, y)
    return this.grid.get(key) || new Map()
  }

  clear() {
    this.grid.clear()
    this.playerCells.clear()
  }
}
export const FastProjectileCheck = new SpatialHash(400)
