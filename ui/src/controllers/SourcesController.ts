import {
  CreateSourcePartData,
  sourcesService,
  UpdateSourceData,
} from "../services/SourcesService.ts"
import { formData } from "../utils/RouterUtil.ts"
import { ActionFunction, defer, LoaderFunction, redirect } from "react-router-dom"

export const index: LoaderFunction = () => {
  return defer({ sources: sourcesService.index() })
}

export const get: LoaderFunction = ({ params }) => {
  return defer({ source: sourcesService.get(params.id as string) })
}

export const create: ActionFunction = async ({ request }) => {
  const source = await sourcesService.create({
    parts: [formData<CreateSourcePartData>(await request.formData(), "text", "translation")],
  })
  return redirect(`/sources/${source.id}`)
}

export const update: ActionFunction = async ({ request, params }) => {
  const source = await sourcesService.update(
    params.id as string,
    formData<UpdateSourceData>(await request.formData(), "name")
  )
  return redirect(`/sources/${source.id}`)
}
