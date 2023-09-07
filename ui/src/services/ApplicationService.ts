import { convertResponse, responseError } from "./Format.ts"

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
    return convertResponse(response)
  }
}

export default ApplicationService
