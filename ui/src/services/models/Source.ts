export const PosPunctuation = "Punctuation"

export interface Token {
  text: string
  partOfSpeech: string
  startIndex: number
  length: number
}
const TokenEmpty = Object.freeze<Token>({
  text: "",
  partOfSpeech: "",
  startIndex: 0,
  length: 0,
})

export interface Text {
  text: string
  translation: string
  previousBreak: boolean
}
const TextEmpty = Object.freeze<Text>({
  text: "",
  translation: "",
  previousBreak: false,
})

export interface TokenizedText extends Text {
  tokens: Token[]
}
const TokenizedTextEmpty = Object.freeze<TokenizedText>({
  ...TextEmpty,
  tokens: [TokenEmpty],
})

export interface SourcePartMedia {
  imageUrl: string
  audioUrl: string
}
const SourcePartMediaEmpty = Object.freeze<SourcePartMedia>({
  imageUrl: "",
  audioUrl: "",
})

export interface SourcePart {
  media: SourcePartMedia
  tokenizedTexts: TokenizedText[]
}

const SourcePartEmpty = Object.freeze<SourcePart>({
  media: SourcePartMediaEmpty,
  tokenizedTexts: [TokenizedTextEmpty],
})

export interface Source {
  id: number
  name: string
  reference: string
  parts: SourcePart[]
  updatedAt: Date
  createdAt: Date
}
export const SourceEmpty = Object.freeze<Source>({
  id: 0,
  name: "",
  reference: "",
  parts: [SourcePartEmpty],
  updatedAt: new Date(0),
  createdAt: new Date(0),
})

export function partString(part: SourcePart): string {
  return part.tokenizedTexts
    .map((tokenizedText): string => {
      const lines = []
      if (tokenizedText.previousBreak) lines.push("")
      lines.push(tokenizedText.text)
      if (tokenizedText.translation !== "") lines.push(tokenizedText.translation)
      return lines.join("\n")
    })
    .join("\n")
}

export function tokenPreviousSpace(tokens: Token[], index: number): boolean {
  if (index === 0) return false
  const currentToken = tokens[index]
  const previousToken = tokens[index - 1]
  return previousToken.startIndex + previousToken.length + 1 === currentToken.startIndex
}

export function tokenPreviousPunct(tokens: Token[], index: number): boolean {
  if (index === 0) return false
  return tokens[index - 1].partOfSpeech === PosPunctuation
}
