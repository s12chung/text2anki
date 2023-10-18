import { PartCreateMultiData } from "../services/PartsService.ts"
import {
  CreateSourceDataEmpty,
  sourcesService,
  UpdateSourceDataEmpty,
} from "../services/SourcesService.ts"
import { formData } from "../utils/RequestUtil.ts"
import { Status405 } from "../utils/StatusUtil.ts"
import { createPrePart } from "./PartsController.ts"
import { ActionFunction, defer, LoaderFunction, redirect } from "react-router-dom"

export const index: LoaderFunction = () => {
  return defer({ sources: sourcesService.index() })
}

export const get: LoaderFunction = ({ params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  return defer({ source: sourcesService.get(params.id) })
}

export const create: ActionFunction = async ({ request }) => {
  const data = formData(await request.formData(), CreateSourceDataEmpty)
  const resp = await checkAndCreatePrePart(data)
  if (resp) return resp
  return redirect(`/sources/${(await sourcesService.create(data)).id}`)
}

export const update: ActionFunction = async ({ request, params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  const source = await sourcesService.update(
    params.id,
    formData(await request.formData(), UpdateSourceDataEmpty),
  )
  return { source }
}

export const destroy: ActionFunction = async ({ params }) => {
  if (!params.id) throw new Response("id not found", Status405) // eslint-disable-line @typescript-eslint/no-throw-literal
  await sourcesService.destroy(params.id)
  // eslint-disable-next-line no-warning-comments
  // TODO: pass notification via. localstorage
  return redirect(`/`)
}

async function checkAndCreatePrePart(data: PartCreateMultiData): Promise<Response | null> {
  if (data.prePartListId || data.parts.length !== 1 || data.parts[0].translation) {
    return null
  }
  return createPrePart(data.parts[0].text)
}
