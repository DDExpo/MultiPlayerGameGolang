export function randomBrightColor(): string {
  const r = 128 + Math.floor(Math.random() * 128)
  const g = 128 + Math.floor(Math.random() * 128)
  const b = 128 + Math.floor(Math.random() * 128)

  return '#' + ((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')
}