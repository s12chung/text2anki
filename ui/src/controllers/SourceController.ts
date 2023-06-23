import { CreateSourceData, sourceService, UpdateSourceData } from "../services/SourceService.ts"
import { formData } from "../utils/RouterUtils.ts"
import { ActionFunction, defer, LoaderFunction, redirect } from "react-router-dom"

export const index: LoaderFunction = () => {
  return defer({ sources: sourceService.index() })
}

export const get: LoaderFunction = ({ params }) => {
  return defer({ source: sourceService.get(params.id as string) })
}

export const create: ActionFunction = async ({ request }) => {
  const source = await sourceService.create(
    formData<CreateSourceData>(await request.formData(), "text")
  )
  return redirect(`/sources/${source.id}`)
}

export const update: ActionFunction = async ({ request, params }) => {
  const source = await sourceService.update(
    params.id as string,
    formData<UpdateSourceData>(await request.formData(), "name")
  )
  return redirect(`/sources/${source.id}`)
}
