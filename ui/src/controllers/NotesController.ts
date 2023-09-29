import { notesService } from "../services/NotesService.ts"
import { CreateNoteDataEmpty } from "../services/models/Note.ts"
import { formData } from "../utils/RequestUtil.ts"
import { ActionFunction, defer, LoaderFunction } from "react-router-dom"

export const create: ActionFunction = async ({ request }) => {
  return { note: notesService.create(formData(await request.formData(), CreateNoteDataEmpty)) }
}

export const index: LoaderFunction = () => {
  return defer({ sources: notesService.index() })
}
