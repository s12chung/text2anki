import { Term, Translation } from "../../services/models/Term.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { pageSize, totalPages } from "../../utils/HtmlUtil.ts"
import { useKeyDownEffect, useTimedState } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useContext, useMemo, useState } from "react"

const maxPageSize = 5

export function useChangeTermWithKeyboard(
  terms: Term[],
  onTermSelect: (term: Term) => void,
): readonly [number, number, number, number, boolean] {
  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const [pageIndex, setPageIndex] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, maxPageSize), [terms])
  const [shake, setShake] = useTimedState(100)

  const { stopKeyboardEvents } = useContext(StopKeyboardContext)

  useKeyDownEffect(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return
      if (terms.length === 1) {
        setShake(true)
        e.preventDefault()
        return
      }

      switch (e.code) {
        case "ArrowUp":
        case "KeyW":
          setTermFocusIndex(
            decrement(termFocusIndex, pageSize(terms.length, maxPageSize, pageIndex)),
          )
          break
        case "ArrowDown":
        case "KeyS":
          setTermFocusIndex(
            increment(termFocusIndex, pageSize(terms.length, maxPageSize, pageIndex)),
          )
          break
        case "ArrowLeft":
        case "KeyA":
          setPageIndex(decrement(pageIndex, pagesLen))
          setTermFocusIndex(0)
          break
        case "ArrowRight":
        case "KeyD":
          setPageIndex(increment(pageIndex, pagesLen))
          setTermFocusIndex(0)
          break
        case "Enter":
        case "Space":
          onTermSelect(terms[termFocusIndex])
          break
        default:
          return
      }
      e.preventDefault()
    },
    [stopKeyboardEvents, termFocusIndex, terms, pageIndex, pagesLen, onTermSelect, setShake],
  )
  return [termFocusIndex, pageIndex, pagesLen, maxPageSize, shake] as const
}

export function otherTranslationTexts(translations: Translation[]): string {
  return unique(translations.map((translation) => translation.text))
    .slice(1, 6)
    .join("; ")
}
