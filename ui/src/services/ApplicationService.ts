import { convertKeys, snakeToCamel } from "../utils/StringUtil.ts"

abstract class ApplicationService {
  protected apiUrl = "http://localhost:3000"
  protected pathPrefix = "/"

  protected async fetch(path?: string, init?: RequestInit): Promise<unknown> {
    if (!path) path = ""

    const response = await fetch(this.apiUrl + this.pathPrefix + path, init)
    if (!response.ok) {
      throw new Error(`Failed to fetch: ${path}`)
    }
    return convertKeys(await response.json(), snakeToCamel)
  }
}

export default ApplicationService
