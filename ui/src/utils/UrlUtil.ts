export function queryString(queryParams: Record<string, string[]>): string {
  const params = new URLSearchParams()
  for (const key in queryParams) {
    if (!Object.hasOwn(queryParams, key)) continue
    const values = queryParams[key]
    for (const value of values) {
      params.append(key, value)
    }
  }
  return params.toString()
}
