import { convertKeys, snakeToCamel } from "./StringUtil.ts"

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

export async function responseError(response: Response): Promise<ResponseError> {
  const body = convertKeys(await response.json(), snakeToCamel) as ResponseErrorBody
  return new ResponseError(response, body)
}

export function printError(err: unknown): Error {
  let error: Error

  if (err instanceof Error) {
    error = err
  } else {
    let errorString: string
    switch (typeof err) {
      case "string":
      case "object":
        if (err instanceof ResponseError) {
          console.error(err) // eslint-disable-line no-console
          return new Error(err.userMessage())
        }
        errorString = JSON.stringify(err)
        break
      default:
        errorString = String(err)
    }
    error = new Error(errorString)
  }

  console.error(error) // eslint-disable-line no-console
  return error
}

export function printAndAlertError(err: unknown): Error {
  const error = printError(err)
  window.alert(error.message) // eslint-disable-line no-alert
  return error
}
