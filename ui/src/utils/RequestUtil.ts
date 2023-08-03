export function headers<T extends Record<keyof T, string | string[]>>(obj: T): Headers {
  return setAppendAble(new Headers(), obj)
}

export function queryString<T extends Record<keyof T, string | string[]>>(queryParams: T): string {
  return setAppendAble(new URLSearchParams(), queryParams).toString()
}

export function queryObject<T extends Record<keyof T, string[] | string | number | boolean>>(
  url: string,
  empty: T
) {
  const params = new URL(url).searchParams

  const obj = { ...empty }
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) continue

    const values = params.getAll(key)
    const objValue = obj[key]

    obj[key] = (
      Array.isArray(objValue) ? values : stringParsedValue(objValue, values[0])
    ) as T[Extract<keyof T, string>]
  }
  return obj
}

type FormTypes = FormDataEntryValue | number | boolean

export function formData<
  T extends Record<keyof T, FormTypes | U[]>,
  U extends Record<keyof U, FormTypes>
>(formData: FormData, empty: T): T {
  const formKeys = Array.from(formData.keys()).sort()

  const obj = { ...empty }
  for (const key in obj) {
    if (!Object.hasOwn(obj, key)) continue

    const objValue = obj[key]
    let value
    if (Array.isArray(objValue)) {
      value = objFromArrayKeys(objValue[0], key, formKeys, formData)
    } else {
      const v = formData.get(key)
      if (v === null) continue
      value = formValue(objValue, key, v)
    }
    obj[key] = value as T[Extract<keyof T, string>]
  }

  return obj
}

interface Appendable {
  append(name: string, value: string): void
}

function setAppendAble<T extends Appendable, U extends Record<keyof U, string | string[]>>(
  appendable: T,
  obj: U
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

function formValue<T extends FormTypes>(
  objValue: T,
  key: string,
  value: FormDataEntryValue
): FormTypes {
  if (objValue instanceof File || value instanceof File) {
    if (!(objValue instanceof File))
      throw new Error(`key, ${key}, is not an instance of File, but the formData at key is`)
    if (!(value instanceof File))
      throw new Error(`formData key, ${key}, is not an instance of File, but the obj at key is`)
    return value
  }
  return stringParsedValue(objValue, value) // guaranteed string parsable due to above
}

// eslint-disable-next-line consistent-return
function stringParsedValue<T extends string | number | boolean>(
  objValue: T,
  value: string
): string | number | boolean {
  // never hits default due to typing
  // eslint-disable-next-line default-case
  switch (typeof objValue) {
    case "string":
      return value
    case "number":
      return parseInt(value, 10)
    case "boolean":
      return value === "true"
  }
}

// eslint-disable-next-line max-params
function objFromArrayKeys<T extends Record<keyof T, FormTypes>>(
  empty: T,
  key: string,
  formKeys: string[],
  formData: FormData
): T[] {
  const keys = formKeys.filter((key) => key.startsWith(key))

  const array: T[] = []
  for (let i = 0; ; i++) {
    const prefix = `${key}[${i}]`
    if (keys.filter((key) => key.startsWith(prefix)).length === 0) break

    const obj = { ...empty }
    for (const objKey in obj) {
      if (!Object.hasOwn(obj, objKey)) continue

      const value = formData.get(`${prefix}.${objKey}`)
      if (value === null) continue
      obj[objKey] = formValue(obj[objKey], objKey, value) as T[Extract<keyof T, string>]
    }
    array.push(obj)
  }
  return array
}
