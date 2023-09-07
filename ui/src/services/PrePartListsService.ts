import ApplicationService from "./ApplicationService.ts"
import { Http, requestInit } from "./Format.ts"

export interface PrePartListSignData {
  preParts: PrePartSignData[]
}

export interface PrePartSignData {
  imageExt?: string
  audioExt?: string
}

export interface PrePartListSignResponse {
  id: string
  preParts: PrePartSignResponse[]
}

export interface PrePartSignResponse {
  imageRequest?: PreSignedHTTPRequest
  audioRequest?: PreSignedHTTPRequest
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
  imageUrl?: string
  audioUrl?: string
}

export interface PrePartListVerifyData {
  text: string
}

export interface PrePartListVerifyResponse {
  extractorType: string
}

export interface PrePartListCreateData {
  extractorType: string
  text: string
}

export interface PrePartListCreateResponse {
  id: string
}

class PrePartListsService extends ApplicationService {
  protected pathPrefix = "/sources/pre_part_lists"

  async sign(data: PrePartListSignData): Promise<PrePartListSignResponse> {
    return (await this.fetch("/sign", requestInit(Http.POST, data))) as PrePartListSignResponse
  }

  async get(id: string | number): Promise<PrePartList> {
    return (await this.fetch(`/${id}`)) as PrePartList
  }

  async verify(data: PrePartListVerifyData): Promise<PrePartListVerifyResponse> {
    return (await this.fetch("/verify", requestInit(Http.POST, data))) as PrePartListVerifyResponse
  }

  async create(data: PrePartListCreateData): Promise<PrePartListCreateResponse> {
    return (await this.fetch("/", requestInit(Http.POST, data))) as PrePartListCreateResponse
  }
}

export const prePartListService = new PrePartListsService()
