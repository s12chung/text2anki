export function headers<T extends Record<keyof T, string | string[]>>(obj: T): Headers {
  return setAppendAble(new Headers(), obj)
}

export function queryString<T extends Record<keyof T, string | string[]>>(queryParams: T): string {
  return setAppendAble(new URLSearchParams(), queryParams).toString()
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

interface Appendable {
  append(name: string, value: string): void
}

function setAppendAble<T extends Appendable, S extends Record<keyof S, string | string[]>>(
  appendable: T,
  obj: S
): T {
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) continue
    const values = obj[key]

    if (typeof values === "string") {
      appendable.append(key, values)
      continue
    }
    for (const value of values) {
      appendable.append(key, value)
    }
  }
  return appendable
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
