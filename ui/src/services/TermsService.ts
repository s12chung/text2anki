import { queryString } from "../utils/RequestUtil.ts"
import ApplicationService from "./ApplicationService.ts"
import { Term, TermEmpty } from "./models/Term.ts"

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
