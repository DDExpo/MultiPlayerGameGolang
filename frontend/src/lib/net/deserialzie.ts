import { ClientData } from "$lib/stores/game.svelte"
import type { BlobOptions } from "buffer"

let off:    number
let view:   DataView
let buffer: ArrayBuffer

export function initReader(dataView: DataView, arrayBuffer: ArrayBuffer, offset: number = 0) {
  off    = offset
  view   = dataView
  buffer = arrayBuffer
}

export function skipBytes(newOffset: number) { off += newOffset }

export function readUint8(): number { return view.getUint8(off++) }

export function readInt16(): number {
  const val = view.getInt16(off, true)
  off += 2
  return val
}

export function readUint16(): number {
  const val = view.getUint16(off, true)
  off += 2
  return val
}

export function readUint32(): number {
  const val = view.getUint16(off, true)
  off += 4
  return val
}


export function readFloat32(): number {
  const val = view.getFloat32(off, true)
  off += 4
  return val
}

export function readString(): string {
  const len = view.getUint16(off, true)
  off += 2
  const str = new TextDecoder().decode(new Uint8Array(buffer, off, len))
  off += len
  return str
}

export function deserializePosition() {
  return [ readFloat32(), readFloat32(), readFloat32() ]
}

export function deserializeCombat(client: boolean) {

  if (client) {
    ClientData.Hp     = readInt16()
    ClientData.Kills  = readUint16()
    ClientData.Damage = readInt16()
  } else {
    return [readInt16(), readUint16(), readInt16()]
  }

}

export function deserializeWeapon(client: boolean) {
  if (client) {
    ClientData.WeaponType  = readUint8()
    ClientData.WeaponSpeed = readUint8()
    ClientData.WeaponWidth = readUint8()
    ClientData.WeaponRange = readInt16()
  } else {
    return [readUint8(), readUint8(), readUint8(), readInt16()]
  }
}
