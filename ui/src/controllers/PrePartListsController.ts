import { prePartListService } from "../services/PrePartListsService.ts"
import { Status405 } from "../utils/StatusUtil.ts"
import { defer, LoaderFunction } from "react-router-dom"

export const sourceIdQueryParam = "source_id"

export const get: LoaderFunction = ({ request, params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  const sourceId = Number(new URL(request.url).searchParams.get(sourceIdQueryParam))
  return defer({ prePartList: prePartListService.get(params.id), sourceId })
}
