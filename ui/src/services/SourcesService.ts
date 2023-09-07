import ApplicationService from "./ApplicationService"
import { Http, requestInit } from "./Format.ts"

export const PosPunctuation = "Punctuation"

export interface Token {
  text: string
  partOfSpeech: string
  startIndex: number
  length: number
}

const TokenEmpty = Object.freeze<Token>({
  text: "",
  partOfSpeech: "",
  startIndex: 0,
  length: 0,
})

export interface Text {
  text: string
  translation: string
  previousBreak: boolean
}

const TextEmpty = Object.freeze<Text>({
  text: "",
  translation: "",
  previousBreak: false,
})

export interface TokenizedText extends Text {
  tokens: Token[]
}

const TokenizedTextEmpty = Object.freeze<TokenizedText>({
  ...TextEmpty,
  tokens: [TokenEmpty],
})

export interface SourcePartMedia {
  imageUrl: string
  audioUrl: string
}

const SourcePartMediaEmpty = Object.freeze<SourcePartMedia>({
  imageUrl: "",
  audioUrl: "",
})

export interface SourcePart {
  media: SourcePartMedia
  tokenizedTexts: TokenizedText[]
}

const SourcePartEmpty = Object.freeze<SourcePart>({
  media: SourcePartMediaEmpty,
  tokenizedTexts: [TokenizedTextEmpty],
})

export interface Source {
  id: number
  name: string
  reference: string
  parts: SourcePart[]
  updatedAt: Date
  createdAt: Date
}

export const SourceEmpty = Object.freeze<Source>({
  id: 0,
  name: "",
  reference: "",
  parts: [SourcePartEmpty],
  updatedAt: new Date(0),
  createdAt: new Date(0),
})

export interface CreateSourcePartData {
  text: string
  translation: string
}

export const CreateSourcePartDataEmpty = Object.freeze<CreateSourcePartData>({
  text: "",
  translation: "",
})

export interface CreateSourceData {
  prePartListId: string
  name: string
  reference: string
  parts: CreateSourcePartData[]
}

export const CreateSourceDataEmpty = Object.freeze<CreateSourceData>({
  prePartListId: "",
  name: "",
  reference: "",
  parts: [CreateSourcePartDataEmpty],
})

export interface UpdateSourceData {
  name: string
  reference: string
}

export const UpdateSourceDataEmpty = Object.freeze<UpdateSourceData>({
  name: "",
  reference: "",
})

class SourcesService extends ApplicationService {
  protected pathPrefix = "/sources"

  async index(): Promise<Source[]> {
    return this.fetch("", [SourceEmpty])
  }

  async get(id: string | number): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty)
  }

  async create(data: CreateSourceData): Promise<Source> {
    return this.fetch("", SourceEmpty, requestInit(Http.POST, data))
  }

  async update(id: string | number, data: UpdateSourceData): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty, requestInit(Http.PATCH, data))
  }

  async destroy(id: string | number): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty, requestInit(Http.DELETE))
  }
}

export const sourcesService = new SourcesService()
