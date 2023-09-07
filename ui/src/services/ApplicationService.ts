import { camelToSnake, convertKeys, snakeToCamel } from "../utils/StringUtil.ts"
import { responseError } from "./Format.ts"

abstract class ApplicationService {
  protected apiUrl = "http://localhost:3000"
  protected pathPrefix = "/"

  protected async fetch(path?: string, init?: RequestInit): Promise<unknown> {
    if (!path) path = ""

    const response = await fetch(this.apiUrl + this.pathPrefix + path, init)
    if (!response.ok) {
      // eslint-disable-next-line @typescript-eslint/no-throw-literal
      throw await responseError(response)
    }
    return convertKeys(await response.json(), snakeToCamel)
  }
}

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

export default ApplicationService
