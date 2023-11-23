import { PosPunctuation, Token } from "../../services/models/Source.ts"
import { useKeyDownEffect } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useContext, useMemo, useState } from "react"

export interface SelectedToken {
  text: string
  partOfSpeech: string
}

export function useFocusTokenWithKeyboard(
  tokens: Token[],
  isTokenSelected: boolean,
  onTokenSelect: (token: SelectedToken) => void,
): readonly [number] {
  const [tokenFocusIndex, setTokenFocusIndex] = useState<number>(0)
  const isAllPunct = useMemo<boolean>(
    () => tokens.every((token) => token.partOfSpeech === PosPunctuation),
    [tokens],
  )

  const { stopKeyboardEvents } = useContext(StopKeyboardContext)
  useKeyDownEffect(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents || isTokenSelected || isAllPunct) return

      switch (e.code) {
        case "ArrowLeft":
        case "KeyA":
          setTokenFocusIndex(skipPunct(tokens, tokenFocusIndex, decrement))
          break
        case "ArrowRight":
        case "KeyD":
          setTokenFocusIndex(skipPunct(tokens, tokenFocusIndex, increment))
          break
        case "Enter":
        case "Space":
          onTokenSelect(tokens[tokenFocusIndex])
          break
        default:
          return
      }

      e.preventDefault()
    },
    [stopKeyboardEvents, isTokenSelected, isAllPunct, tokens, tokenFocusIndex, onTokenSelect],
  )
  return [tokenFocusIndex] as const
}

function skipPunct(
  tokens: Token[],
  index: number,
  change: (index: number, length: number) => number,
): number {
  index = change(index, tokens.length)
  if (tokens[index].partOfSpeech !== PosPunctuation) return index
  return skipPunct(tokens, index, change)
}
