import { NoteUsage } from "../../services/NotesService.ts"
import { Source, SourcePart, Token } from "../../services/SourcesService.ts"
import { useKeyDownEffect } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useCallback, useContext, useState } from "react"

export function useFocusTextWithKeyboard(
  parts: SourcePart[],
  entered: boolean,
  onEscape: () => void
): readonly [
  number,
  number,
  (element: HTMLElement) => void,
  (partFocusIndex: number, textFocusIndex: number) => void
] {
  const [focusElement, focusLastElement] = useFocusElement()
  const [partFocusIndex, textFocusIndex, decrementText, incrementText, setText] =
    useChangeFocus(parts)

  const { stopKeyboardEvents } = useContext(StopKeyboardContext)
  useKeyDownEffect(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return

      switch (e.code) {
        case "Escape":
          if (!entered) return
          focusLastElement()
          onEscape()
          break
        default:
      }

      if (entered) return

      switch (e.code) {
        case "ArrowUp":
        case "KeyW":
          decrementText()
          break
        case "ArrowDown":
        case "KeyS":
          incrementText()
          break

        default:
          return
      }

      e.preventDefault()
    },
    [stopKeyboardEvents, entered, focusLastElement, onEscape, decrementText, incrementText]
  )
  return [partFocusIndex, textFocusIndex, focusElement, setText] as const
}

export interface ITermsComponentProps {
  token: Token
  usage: NoteUsage
}

// eslint-disable-next-line max-params
export function getTermProps(
  source: Source,
  partFocusIndex: number,
  textFocusIndex: number,
  tokenFocusIndex: number
): ITermsComponentProps {
  const tokenizedText = source.parts[partFocusIndex].tokenizedTexts[textFocusIndex]
  return {
    token: tokenizedText.tokens[tokenFocusIndex],
    usage: {
      sourceName: source.name,
      sourceReference: source.reference,

      usage: tokenizedText.text,
      usageTranslation: tokenizedText.translation,
    },
  }
}

function useFocusElement(): readonly [(element: HTMLElement) => void, () => void] {
  const [lastFocusedElement, setLastFocusedElement] = useState<HTMLElement | null>(null)
  const focusElement = useCallback((element: HTMLElement) => {
    setLastFocusedElement(element)
    element.focus()
  }, [])

  return [focusElement, () => lastFocusedElement?.focus()] as const
}

function useChangeFocus(
  parts: SourcePart[]
): readonly [
  number,
  number,
  () => void,
  () => void,
  (partFocusIndex: number, textFocusIndex: number) => void
] {
  const [partFocusIndex, setPartFocusIndex] = useState<number>(0)
  const [textFocusIndex, setTextFocusIndex] = useState<number>(0)

  const decrementText = useCallback(() => {
    const currentTokenizedTextsLen = parts[partFocusIndex].tokenizedTexts.length
    const result = decrement(textFocusIndex, currentTokenizedTextsLen)
    if (result !== currentTokenizedTextsLen - 1) {
      setTextFocusIndex(result)
      return
    }

    const partIndex = decrement(partFocusIndex, parts.length)
    setPartFocusIndex(partIndex)
    setTextFocusIndex(parts[partIndex].tokenizedTexts.length - 1)
  }, [parts, partFocusIndex, textFocusIndex])

  const incrementText = useCallback(() => {
    const currentTokenizedTextsLen = parts[partFocusIndex].tokenizedTexts.length
    const result = increment(textFocusIndex, currentTokenizedTextsLen)
    if (result !== 0) {
      setTextFocusIndex(result)
      return
    }

    const partIndex = increment(partFocusIndex, parts.length)
    setPartFocusIndex(partIndex)
    setTextFocusIndex(0)
  }, [parts, partFocusIndex, textFocusIndex, setPartFocusIndex, setTextFocusIndex])

  const setText = useCallback((partFocusIndex: number, textFocusIndex: number) => {
    setPartFocusIndex(partFocusIndex)
    setTextFocusIndex(textFocusIndex)
  }, [])

  return [partFocusIndex, textFocusIndex, decrementText, incrementText, setText] as const
}
