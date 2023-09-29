import { CommonLevel } from "./Lang.ts"

export interface Translation {
  text: string
  explanation: string
}
const TranslationEmpty = Object.freeze<Translation>({
  text: "",
  explanation: "",
})

export interface Term {
  id: number
  text: string
  variants: string[]
  partOfSpeech: string
  commonLevel: CommonLevel
  translations: Translation[]
  dictionarySource: string
}
export const TermEmpty = Object.freeze<Term>({
  id: 0,
  text: "",
  variants: [""],
  partOfSpeech: "",
  commonLevel: 0,
  translations: [TranslationEmpty],
  dictionarySource: "",
})
