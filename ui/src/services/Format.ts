import { printError as _printError } from "../utils/ErrorUtil.ts"
import { camelToSnake, ConvertKeyFunc, convertKeys, snakeToCamel } from "../utils/StringUtil.ts"

export function requestInit<T extends { [K in keyof T]: unknown }>(
  method: Http,
  data?: T
): RequestInit {
  return {
    method,
    ...(data
      ? {
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(convertKeys(data, camelToSnake)),
        }
      : {}),
  }
}

type JSONTypes = string | number | boolean | null | EmptyObj | JSONTypes[]
type EmptyTypes = JSONTypes | Date | Date[]

type JSONObj = { [K in keyof unknown]: JSONTypes }
export type EmptyObj = { [K in keyof unknown]: EmptyTypes }

export async function convertResponse<T extends EmptyObj>(
  response: Response,
  empty: T
): Promise<T> {
  const data = await response.json()
  if (typeof data !== "object" || data === null) {
    throw new Error("convertResponse data is not an object or null")
  }
  return convertData(data as JSONObj | JSONObj[], empty, snakeToCamel)
}

export function convertData<T extends EmptyObj | EmptyObj[]>(
  data: JSONObj | JSONObj[],
  empty: T,
  convertKey: ConvertKeyFunc
): T {
  if (Array.isArray(data)) {
    if (!Array.isArray(empty)) throw new Error(`data is array, but empty is not: ${String(empty)}`)
    return convertArray(data, empty[0], convertKey) as T // T = EmptyObj[] => empty[0] = EmptyObj => returns EmptyObj[]
  }
  return convertObject(data, empty, convertKey)
}

export function convertArray<T extends EmptyObj>(
  data: JSONObj[],
  empty: T,
  convertKey: ConvertKeyFunc
): T[] {
  return data.map((item) => handleValue(item, empty, convertKey))
}

export function convertObject<T extends EmptyObj>(
  data: JSONObj,
  empty: T,
  convertKey: ConvertKeyFunc
): T {
  const obj = { ...empty }
  for (const key in data) {
    if (!Object.hasOwn(data, key)) continue

    const convertedKey = convertKey(key) as keyof T // always a key
    if (!Object.hasOwn(obj, convertedKey)) {
      throw new Error(`converted key not matching object: ${String(convertedKey)}`)
    }
    obj[convertedKey] = handleValue(
      data[key as keyof EmptyObj], // always a key
      empty[convertedKey] as JSONTypes, // T = EmptyObj => T[keyof T] => JSONTypes
      convertKey
    ) as T[keyof T] // empty[convertedKey] = T[keyof T]
  }
  return obj
}

function handleValue<T extends EmptyTypes>(
  data: JSONTypes,
  empty: T,
  convertKey: ConvertKeyFunc
): T {
  if (data === null) return empty
  if (empty === null) throw new Error("empty is null")

  if (typeof data !== "object") {
    if (empty instanceof Date && typeof data === "string") return new Date(data) as unknown as T // T = empty = Date via guard
    if (typeof data !== typeof empty)
      throw new Error(`data (${String(data)}) type not matching empty {${String(empty)})`)
    return data as T // `typeof data !== typeof empty` ensures same
  }
  if (typeof empty !== "object") throw new Error(`empty is not an object: ${String(empty)}`)
  return convertData(data, empty, convertKey)
}

export async function responseError(response: Response): Promise<ResponseError> {
  return new ResponseError(response, await convertResponse(response, ResponseErrorBodyEmpty))
}

export class ResponseError {
  public headers: Headers
  public status: number
  public statusText: string
  public url: string
  public body: ResponseErrorBody
  constructor(response: Response, body: ResponseErrorBody) {
    this.headers = response.headers
    this.status = response.status
    this.statusText = response.statusText
    this.url = response.url
    this.body = body
  }

  userMessage(): string {
    return this.body.error
  }
}

export interface ResponseErrorBody {
  error: string
  code: number
  statusText: string
}

export const ResponseErrorBodyEmpty = Object.freeze<ResponseErrorBody>({
  error: "",
  code: 0,
  statusText: "",
})

export function printError(err: unknown): Error {
  if (err instanceof ResponseError) {
    console.error(err) // eslint-disable-line no-console
    return new Error(err.userMessage())
  }
  return _printError(err)
}

export function printAndAlertError(err: unknown): Error {
  const error = printError(err)
  window.alert(error.message) // eslint-disable-line no-alert
  return error
}

export enum Http {
  GET = "GET",
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE",
  PATCH = "PATCH",
  HEAD = "HEAD",
  OPTIONS = "OPTIONS",
  CONNECT = "CONNECT",
  TRACE = "TRACE",
}
