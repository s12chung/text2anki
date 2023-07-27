export function formData<T extends Record<keyof T, FormDataEntryValue | number | boolean>>(
  formData: FormData,
  empty: T
): T {
  const obj = { ...empty }
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) continue

    const value = formData.get(key)
    if (value === null) continue

    if (value instanceof File) {
      if (!(obj[key] instanceof File))
        throw new Error(`key, ${key}, is not an instance of File, but the formData at key is`)
      obj[key] = value as T[Extract<keyof T, string>]
      continue
    }
    setStringParsable(obj as Record<keyof T, string | number | boolean>, key, value) // guaranteed string parsable due to above
  }

  return obj
}

export function queryString<T extends Record<keyof T, string | string[]>>(queryParams: T): string {
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

export function queryObject<T extends Record<keyof T, string | string[] | number | boolean>>(
  url: string,
  empty: T
) {
  const params = new URL(url).searchParams

  const obj = { ...empty }
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) continue

    const values = params.getAll(key as string)
    if (Array.isArray(obj[key])) {
      obj[key] = values as T[Extract<keyof T, string>]
      continue
    }
    // eslint-disable-next-line prefer-destructuring
    setStringParsable(obj as Record<keyof T, string | number | boolean>, key, values[0]) // guaranteed string parsable due to above
  }
  return obj
}

function setStringParsable<T extends Record<keyof T, string | number | boolean>>(
  obj: T,
  key: keyof T,
  value: string
) {
  switch (typeof obj[key]) {
    case "string":
      obj[key] = value as T[Extract<keyof T, string>]
      return
    case "number":
      obj[key] = parseInt(value, 10) as T[Extract<keyof T, string>]
      return
    case "boolean":
      obj[key] = (value === "true") as T[Extract<keyof T, string>]
      return
    default:
      throw new Error(
        `key, ${String(key)}, is not a valid type ${obj[key].toString()} (${typeof obj[key]})`
      )
  }
}
