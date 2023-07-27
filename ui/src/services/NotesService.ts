import ApplicationService, { Http, requestInit } from "./ApplicationService.ts"
import { CommonLevel } from "./LangService.ts"
import { Term, Translation } from "./TermsService.ts"

export interface Note extends CreateNoteData {
  id: number
  downloaded: boolean
}

export const NotUsageEmpty = Object.freeze<NoteUsage>({
  usage: "",
  usageTranslation: "",
})

export const CreateNoteDataEmpty = Object.freeze<CreateNoteData>({
  text: "",
  partOfSpeech: "",
  translation: "",
  commonLevel: CommonLevel.Unique,
  explanation: "",
  ...NotUsageEmpty,
  dictionarySource: "",
  notes: "",
})

export interface CreateNoteData extends NoteUsage {
  text: string
  partOfSpeech: string
  translation: string
  commonLevel: CommonLevel
  explanation: string
  dictionarySource: string
  notes: string
}

export interface NoteUsage {
  usage: string
  usageTranslation: string
}

export function createNoteDataFromTerm(
  term: Term,
  usage: NoteUsage,
  translationIndex?: number
): CreateNoteData {
  if (!translationIndex) translationIndex = 0
  const translation: Translation = term.translations[translationIndex]
  return {
    text: term.text,
    partOfSpeech: term.partOfSpeech,
    translation: translation.text,
    commonLevel: term.commonLevel,
    explanation: translation.explanation,
    usage: usage.usage,
    usageTranslation: usage.usageTranslation,
    dictionarySource: term.dictionarySource,
    notes: "",
  }
}

class NotesService extends ApplicationService {
  protected pathPrefix = "/notes"

  async create(data: CreateNoteData): Promise<Note> {
    return (await this.fetch("", requestInit(Http.POST, data))) as Note
  }
}

export const notesService = new NotesService()
