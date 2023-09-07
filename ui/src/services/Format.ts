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

export type RecursiveObj<T> = {
  [K in keyof T]: RecursiveObj<T[K]> | RecursiveObj<T[K]>[]
}

export async function convertResponse<U extends RecursiveObj<unknown>>(
  response: Response,
  empty: U
): Promise<U> {
  const data = await response.json()
  if (typeof data !== "object" || data === null) {
    throw new Error("convertResponse data is not an object or null")
  }
  return convertData(data, empty, snakeToCamel) as U
}

export function convertData(data: unknown, empty: unknown, convertKey: ConvertKeyFunc): unknown {
  if (typeof data !== "object" || data === null) {
    return data
  }
  if (Array.isArray(data)) {
    if (!Array.isArray(empty)) throw new Error(`data is array, but empty is not: ${String(empty)}`)
    return convertArray(data, empty[0] as RecursiveObj<unknown>, convertKey)
  }
  return convertObject(data, empty as RecursiveObj<unknown>, convertKey)
}

export function convertArray<T extends RecursiveObj<unknown>, U extends RecursiveObj<unknown>>(
  data: T[],
  empty: U,
  convertKey: ConvertKeyFunc
): U[] {
  return data.map((item) => convertObject(item, empty, convertKey))
}

export function convertObject<T extends RecursiveObj<unknown>, U extends RecursiveObj<unknown>>(
  data: T,
  empty: U,
  convertKey: ConvertKeyFunc
): U {
  const obj = { ...empty }
  for (const key in data) {
    if (!Object.hasOwn(data, key)) continue

    const convertedKey = convertKey(key) as keyof U
    if (!Object.hasOwn(obj, convertedKey)) {
      throw new Error(`converted key not matching object: ${String(convertedKey)}`)
    }
    obj[convertedKey] = convertData(
      data[key],
      empty[convertedKey] as RecursiveObj<unknown>,
      convertKey
    ) as U[keyof U]
  }
  return obj
}

export async function responseError(response: Response): Promise<ResponseError> {
  return new ResponseError(
    response,
    (await convertResponse(response, ResponseErrorBodyEmpty)) as ResponseErrorBody
  )
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
