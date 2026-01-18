import { MsgType } from "$lib/Consts"
import { inputSeq } from "$lib/stores/game.svelte"

export function serializeChatMsg(text: string, color: string): ArrayBuffer {
  const encoder    = new TextEncoder()
  const textBytes  = encoder.encode(text)
  const colorBytes = encoder.encode(color)

  // 1(type) + 1 + text + 1 + color
  const buffer = new ArrayBuffer(
    1 +
    (1 + textBytes.length) +
    (1 + colorBytes.length)
  )

  const view = new DataView(buffer)
  let offset = 0

  view.setUint8(offset++, MsgType.USER_CHAT)

  view.setUint8(offset++, textBytes.length)
  new Uint8Array(buffer, offset, textBytes.length).set(textBytes)
  offset += textBytes.length

  view.setUint8(offset++, colorBytes.length)
  new Uint8Array(buffer, offset, colorBytes.length).set(colorBytes)
  offset += colorBytes.length

  return buffer
}

export function serializeInputMsg(dx: number, dy: number, isDash: boolean, angle: number) {
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
