import { CreateNoteData, CreateNoteDataKeys, notesService } from "../services/NotesService.ts"
import { formData } from "../utils/RequestUtil.ts"
import { ActionFunction } from "react-router-dom"

export const create: ActionFunction = async ({ request }) => {
  const data = formData(await request.formData(), ...CreateNoteDataKeys)
  data.commonLevel = parseInt(data.commonLevel as string, 10)
  return { note: notesService.create(data as CreateNoteData) }
}
