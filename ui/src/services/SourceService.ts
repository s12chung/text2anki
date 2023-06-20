import ApplicationService from "./ApplicationService"

export interface Token {
  text: string
  partOfSpeech: string
  startIndex: number
  length: number
}

export interface TokenizedText {
  text: string
  translation: string
  tokens: Token[]
}

export interface Source {
  id: number
  name: string
  tokenizedTexts: TokenizedText[]
  updatedAt: Date
  createdAt: Date
}

class SourceService extends ApplicationService {
  protected pathPrefix = "/sources"

  async index(): Promise<Source[]> {
    return (await this.fetch()) as Source[]
  }

  async get(id: string): Promise<Source> {
    return (await this.fetch(`/${id}`)) as Source
  }
}

const sourceService = new SourceService()

export default sourceService
