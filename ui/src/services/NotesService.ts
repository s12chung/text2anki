import ApplicationService from "./ApplicationService.ts"
import { Http, requestInit } from "./Format.ts"
import { CreateNoteData, Note, NoteEmpty } from "./models/Note.ts"

class NotesService extends ApplicationService {
  protected pathPrefix = "/notes"

  downloadUrl(): string {
    return this.pathUrl("/download")
  }

  async index(): Promise<Note[]> {
    return this.fetch("", [NoteEmpty])
  }

  async create(data: CreateNoteData): Promise<Note> {
    return this.fetch("", NoteEmpty, requestInit(Http.POST, data))
  }
}

export const notesService = new NotesService()
