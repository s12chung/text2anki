import { printError as _printError } from "../utils/ErrorUtil.ts"
import { camelToSnake, convertKeys, snakeToCamel } from "../utils/StringUtil.ts"

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

export async function convertResponse(response: Response): Promise<unknown> {
  return convertKeys(await response.json(), snakeToCamel)
}

export async function responseError(response: Response): Promise<ResponseError> {
  return new ResponseError(response, (await convertResponse(response)) as ResponseErrorBody)
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
