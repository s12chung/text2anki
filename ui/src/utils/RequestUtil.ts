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
    const values = params.getAll(key as string)
    if (values.length === 1) {
      obj[key] = values[0] as T[keyof T]
    }
    obj[key] = values as T[keyof T]
  }
  return obj
}
