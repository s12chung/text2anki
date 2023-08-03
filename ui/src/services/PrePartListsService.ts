import { queryString } from "../utils/RequestUtil.ts"
import ApplicationService from "./ApplicationService.ts"

export interface PrePartListSignResponse {
  id: string
  requests: PreSignedHTTPRequest[]
}

export interface PreSignedHTTPRequest {
  url: string
  method: string
  signedHeader: Record<string, string[]>
}

export interface PrePartList {
  id: string
  preParts: PrePart[]
}

export interface PrePart {
  url: string
}

class PrePartListsService extends ApplicationService {
  protected pathPrefix = "/sources/pre_part_lists"

  async sign(exts: string[]): Promise<PrePartListSignResponse> {
    return (await this.fetch(`/sign?${queryString({ exts })}`)) as PrePartListSignResponse
  }

  async get(id: string | number): Promise<PrePartList> {
    return (await this.fetch(`/${id}`)) as PrePartList
  }
}

export const prePartListService = new PrePartListsService()
