import { queryString } from "../utils/UrlUtil.ts"
import ApplicationService from "./ApplicationService.ts"

export interface Term {
  id: number
  text: string
  variants: string[]
  partOfSpeech: string
  commonLevel: number
  translations: Translation[]
}

export interface Translation {
  text: string
  explanation: string
}

class TermsService extends ApplicationService {
  protected pathPrefix = "/terms"

  async search(query: string, pos: string): Promise<Term[]> {
    return (await this.fetch(`/search?${queryString({ query: [query], pos: [pos] })}`)) as Term[]
  }
}

export const termsService = new TermsService()
