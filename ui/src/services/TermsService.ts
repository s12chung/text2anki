import { queryString } from "../utils/RequestUtil.ts"
import ApplicationService from "./ApplicationService.ts"
import { CommonLevel } from "./Lang.ts"

export interface Translation {
  text: string
  explanation: string
}

const TranslationEmpty = Object.freeze<Translation>({
  text: "",
  explanation: "",
})

export interface Term {
  id: number
  text: string
  variants: string[]
  partOfSpeech: string
  commonLevel: CommonLevel
  translations: Translation[]
  dictionarySource: string
}

const TermEmpty = Object.freeze<Term>({
  id: 0,
  text: "",
  variants: [""],
  partOfSpeech: "",
  commonLevel: 0,
  translations: [TranslationEmpty],
  dictionarySource: "",
})

export interface TermsSearchData {
  query: string
  pos: string | string[]
}

export const TermsSearchDataEmpty = Object.freeze<TermsSearchData>({
  query: "",
  pos: [],
})

class TermsService extends ApplicationService {
  protected pathPrefix = "/terms"

  async search(data: TermsSearchData): Promise<Term[]> {
    return this.fetch(`/search?${queryString(data)}`, [TermEmpty])
  }
}

export const termsService = new TermsService()
