export function paginate<T>(array: T[], maxPageSize: number, pageIndex: number): T[] {
  if (array.length === 0) return array
  if (maxPageSize <= 0) throw new Error("Invalid pageSize, it must be a positive integer.")
  if (pageIndex < 0) throw new Error("Invalid pageIndex, it must be a positive integer.")

  const start = pageIndex * maxPageSize
  const end = (pageIndex + 1) * maxPageSize

  if (start >= array.length) throw new Error("Invalid pageIndex, out of array bounds.")

  return array.slice(start, end)
}

export function pageSize(length: number, maxPageSize: number, pageIndex: number): number {
  if (length === 0) return 0
  if (maxPageSize <= 0) throw new Error("Invalid pageSize, it must be > 0.")
  if (pageIndex < 0) throw new Error("Invalid pageIndex, it must be a positive integer.")

  const start = pageIndex * maxPageSize
  const end = Math.min((pageIndex + 1) * maxPageSize, length)

  if (start >= length) throw new Error("Invalid pageIndex, out of array bounds.")

  return end - start
}

export function totalPages(array: unknown[], pageSize: number): number {
  if (pageSize <= 0) throw new Error("Invalid pageSize, it must be a positive integer.")
  return Math.ceil(array.length / pageSize)
}
