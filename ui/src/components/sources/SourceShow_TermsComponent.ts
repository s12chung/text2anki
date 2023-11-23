import { Term, Translation } from "../../services/models/Term.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { pageSize, totalPages } from "../../utils/HtmlUtil.ts"
import { useKeyDownEffect, useTimedState } from "../../utils/JSXUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { StopKeyboardContext } from "./SourceShow_SourceComponent.ts"
import { useCallback, useContext, useMemo, useState } from "react"

const maxPageSize = 5

export function useChangeTermWithKeyboard(
  terms: Term[],
  onTermSelect: (term: Term) => void,
): readonly [number, number, number, number, boolean] {
  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const [pageIndex, setPageIndex] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, maxPageSize), [terms])
  const [shake, setShake] = useTimedState(100)

  const changeTermFocusIndex = useCallback(
    (e: KeyboardEvent, change: (index: number, length: number) => number) => {
      if (terms.length === 1) {
        setShake(true)
        e.preventDefault()
        return
      }
      setTermFocusIndex(change(termFocusIndex, pageSize(terms.length, maxPageSize, pageIndex)))
    },
    [pageIndex, setShake, termFocusIndex, terms.length],
  )
  const changePage = useCallback(
    (e: KeyboardEvent, change: (index: number, length: number) => number) => {
      if (pagesLen === 1) {
        setShake(true)
        e.preventDefault()
        return
      }
      setPageIndex(change(pageIndex, pagesLen))
      setTermFocusIndex(0)
    },
    [pageIndex, pagesLen, setShake],
  )

  const { stopKeyboardEvents } = useContext(StopKeyboardContext)

  useKeyDownEffect(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return

      switch (e.code) {
        case "ArrowUp":
        case "KeyW":
          changeTermFocusIndex(e, decrement)
          break
        case "ArrowDown":
        case "KeyS":
          changeTermFocusIndex(e, increment)
          break
        case "ArrowLeft":
        case "KeyA":
          changePage(e, decrement)
          break
        case "ArrowRight":
        case "KeyD":
          changePage(e, increment)
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
    [changePage, changeTermFocusIndex, onTermSelect, stopKeyboardEvents, termFocusIndex, terms],
  )
  return [termFocusIndex, pageIndex, pagesLen, maxPageSize, shake] as const
}

export function otherTranslationTexts(translations: Translation[]): string {
  return unique(translations.map((translation) => translation.text))
    .slice(1, 6)
    .join("; ")
}
