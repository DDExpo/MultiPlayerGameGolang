
export function readFloat32(view: DataView, offsetRef: { v: number }): number {
  const value = view.getFloat32(offsetRef.v, true)
  offsetRef.v += 4
  return value
}

export function readInt32(view: DataView, offsetRef: { v: number }): number {
  const value = view.getInt32(offsetRef.v, true)
  offsetRef.v += 4
  return value
}

export function readString( view: DataView, buffer: ArrayBuffer, offsetRef: { v: number }): string {
  const len = view.getUint8(offsetRef.v++)
  const bytes = new Uint8Array(buffer, offsetRef.v, len)
  offsetRef.v += len
  return new TextDecoder().decode(bytes)
}
