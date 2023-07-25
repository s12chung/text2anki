/* eslint-disable max-lines */
import { Source, Token } from "../../services/SourcesService.ts"
import { Term } from "../../services/TermsService.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { paginate, totalPages } from "../../utils/HtmlUtil.ts"
import { queryString } from "../../utils/UrlUtil.ts"
import AwaitError from "../AwaitError.tsx"
import SlideOver from "../SlideOver.tsx"
import NoteForm from "../note/NoteForm.tsx"
import React, { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { Await, Link, useFetcher } from "react-router-dom"

export interface ISourceShowData {
  source: Promise<Source>
}
interface ISourceShowProps {
  data: ISourceShowData
}

const SourceShow: React.FC<ISourceShowProps> = ({ data }) => {
  return (
    <React.Suspense fallback={<div>Loading....</div>}>
      <Await resolve={data.source} errorElement={<AwaitError />}>
        {(source: Source) => <SourceComponent source={source} />}
      </Await>
    </React.Suspense>
  )
}

function tokenPreviousSpace(tokens: Token[], index: number): boolean {
  if (index === 0) {
    return false
  }
  const currentToken = tokens[index]
  const previousToken = tokens[index - 1]
  return previousToken.startIndex + previousToken.length + 1 === currentToken.startIndex
}

function increment(index: number, length: number): number {
  if (index === -1) return 0
  return index < length - 1 ? index + 1 : 0
}

function decrement(index: number, length: number): number {
  if (index === -1) return 0
  return index > 0 ? index - 1 : length - 1
}

// eslint-disable-next-line max-lines-per-function
const SourceComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [textFocusIndex, setTextFocusIndex] = useState<number>(-1)
  const [tokenFocusIndex, setTokenFocusIndex] = useState<number>(-1)
  const [selectedToken, setSelectedToken] = useState<Token | null>(null)

  const textRefs = useRef<(HTMLDivElement | null)[]>([])
  const tokenRefs = useRef<(HTMLDivElement | null)[][]>([])

  // eslint-disable-next-line prefer-destructuring
  const tokenizedTexts = source.parts[0].tokenizedTexts

  useEffect(() => {
    setSelectedToken(null)
    const textElement = textRefs.current[textFocusIndex]
    if (!textElement) return
    textElement.focus()
    window.scrollTo({
      top: textElement.getBoundingClientRect().top + window.scrollY - 150,
      behavior: "smooth",
    })
    if (tokenFocusIndex === -1) return
    tokenRefs.current[textFocusIndex][tokenFocusIndex]?.focus()
  }, [textFocusIndex, tokenFocusIndex])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      const textLen = tokenizedTexts.length
      const tokenLen = tokenizedTexts[textFocusIndex]?.tokens.length

      switch (event.code) {
        case "Escape":
          if (selectedToken === null) return
          setSelectedToken(null)
          tokenRefs.current[textFocusIndex][tokenFocusIndex]?.focus()
          break
        default:
      }

      if (selectedToken !== null) return

      switch (event.code) {
        case "ArrowUp":
        case "KeyW":
          setTextFocusIndex(decrement(textFocusIndex, textLen))
          setTokenFocusIndex(-1)
          break
        case "ArrowDown":
        case "KeyS":
          setTextFocusIndex(increment(textFocusIndex, textLen))
          setTokenFocusIndex(-1)
          break
        case "ArrowLeft":
        case "KeyA":
          setTokenFocusIndex(decrement(tokenFocusIndex, tokenLen))
          if (textFocusIndex === -1) setTextFocusIndex(0)
          break
        case "ArrowRight":
        case "KeyD":
          setTokenFocusIndex(increment(tokenFocusIndex, tokenLen))
          if (textFocusIndex === -1) setTextFocusIndex(0)
          break
        case "Enter":
        case "Space":
          if (tokenFocusIndex === -1) return
          setSelectedToken(tokenizedTexts[textFocusIndex].tokens[tokenFocusIndex])
          break
        default:
          return
      }

      event.preventDefault()
    },
    [tokenizedTexts, textFocusIndex, tokenFocusIndex, selectedToken]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const handleTextClick = (index: number) => setTextFocusIndex(index)
  const handleTokenClick = (index: number) => setTokenFocusIndex(index)

  const termsFocus = selectedToken !== null
  const tokenizedTextClass = (b: boolean) =>
    `group py-2 focin:py-4 focin:bg-gray-std ${b ? "py-4 bg-gray-std" : ""}`

  const textClass = (b: boolean) => `ko-sans text-2xl focgrin:text-light ${b ? "text-light" : ""}`
  const translationClass = (b: boolean) => `text-lg focgrin:text-2xl ${b ? "text-2xl" : "text-lg"}`
  const tokenClass = (b: boolean) => `focus:text-white focus:bg-ink ${b ? "text-white bg-ink" : ""}`

  return (
    <>
      <div className="grid-std flex-std my-std">
        <div className="flex-grow">
          <h2>{source.name}</h2>
        </div>
        <div className="flex">
          <Link to={`/sources/${source.id}/edit`} className="btn">
            Edit
          </Link>
        </div>
      </div>

      <div className="text-center">
        {tokenizedTexts.map((tokenizedText, textIndex) => {
          const textFocus = textIndex === textFocusIndex
          return (
            /* eslint-disable-next-line react/no-array-index-key */
            <div key={`${tokenizedText.text}-${textIndex}`}>
              {Boolean(tokenizedText.previousBreak) && <div className="text-4xl">&nbsp;</div>}
              <div
                ref={(ref) => (textRefs.current[textIndex] = ref)}
                tabIndex={-1}
                className={tokenizedTextClass(textFocus)}
                onClick={() => handleTextClick(textIndex)}
              >
                <div className={textClass(textFocus)}>{tokenizedText.text}</div>
                {textIndex === textFocusIndex && (
                  <div className="ko-sans text-4xl justify-center mb-2 child:py-2 flex">
                    {tokenizedText.tokens.map((token, index) => {
                      const previousSpace = tokenPreviousSpace(tokenizedText.tokens, index)
                      return (
                        /* eslint-disable-next-line react/no-array-index-key */
                        <React.Fragment key={`${token.text}-${token.partOfSpeech}-${index}`}>
                          {!previousSpace && index !== 0 && <div>&middot;</div>}
                          {Boolean(previousSpace) && index !== 0 && <div>&nbsp;&nbsp;</div>}
                          <div
                            ref={(ref) => {
                              if (!tokenRefs.current[textIndex]) tokenRefs.current[textIndex] = []
                              tokenRefs.current[textIndex][index] = ref
                            }}
                            className={tokenClass(index === tokenFocusIndex && textFocus)}
                            tabIndex={-1}
                            onClick={() => handleTokenClick(index)}
                          >
                            <div>{token.text}</div>
                          </div>
                        </React.Fragment>
                      )
                    })}
                  </div>
                )}
                <div className={translationClass(textFocus)}>{tokenizedText.translation}</div>
                {textFocus && termsFocus ? <TermsComponent token={selectedToken} /> : null}
              </div>
            </div>
          )
        })}
      </div>
    </>
  )
}

