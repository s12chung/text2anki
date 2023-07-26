export function formData<T extends Record<keyof T, FormDataEntryValue>>(
  formData: FormData,
  ...keys: (keyof T)[]
): T {
  const obj = {} as T
  for (const key of keys) {
    obj[key] = formData.get(key as string) as T[keyof T]
  }
  return obj
}

export function queryString<T extends Record<K, string | string[]>, K extends keyof T>(
  queryParams: T
): string {
  const params = new URLSearchParams()
  for (const key in queryParams) {
    if (!Object.hasOwn(queryParams, key)) continue
    const values = queryParams[key]

    if (typeof values === "string") {
      params.append(key, values)
      continue
    }

    for (const value of values) {
      params.append(key, value)
    }
  }
  return params.toString()
}

export function queryObject<T extends Record<keyof T, string | string[]>>(
  url: string,
  ...keys: (keyof T)[]
) {
  const params = new URL(url).searchParams
  const obj = {} as T
  for (const key of keys) {
    obj[key] = params.get(key as string) as T[keyof T]
  }
  return obj
}
