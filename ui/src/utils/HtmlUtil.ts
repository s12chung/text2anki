export function join(klass: string, joiner: string): string {
  const classes = [klass] as string[]
  if (joiner !== "") classes.push(joiner)
  return classes.join(" ")
}

export function joinIf(klass: string, condition: boolean, joiner: string): string {
  if (!condition) return klass
  return join(klass, joiner)
}

export function paginate<T>(array: T[], pageSize: number, page: number): T[] {
  if (array.length === 0) return array
  if (pageSize <= 0) throw new Error("Invalid pageSize, it must be a positive integer.")
  if (page <= 0) throw new Error("Invalid page, it must be a positive integer.")

  const start = (page - 1) * pageSize
  const end = page * pageSize

  if (start >= array.length || start < 0) throw new Error("Invalid page, out of array bounds.")

  return array.slice(start, end)
}
