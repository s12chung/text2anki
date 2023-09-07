import { convertResponse, RecursiveObj, responseError } from "./Format.ts"

abstract class ApplicationService {
  protected apiUrl = "http://localhost:3000"
  protected pathPrefix = "/"

  protected async fetch<U extends RecursiveObj<unknown>>(
    path: string,
    empty: U,
    init?: RequestInit
  ): Promise<U> {
    const response = await fetch(this.apiUrl + this.pathPrefix + path, init)
    if (!response.ok) {
      // eslint-disable-next-line @typescript-eslint/no-throw-literal
      throw await responseError(response)
    }
    return convertResponse(response, empty)
  }
}

export default ApplicationService
