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

  const source = await partsService.create(sourceId, data)
  return { source }
}

export const multi: ActionFunction = async ({ request, params }) => {
  const sourceId = getSourceId(params)
  await partsService.multi(sourceId, formData(await request.formData(), PartCreateMultiDataEmpty))
  return redirect(`/sources/${sourceId}`)
}

export const update: ActionFunction = async ({ request, params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  if (!params.sourceId) throw new Response("sourceId not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  const source = await partsService.update(
    params.sourceId,
    params.id,
    formData(await request.formData(), PartDataEmpty)
  )
  return { source }
}

export const destroy: ActionFunction = async ({ params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  if (!params.sourceId) throw new Response("sourceId not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  const source = await partsService.destroy(params.sourceId, params.id)
  return { source }
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
