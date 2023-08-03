export function increment(index: number, length: number): number {
  if (index === -1) return 0
  return index < length - 1 ? index + 1 : 0
}

export function decrement(index: number, length: number): number {
  if (index === -1) return 0
  return index > 0 ? index - 1 : length - 1
}
