import { MsgType } from "$lib/Consts"
import { inputSeq } from "$lib/stores/game.svelte"

export function createBinaryUserMsg(username: string): ArrayBuffer {
  const encoder   = new TextEncoder()
  const userBytes = encoder.encode(username)
  const buffer = new ArrayBuffer(2 + userBytes.length)
  const view   = new DataView(buffer)

  // 1 byte type + 1 length + username
  let offset   = 0

  view.setUint8(offset++, MsgType.USER_REG)
  view.setUint8(offset++, userBytes.length)
  new Uint8Array(buffer, offset).set(userBytes)
  return buffer
}

export function createBinaryChatMsg(username: string, text: string, timestamp: string, color: string): ArrayBuffer {
  const encoder = new TextEncoder()
  const userBytes = encoder.encode(username)
  const textBytes = encoder.encode(text)
  const tsBytes = encoder.encode(timestamp)
  const colorBytes = encoder.encode(color)

  // 1(type)
  // + 1 + username
  // + 1 + text
  // + 1 + timestamp
  // + 1 + color
  const buffer = new ArrayBuffer(
    1 +
    (1 + userBytes.length) +
    (1 + textBytes.length) +
    (1 + tsBytes.length) +
    (1 + colorBytes.length)
  )

  const view = new DataView(buffer)
  let offset = 0

  view.setUint8(offset++, MsgType.USER_CHAT)

  view.setUint8(offset++, userBytes.length)
  new Uint8Array(buffer, offset, userBytes.length).set(userBytes)
  offset += userBytes.length

  view.setUint8(offset++, textBytes.length)
  new Uint8Array(buffer, offset, textBytes.length).set(textBytes)
  offset += textBytes.length

  view.setUint8(offset++, tsBytes.length)
  new Uint8Array(buffer, offset, tsBytes.length).set(tsBytes)
  offset += tsBytes.length

  view.setUint8(offset++, color.length)
  new Uint8Array(buffer, offset, color.length).set(colorBytes)
  offset += colorBytes.length

  return buffer
}

export function createInputMsg(dx: number, dy: number, isDash: boolean, angle: number) {
  const buffer = new ArrayBuffer(1 + 2 + 1 + 1 + 4 + 1)
  const view = new DataView(buffer)
  let o = 0

  view.setUint8(o++, MsgType.USER_INPUT)
  view.setUint16(o, ++inputSeq.value, true)
  o += 2
  view.setInt8(o++, dx)
  view.setInt8(o++, dy)
  view.setFloat32(o, angle, true)
  o += 4  
  view.setUint8(o++, isDash ? 1 : 0)

  return buffer
}

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
