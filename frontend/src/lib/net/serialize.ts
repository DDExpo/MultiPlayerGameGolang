import { MsgType } from "$lib/Consts"
import { inputSeq } from "$lib/stores/game.svelte"

export function serializeChatMsg(text: string, color: string): ArrayBuffer {
  const encoder    = new TextEncoder()
  const textBytes  = encoder.encode(text)
  const colorBytes = encoder.encode(color)

  const buffer = new ArrayBuffer(
    1 +
    2 + textBytes.length +
    1 + colorBytes.length
  )
  const view = new DataView(buffer)
  let offset = 0

  view.setUint8(offset, MsgType.USER_CHAT)
  offset++
  view.setUint16(offset, textBytes.length, true)
  offset += 2
  new Uint8Array(buffer, offset, textBytes.length).set(textBytes)
  offset += textBytes.length

  view.setUint8(offset++, colorBytes.length)
  new Uint8Array(buffer, offset, colorBytes.length).set(colorBytes)
  offset += colorBytes.length

  return buffer
}

export function serializeInputMsg(dx: number, dy: number, isDash: boolean, angle: number) {
  const buffer = new ArrayBuffer(1 + 1 + 1 + 4 + 1)
  const view = new DataView(buffer)
  let offset = 0

  view.setUint8(offset++, MsgType.USER_INPUT)
  view.setInt8(offset++, dx)
  view.setInt8(offset++, dy)
  view.setFloat32(offset, angle, true)
  offset += 4  
  view.setUint8(offset++, isDash ? 1 : 0)

  return buffer
}
