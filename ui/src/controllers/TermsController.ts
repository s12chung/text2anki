import { TermsSearchData, TermsSearchDataKeys, termsService } from "../services/TermsService.ts"
import { queryObject } from "../utils/RequestUtil.ts"
import { defer, LoaderFunction } from "react-router-dom"

export const search: LoaderFunction = ({ request }) => {
  return defer({
    terms: termsService.search(queryObject<TermsSearchData>(request.url, ...TermsSearchDataKeys)),
  })
}
