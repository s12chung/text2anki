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

export interface TermsSearchData {
  query: string
  pos: string | string[]
}

class TermsService extends ApplicationService {
  protected pathPrefix = "/terms"

  async search(data: TermsSearchData): Promise<Term[]> {
    return (await this.fetch(`/search?${queryString(data)}`)) as Term[]
  }
}

export const termsService = new TermsService()
