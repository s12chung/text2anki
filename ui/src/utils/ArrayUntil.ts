export function unique<T>(items: T[]): T[] {
  return [...new Set(items)] // will keep order
}

export function filterKeys<K extends string | number | symbol>(
  original: K[],
  ...others: K[][]
): K[] {
  const otherMap: { [key in K]?: boolean } = {}
  for (const other of others) {
    for (const key of other) {
      otherMap[key] = true
    }
  }
  return original.filter((val): boolean => !otherMap[val])
}
