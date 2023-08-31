import { prePartListService } from "../services/PrePartListsService.ts"
import {
  CreateSourceData,
  CreateSourceDataEmpty,
  sourcesService,
  UpdateSourceDataEmpty,
} from "../services/SourcesService.ts"
import { formData } from "../utils/RequestUtil.ts"
import { Status405 } from "../utils/StatusUtil.ts"
import { ActionFunction, defer, LoaderFunction, redirect } from "react-router-dom"

export const index: LoaderFunction = () => {
  return defer({ sources: sourcesService.index() })
}

export const get: LoaderFunction = ({ params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  return defer({ source: sourcesService.get(params.id) })
}

export const edit = get

export const create: ActionFunction = async ({ request }) => {
  const data = formData(await request.formData(), CreateSourceDataEmpty)
  const resp = await createPrePart(data)
  if (resp) return resp
  return redirect(`/sources/${(await sourcesService.create(data)).id}`)
}

export const update: ActionFunction = async ({ request, params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  const source = await sourcesService.update(
    params.id,
    formData(await request.formData(), UpdateSourceDataEmpty)
  )
  return redirect(`/sources/${source.id}`)
}

export const destroy: ActionFunction = async ({ params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  await sourcesService.destroy(params.id)
  // eslint-disable-next-line no-warning-comments
  // TODO: pass notification via. localstorage
  return redirect(`/`)
}

async function createPrePart(data: CreateSourceData): Promise<Response | null> {
  if (data.prePartListId || data.parts.length !== 1 || data.parts[0].translation) {
    return null
  }
  // eslint-disable-next-line prefer-destructuring
  const text = data.parts[0].text.trim()
  if (text.includes("\n") || text.includes("\r")) {
    return null
  }
  const { extractorType } = await prePartListService.verify({ text })

  if (extractorType === "") {
    return null
  }
  return redirect(
    `/sources/pre_part_lists/${(await prePartListService.create({ extractorType, text })).id}`
  )
}