interface ITermsShowData {
  terms: Term[]
}

const pageSize = 5

// eslint-disable-next-line max-lines-per-function
const TermsComponent: React.FC<{ token: Token }> = ({ token }) => {
  const fetcher = useFetcher<ITermsShowData>()
  const terms = useMemo<Term[]>(() => (fetcher.data ? fetcher.data.terms : []), [fetcher.data])

  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const termRefs = useRef<(HTMLDivElement | null)[]>([])

  const [page, setPage] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, pageSize), [terms])

  const [showCreateNote, setShowCreateNote] = useState<boolean>(false)
  const onCloseCreateNote = () => setShowCreateNote(false)

  useEffect(() => {
    if (fetcher.state !== "idle" || fetcher.data) return
    fetcher.load(`/terms/search?${queryString({ query: token.text, pos: token.partOfSpeech })}`)
  }, [fetcher, token])

  useEffect(() => {
    const termElement = termRefs.current[termFocusIndex]
    if (!termElement) return
    termElement.focus()
  }, [terms, page, termFocusIndex]) // trigger from terms/page to do initial focus

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      switch (event.code) {
        case "ArrowUp":
        case "KeyW":
          setTermFocusIndex(decrement(termFocusIndex, pageSize))
          break
        case "ArrowDown":
        case "KeyS":
          setTermFocusIndex(increment(termFocusIndex, pageSize))
          break
        case "ArrowLeft":
        case "KeyA":
          setPage(decrement(page, pagesLen))
          setTermFocusIndex(0)
          break
        case "ArrowRight":
        case "KeyD":
          setPage(increment(page, pagesLen))
          setTermFocusIndex(0)
          break
        case "Enter":
        case "Space":
          setShowCreateNote(true)
          break
        default:
          return
      }
      event.preventDefault()
    },
    [termFocusIndex, page, pagesLen]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const topLevelClass = "grid-std text-left text-lg py-2 space-y-2"
  if (!fetcher.data) {
    return <div className={topLevelClass}>Loading...</div>
  }

  return (
    <div className={topLevelClass}>
      {terms.length === 0 ? (
        <div>No terms found</div>
      ) : (
        <div>
          {paginate(terms, pageSize, page).map((term, index) => (
            <div
              key={term.id}
              ref={(ref) => (termRefs.current[index] = ref)}
              tabIndex={-1}
              className="focus:underline"
            >
              <div className="text-xl">
                {term.text}&nbsp;
                <span className="text-light text-base">{term.partOfSpeech}</span>
                {term.commonLevel !== 0 && (
                  <span className="relative top-2">&nbsp;{"*".repeat(term.commonLevel)}</span>
                )}
                : {term.translations[0].text} &mdash; {term.translations[0].explanation}
              </div>
              <div className="ml-8">
                {unique(term.translations.map((translation) => translation.text))
                  .slice(1, 6)
                  .join("; ")}
              </div>
            </div>
          ))}
          <div className="text-center space-x-2">
            {new Array(pagesLen).fill(null).map((_, index) => (
              /* eslint-disable-next-line react/no-array-index-key */
              <span key={index} className={index === page ? "" : "text-light"}>
                {index === page ? <>&#x2716;</> : <>&bull;</>}
              </span>
            ))}
          </div>
        </div>
      )}

      <SlideOver.Dialog show={showCreateNote} onClose={onCloseCreateNote}>
        <SlideOver.Header title="Create Note" onClose={onCloseCreateNote} />
        <NoteForm onClose={onCloseCreateNote} />
      </SlideOver.Dialog>
    </div>
  )
}

export default SourceShow
