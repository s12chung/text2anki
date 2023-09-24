import {
  PartCreateMultiDataEmpty,
  PartData,
  PartDataEmpty,
  partsService,
} from "../services/PartsService.ts"
import { prePartListService } from "../services/PrePartListsService.ts"
import { formData, queryString } from "../utils/RequestUtil.ts"
import { Status405 } from "../utils/StatusUtil.ts"
import { sourceIdQueryParam } from "./PrePartListsController.ts"
import { ActionFunction, Params, redirect } from "react-router-dom"

export const create: ActionFunction = async ({ request, params }) => {
  const sourceId = getSourceId(params)
  const data = formData(await request.formData(), PartDataEmpty)
  const resp = await checkAndCreatePrePart(data, sourceId)
  if (resp) return resp
  await partsService.create(sourceId, data)
  // eslint-disable-next-line no-warning-comments
  // TODO: handle proper loading of sources
  window.location.reload()
  return redirect(`/sources/${sourceId}`)
}

export const multi: ActionFunction = async ({ request, params }) => {
  const sourceId = getSourceId(params)
  await partsService.multi(sourceId, formData(await request.formData(), PartCreateMultiDataEmpty))
  return redirect(`/sources/${sourceId}`)
}

function getSourceId(params: Params): number {
  const sourceId = Number(params.sourceId)
  if (!sourceId) throw new Response("sourceId not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  return sourceId
}

async function checkAndCreatePrePart(data: PartData, sourceId: number): Promise<Response | null> {
  if (data.translation) return null
  return createPrePart(data.text, sourceId)
}

export async function createPrePart(text: string, sourceId?: number): Promise<Response | null> {
  text = text.trim()
  if (text.includes("\n") || text.includes("\r")) {
    return null
  }
  const { extractorType } = await prePartListService.verify({ text })
  if (extractorType === "") {
    return null
  }

  const prePartListId = (await prePartListService.create({ extractorType, text })).id
  const query = sourceId ? `?${queryString({ [sourceIdQueryParam]: String(sourceId) })}` : ""
  return redirect(`/sources/pre_part_lists/${prePartListId}${query}`)
}
