import { Term, Translation } from "../../services/TermsService.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { pageSize, totalPages } from "../../utils/HtmlUtil.ts"
import { useKeyDownEffect, useTimedState } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useContext, useEffect, useMemo, useState } from "react"

const maxPageSize = 5

export function useChangeTermWithKeyboard(
  terms: Term[],
  onEnter: (term: Term) => void,
  isEntered: () => boolean
): readonly [number, number, number, number, boolean] {
  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const [pageIndex, setPageIndex] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, maxPageSize), [terms])
  const [shake, setShake] = useTimedState(100)

  const { stopKeyboardEvents, setStopKeyboardEvents } = useContext(StopKeyboardContext)

  useEffect(() => setStopKeyboardEvents(isEntered()), [isEntered, setStopKeyboardEvents])
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
            decrement(termFocusIndex, pageSize(terms.length, maxPageSize, pageIndex))
          )
          break
        case "ArrowDown":
        case "KeyS":
          setTermFocusIndex(
            increment(termFocusIndex, pageSize(terms.length, maxPageSize, pageIndex))
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
          onEnter(terms[termFocusIndex])
          break
        default:
          return
      }
      e.preventDefault()
    },
    [stopKeyboardEvents, termFocusIndex, terms, pageIndex, pagesLen, onEnter, setShake]
  )
  return [termFocusIndex, pageIndex, pagesLen, maxPageSize, shake] as const
}

export function otherTranslationTexts(translations: Translation[]): string {
  return unique(translations.map((translation) => translation.text))
    .slice(1, 6)
    .join("; ")
}
