import { MsgType } from "$lib/Consts"

export function createBinaryUserMsg(username: string): ArrayBuffer {
  const encoder   = new TextEncoder()
  const userBytes = encoder.encode(username)
  const buffer = new ArrayBuffer(2 + userBytes.length)
  const view   = new DataView(buffer)

  // 1 byte type + 1 length + username
  let offset   = 0

  view.setUint8(offset++, MsgType.USER)
  view.setUint8(offset++, userBytes.length)
  new Uint8Array(buffer, offset).set(userBytes)
  return buffer
}

export function createBinaryChatMsg(username: string, text: string, timestamp: string): ArrayBuffer {
  const encoder = new TextEncoder()
  const userBytes = encoder.encode(username)
  const t = encoder.encode(text)
  const ts = encoder.encode(timestamp)

  // 1 (type) + (1 + userBytes.length) username + (1 + t.length) text + (1 + ts.length) timestamp
  const buffer = new ArrayBuffer(1 + (1 + userBytes.length) + (1 + t.length) + (1 + ts.length))
  const view = new DataView(buffer)

  let offset = 0
  view.setUint8(offset++, MsgType.CHAT)
  view.setUint8(offset++, userBytes.length)
  new Uint8Array(buffer, offset, userBytes.length).set(userBytes)

  offset += userBytes.length
  view.setUint8(offset++, t.length)
  new Uint8Array(buffer, offset, t.length).set(t)

  offset += t.length
  view.setUint8(offset++, ts.length)
  new Uint8Array(buffer, offset, ts.length).set(ts)

  return buffer
}

export function createBinaryUserStateMsg(username: string, x: number, y: number, angle: number): ArrayBuffer {
  const encoder   = new TextEncoder()
  const userBytes = encoder.encode(username)
  const buffer = new ArrayBuffer(1 + 1 + userBytes.length + 12)
  const view   = new DataView(buffer)
  
  // 1 byte type  + 1 length + username + 4 (x) + 4 (y) + 4 (angle)
  let offset = 0

  view.setUint8(offset++, MsgType.USER_STATE)
  view.setUint8(offset++, userBytes.length)
  new Uint8Array(buffer, offset, userBytes.length).set(userBytes)
  offset += userBytes.length

  view.setFloat32(offset, x, true); offset += 4
  view.setFloat32(offset, y, true); offset += 4
  view.setFloat32(offset, angle, true)

  return buffer
}

export function readFloat32(view: DataView, offsetRef: { v: number }): number {
  const value = view.getFloat32(offsetRef.v, true)
  offsetRef.v += 4
  return value
}

export function readString( view: DataView, buffer: ArrayBuffer, offsetRef: { v: number }): string {
  const len = view.getUint8(offsetRef.v++)
  const bytes = new Uint8Array(buffer, offsetRef.v, len)
  offsetRef.v += len
  return new TextDecoder().decode(bytes)
}