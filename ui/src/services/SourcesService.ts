import ApplicationService, { Http, requestInit } from "./ApplicationService"

export interface Token {
  text: string
  partOfSpeech: string
  startIndex: number
  length: number
}

export const PosPunctuation = "Punctuation"

export interface Text {
  text: string
  translation: string
  previousBreak: boolean
}

export interface TokenizedText extends Text {
  tokens: Token[]
}

export interface SourcePart {
  media?: SourcePartMedia
  tokenizedTexts: TokenizedText[]
}

export interface SourcePartMedia {
  imageUrl?: string
  audioUrl?: string
}

export interface Source {
  id: number
  name: string
  reference: string
  parts: SourcePart[]
  updatedAt: Date
  createdAt: Date
}

export const CreateSourceDataEmpty = Object.freeze<CreateSourceData>({
  prePartListId: "",
  parts: [
    Object.freeze({
      text: "",
      translation: "",
    }),
  ],
})

export interface CreateSourceData {
  prePartListId?: string
  parts: CreateSourcePartData[]
}

export interface CreateSourcePartData {
  text: string
  translation?: string
}

export const UpdateSourceDataEmpty = Object.freeze<UpdateSourceData>({
  name: "",
  reference: "",
})

export interface UpdateSourceData {
  name: string
  reference?: string
}

class SourcesService extends ApplicationService {
  protected pathPrefix = "/sources"

  async index(): Promise<Source[]> {
    return (await this.fetch()) as Source[]
  }

  async get(id: string | number): Promise<Source> {
    return (await this.fetch(`/${id}`)) as Source
  }

  async create(data: CreateSourceData): Promise<Source> {
    return (await this.fetch("", requestInit(Http.POST, data))) as Source
  }

  async update(id: string | number, data: UpdateSourceData): Promise<Source> {
    return (await this.fetch(`/${id}`, requestInit(Http.PATCH, data))) as Source
  }
}

export const sourcesService = new SourcesService()
