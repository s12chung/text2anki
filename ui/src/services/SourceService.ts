import ApplicationService, { Http, requestInit } from "./ApplicationService"

export interface Token {
  text: string
  partOfSpeech: string
  startIndex: number
  length: number
}

export interface Text {
  text: string
  translation: string
  lastEmptyLine: boolean
}

export interface TokenizedText extends Text {
  tokens: Token[]
}

export interface Source {
  id: number
  name: string
  tokenizedTexts: TokenizedText[]
  updatedAt: Date
  createdAt: Date
}

export interface CreateSourceData {
  text: string
  translation: string
}

export interface UpdateSourceData {
  name: string
}

class SourceService extends ApplicationService {
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

export const sourceService = new SourceService()
