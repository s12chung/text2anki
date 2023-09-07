import ApplicationService, { Http, requestInit } from "./ApplicationService.ts"
import { CommonLevel } from "./Lang.ts"
import { Term, Translation } from "./TermsService.ts"

export interface Note extends CreateNoteData {
  id: number
  downloaded: boolean
}

export interface NoteUsage {
  usage: string
  usageTranslation: string

  sourceName: string
  sourceReference: string
}

export const NoteUsageEmpty = Object.freeze<NoteUsage>({
  usage: "",
  usageTranslation: "",

  sourceName: "",
  sourceReference: "",
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

export const CreateNoteDataEmpty = Object.freeze<CreateNoteData>({
  text: "",
  partOfSpeech: "",
  translation: "",
  explanation: "",
  commonLevel: CommonLevel.Unique,

  ...NoteUsageEmpty,
  dictionarySource: "",
  notes: "",
})

// eslint-disable-next-line max-params
export function createNoteDataFromSourceTerm(
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
    explanation: translation.explanation,
    commonLevel: term.commonLevel,

    usage: usage.usage,
    usageTranslation: usage.usageTranslation,

    sourceName: usage.sourceName,
    sourceReference: usage.sourceReference,
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
