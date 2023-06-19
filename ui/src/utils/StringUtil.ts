export function snakeToCamel(str: string): string {
  return str.replace(/_[a-z]/gu, (group) => group.charAt(1).toUpperCase())
}

type ConvertKeyFunc = (str: string) => string

export function convertKeys(data: unknown, convertKey: ConvertKeyFunc): unknown {
  if (typeof data !== "object" || data === null) {
    return data
  }
  if (Array.isArray(data)) {
    return data.map((item) => convertKeys(item, convertKey))
  }

  const typedData = data as Record<string, unknown>,
    mappedData: Record<string, unknown> = {}

  for (const key in typedData) {
    if (Object.hasOwn(typedData, key)) {
      mappedData[convertKey(key)] = convertKeys(typedData[key], convertKey)
    }
  }
  return mappedData
}
