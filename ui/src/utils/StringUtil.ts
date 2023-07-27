export function snakeToCamel(str: string): string {
  return str.replace(/_[a-z]/gu, (group) => group.charAt(1).toUpperCase())
}

export function camelToSnake(str: string): string {
  return str.replace(/[A-Z]/gu, (group) => `_${group.toLowerCase()}`)
}

export function camelToTitle(str: string): string {
  str = str.replace(/[A-Z]/gu, (word) => ` ${word}`)
  return str.charAt(0).toUpperCase() + str.slice(1)
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
