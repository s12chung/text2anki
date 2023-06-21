import { CreateSourceData, sourceService } from "../services/SourceService.ts"
import { formData } from "../utils/RouterUtils.ts"
import { ActionFunction, ActionFunctionArgs, defer, redirect } from "react-router-dom"

export function index() {
  return defer({ sources: sourceService.index() })
}

export const create: ActionFunction = async ({ request }) => {
  const source = await sourceService.create(
    formData<CreateSourceData>(await request.formData(), "text")
  )
  return redirect(`/sources/${source.id}`)
}

export function get({ params }: ActionFunctionArgs) {
  return defer({ source: sourceService.get(params.id as string) })
}
