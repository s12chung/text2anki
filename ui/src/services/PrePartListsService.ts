import ApplicationService from "./ApplicationService.ts"
import { Http, requestInit } from "./Format.ts"
import { PrePartList, PrePartListEmpty } from "./models/PrePartList.ts"

interface PrePartListSignData {
  preParts: PrePartSignData[]
}
export interface PrePartSignData {
  imageExt?: string
  audioExt?: string
}

interface PreSignedHTTPRequest {
  url: string
  method: string
  signedHeader: Record<string, string[]>
}
const PreSignedHTTPRequestEmpty = Object.freeze<PreSignedHTTPRequest>({
  url: "",
  method: "",
  signedHeader: Object.freeze({}),
})

interface PrePartSignResponse {
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

interface PrePartListVerifyData {
  text: string
}
export interface PrePartListVerifyResponse {
  extractorType: string
}

const PrePartListVerifyResponseEmpty = Object.freeze<PrePartListVerifyResponse>({
  extractorType: "",
})
interface PrePartListCreateData {
  extractorType: string
  text: string
}

interface PrePartListCreateResponse {
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
