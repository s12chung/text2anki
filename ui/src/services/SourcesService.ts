import ApplicationService from "./ApplicationService"
import { Http, requestInit } from "./Format.ts"
import { PartCreateMultiData, PartCreateMultiDataEmpty } from "./PartsService.ts"

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

export interface CreateSourceData extends PartCreateMultiData {
  name: string
  reference: string
}

export const CreateSourceDataEmpty = Object.freeze<CreateSourceData>({
  name: "",
  reference: "",
  ...PartCreateMultiDataEmpty,
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

export function tokenPreviousSpace(tokens: Token[], index: number): boolean {
  if (index === 0) return false
  const currentToken = tokens[index]
  const previousToken = tokens[index - 1]
  return previousToken.startIndex + previousToken.length + 1 === currentToken.startIndex
}

export function tokenPreviousPunct(tokens: Token[], index: number): boolean {
  if (index === 0) return false
  return tokens[index - 1].partOfSpeech === PosPunctuation
}
