import ApplicationService from "./ApplicationService.ts"
import { Http, requestInit } from "./Format.ts"

export interface PrePartListSignData {
  preParts: PrePartSignData[]
}

export interface PrePartSignData {
  imageExt?: string
  audioExt?: string
}

const PreSignedHTTPRequestEmpty = Object.freeze<PreSignedHTTPRequest>({
  url: "",
  method: "",
  signedHeader: Object.freeze({}),
})

export interface PreSignedHTTPRequest {
  url: string
  method: string
  signedHeader: Record<string, string[]>
}

export interface PrePartSignResponse {
  imageRequest: PreSignedHTTPRequest
  audioRequest: PreSignedHTTPRequest
}

const PrePartSignResponseEmpty = Object.freeze<PrePartSignResponse>({
  imageRequest: PreSignedHTTPRequestEmpty,
  audioRequest: PreSignedHTTPRequestEmpty,
})

interface PrePartListSignResponse {
  id: string
  preParts: PrePartSignResponse[]
}

const PrePartListSignResponseEmpty = Object.freeze<PrePartListSignResponse>({
  id: "",
  preParts: [PrePartSignResponseEmpty],
})

export interface PrePart {
  imageUrl: string
  audioUrl: string
}

const PrePartEmpty = Object.freeze<PrePart>({
  imageUrl: "",
  audioUrl: "",
})

export interface PrePartList {
  id: string
  preParts: PrePart[]
}

const PrePartListEmpty = Object.freeze<PrePartList>({
  id: "",
  preParts: [PrePartEmpty],
})

export interface PrePartListVerifyData {
  text: string
}

export interface PrePartListVerifyResponse {
  extractorType: string
}

const PrePartListVerifyResponseEmpty = Object.freeze<PrePartListVerifyResponse>({
  extractorType: "",
})

export interface PrePartListCreateData {
  extractorType: string
  text: string
}

export interface PrePartListCreateResponse {
  id: string
}

const PrePartListCreateResponseEmpty = Object.freeze<PrePartListCreateResponse>({
  id: "",
})

class PrePartListsService extends ApplicationService {
  protected pathPrefix = "/sources/pre_part_lists"

  async sign(data: PrePartListSignData): Promise<PrePartListSignResponse> {
    return this.fetch("/sign", PrePartListSignResponseEmpty, requestInit(Http.POST, data))
  }

  async get(id: string | number): Promise<PrePartList> {
    return this.fetch(`/${id}`, PrePartListEmpty)
  }

  async verify(data: PrePartListVerifyData): Promise<PrePartListVerifyResponse> {
    return this.fetch("/verify", PrePartListVerifyResponseEmpty, requestInit(Http.POST, data))
  }

  async create(data: PrePartListCreateData): Promise<PrePartListCreateResponse> {
    return this.fetch("/", PrePartListCreateResponseEmpty, requestInit(Http.POST, data))
  }
}

export const prePartListService = new PrePartListsService()
