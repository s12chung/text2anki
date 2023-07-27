import { queryString } from "../utils/RequestUtil.ts"
import ApplicationService from "./ApplicationService.ts"
import { CommonLevel } from "./LangService.ts"

export interface Term {
  id: number
  text: string
  variants: string[]
  partOfSpeech: string
  commonLevel: CommonLevel
  translations: Translation[]
  dictionarySource: string
}

export interface Translation {
  text: string
  explanation: string
}

export const TermsSearchDataEmpty = Object.freeze<TermsSearchData>({
  query: "",
  pos: [],
})

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
