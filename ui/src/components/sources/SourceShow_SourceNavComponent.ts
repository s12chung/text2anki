import { Source, SourcePart } from "../../services/models/Source.ts"
import { useKeyDownEffect } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useCallback, useContext, useState } from "react"

// eslint-disable-next-line max-params
export function useFocusTextWithKeyboard(
  parts: SourcePart[],
  isTokenSelected: boolean,
  onCreateTextNote: () => void,
  onEscape: () => void,
): readonly [number, number, (partFocusIndex: number, textFocusIndex: number) => void] {
  const [partFocusIndex, textFocusIndex, decrementText, incrementText, setText] =
    useChangeFocus(parts)

  const { stopKeyboardEvents } = useContext(StopKeyboardContext)
  useKeyDownEffect(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return

      switch (e.code) {
        case "Escape":
          if (!isTokenSelected) return
          onEscape()
          e.preventDefault()
          return
        default:
      }

      if (isTokenSelected) return

      switch (e.code) {
        case "ArrowUp":
        case "KeyW":
          decrementText()
          break
        case "ArrowDown":
        case "KeyS":
          incrementText()
          break
        case "KeyN":
          onCreateTextNote()
          break
        default:
          return
      }

      e.preventDefault()
    },
    [stopKeyboardEvents, isTokenSelected, onEscape, decrementText, incrementText, onCreateTextNote],
  )
  return [partFocusIndex, textFocusIndex, setText] as const
}

export function getUsage(source: Source, partFocusIndex: number, textFocusIndex: number) {
  const tokenizedText = source.parts[partFocusIndex].tokenizedTexts[textFocusIndex]
  return {
    sourceName: source.name,
    sourceReference: source.reference,

    usage: tokenizedText.text,
    usageTranslation: tokenizedText.translation,
  }
}

function useChangeFocus(
  parts: SourcePart[],
): readonly [
  number,
  number,
  () => void,
  () => void,
  (partFocusIndex: number, textFocusIndex: number) => void,
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
