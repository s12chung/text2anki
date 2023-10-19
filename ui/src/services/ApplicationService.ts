import { convertResponse, EmptyObj, responseError } from "./Format.ts"

abstract class ApplicationService {
  protected apiUrl = "http://localhost:3000"
  protected pathPrefix = "/"

  protected pathUrl(path: string): string {
    return this.apiUrl + this.pathPrefix + path
  }

  protected async fetch<T extends EmptyObj>(
    path: string,
    empty: T,
    init?: RequestInit,
  ): Promise<T> {
    const response = await fetch(this.pathUrl(path), init)
    if (!response.ok) {
      // eslint-disable-next-line @typescript-eslint/no-throw-literal
      throw await responseError(response)
    }
    return convertResponse(response, empty)
  }
}

export default ApplicationService
