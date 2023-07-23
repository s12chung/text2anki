export function paginate<T>(array: T[], pageSize: number, page: number): T[] {
  if (array.length === 0) return array
  if (pageSize <= 0) throw new Error("Invalid pageSize, it must be a positive integer.")
  if (page < 0) throw new Error("Invalid page, it must be a positive integer.")

  const start = page * pageSize
  const end = (page + 1) * pageSize

  if (start >= array.length || start < 0) throw new Error("Invalid page, out of array bounds.")

  return array.slice(start, end)
}

export function totalPages(array: unknown[], pageSize: number): number {
  if (pageSize <= 0) throw new Error("Invalid pageSize, it must be a positive integer.")
  return Math.ceil(array.length / pageSize)
}
