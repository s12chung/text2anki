/* eslint-disable max-lines */
import { CommonLevel } from "../../services/LangService.ts"
import { CreateNoteData, createNoteDataFromTerm, NoteUsage } from "../../services/NotesService.ts"
import { Source, Token, TokenizedText } from "../../services/SourcesService.ts"
import { Term } from "../../services/TermsService.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { paginate, totalPages } from "../../utils/HtmlUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { queryString } from "../../utils/RequestUtil.ts"
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

function getTermsComponentProps(
  tokenizedText: TokenizedText,
  tokenFocusIndex: number
): ITermsComponentProps {
  return {
    token: tokenizedText.tokens[tokenFocusIndex],
    usage: {
      usage: tokenizedText.text,
      usageTranslation: tokenizedText.translation,
    },
  }
}

let openModal = false

// eslint-disable-next-line max-lines-per-function,max-statements
const SourceComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [partFocusIndex, setPartFocusIndex] = useState<number>(0)
  const [textFocusIndex, setTextFocusIndex] = useState<number>(-1)
  const [tokenFocusIndex, setTokenFocusIndex] = useState<number>(-1)
  const [termsComponentProps, setTermsComponentProps] = useState<ITermsComponentProps | null>(null)

  const textRefs = useRef<(HTMLDivElement | null)[][]>([])
  const tokenRefs = useRef<(HTMLDivElement | null)[][][]>([])

  // eslint-disable-next-line prefer-destructuring
  const currentTokenizedTexts = useMemo<TokenizedText[]>(
    () => source.parts[partFocusIndex].tokenizedTexts,
    [partFocusIndex, source.parts]
  )

  const partsLength = source.parts.length
  const decrementText = useCallback(() => {
    const result = decrement(textFocusIndex, currentTokenizedTexts.length)
    if (result === currentTokenizedTexts.length - 1) {
      setPartFocusIndex(decrement(partFocusIndex, partsLength))
    }
    setTextFocusIndex(result)
    setTokenFocusIndex(-1)
  }, [currentTokenizedTexts.length, partFocusIndex, partsLength, textFocusIndex])
  const incrementText = useCallback(() => {
    const result = increment(textFocusIndex, currentTokenizedTexts.length)
    if (result === 0) {
      setPartFocusIndex(increment(partFocusIndex, partsLength))
    }
    setTextFocusIndex(result)
    setTokenFocusIndex(-1)
  }, [currentTokenizedTexts.length, partFocusIndex, partsLength, textFocusIndex])

  useEffect(() => {
    setTermsComponentProps(null)
    const textElement = textRefs.current[partFocusIndex][textFocusIndex]
    if (!textElement) return
    textElement.focus()
    window.scrollTo({
      top: textElement.getBoundingClientRect().top + window.scrollY - 150,
      behavior: "smooth",
    })
    if (tokenFocusIndex === -1) return
    tokenRefs.current[partFocusIndex][textFocusIndex][tokenFocusIndex]?.focus()
  }, [partFocusIndex, textFocusIndex, tokenFocusIndex])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (openModal) return

      switch (e.code) {
        case "Escape":
          if (termsComponentProps === null) return
          setTermsComponentProps(null)
          tokenRefs.current[partFocusIndex][textFocusIndex][tokenFocusIndex]?.focus()
          break
        default:
      }

      if (termsComponentProps !== null) return

      const tokenLen = currentTokenizedTexts[textFocusIndex]?.tokens.length
      switch (e.code) {
        case "ArrowUp":
        case "KeyW":
          decrementText()
          break
        case "ArrowDown":
        case "KeyS":
          incrementText()
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
          setTermsComponentProps(
            getTermsComponentProps(currentTokenizedTexts[textFocusIndex], tokenFocusIndex)
          )
          break
        default:
          return
      }

      e.preventDefault()
    },
    [
      termsComponentProps,
      currentTokenizedTexts,
      partFocusIndex,
      textFocusIndex,
      tokenFocusIndex,
      decrementText,
      incrementText,
    ]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const textOnClick = (index: number) => setTextFocusIndex(index)
  const tokenOnClick = (index: number) => setTokenFocusIndex(index)

  const termsFocus = termsComponentProps !== null
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
        {source.parts.map((part, partIndex) =>
          part.tokenizedTexts.map((tokenizedText, textIndex) => {
            const textFocus = partIndex === partFocusIndex && textIndex === textFocusIndex
            return (
              /* eslint-disable-next-line react/no-array-index-key */
              <div key={`${tokenizedText.text}-${textIndex}`}>
                {Boolean(tokenizedText.previousBreak) && <div className="text-4xl">&nbsp;</div>}
                <div
                  ref={(ref) => {
                    if (!textRefs.current[partIndex]) textRefs.current[partIndex] = []
                    textRefs.current[partIndex][textIndex] = ref
                  }}
                  tabIndex={-1}
                  className={tokenizedTextClass(textFocus)}
                  onClick={() => textOnClick(textIndex)}
                >
                  <div className={textClass(textFocus)}>{tokenizedText.text}</div>
                  {textFocus ? (
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
                                if (!tokenRefs.current[partIndex]) tokenRefs.current[partIndex] = []
                                if (!tokenRefs.current[partIndex][textIndex])
                                  tokenRefs.current[partIndex][textIndex] = []
                                tokenRefs.current[partIndex][textIndex][index] = ref
                              }}
                              className={tokenClass(index === tokenFocusIndex && textFocus)}
                              tabIndex={-1}
                              onClick={() => tokenOnClick(index)}
                            >
                              <div>{token.text}</div>
                            </div>
                          </React.Fragment>
                        )
                      })}
                    </div>
                  ) : null}
                  <div className={translationClass(textFocus)}>{tokenizedText.translation}</div>
                  {textFocus && termsFocus ? (
                    <TermsComponent
                      token={termsComponentProps.token}
                      usage={termsComponentProps.usage}
                    />
                  ) : null}
                </div>
              </div>
            )
          })
        )}
      </div>
    </>
  )
}

interface ITermsComponentProps {
  token: Token
  usage: NoteUsage
}

interface ITermsShowData {
  terms: Term[]
}

const pageSize = 5

// eslint-disable-next-line max-lines-per-function
const TermsComponent: React.FC<ITermsComponentProps> = ({ token, usage }) => {
  const fetcher = useFetcher<ITermsShowData>()
  const terms = useMemo<Term[]>(() => (fetcher.data ? fetcher.data.terms : []), [fetcher.data])

  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const termRefs = useRef<(HTMLDivElement | null)[]>([])

  const [page, setPage] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, pageSize), [terms])

  const [createNoteData, setCreateNoteData] = useState<CreateNoteData | null>(null)
  const onCloseCreateNote = () => setCreateNoteData(null)

  useEffect(() => {
    if (fetcher.state !== "idle" || fetcher.data) return
    fetcher.load(`/terms/search?${queryString({ query: token.text, pos: token.partOfSpeech })}`)
  }, [fetcher, token])

  useEffect(() => {
    const termElement = termRefs.current[termFocusIndex]
    if (!termElement) return
    termElement.focus()
  }, [terms, page, termFocusIndex]) // trigger from terms/page to do initial focus

  useEffect(() => {
    openModal = createNoteData !== null
  }, [createNoteData])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (openModal) return

      switch (e.code) {
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
          setCreateNoteData(createNoteDataFromTerm(terms[termFocusIndex], usage))
          break
        default:
          return
      }
      e.preventDefault()
    },
    [termFocusIndex, page, pagesLen, terms, usage]
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
                {term.commonLevel !== CommonLevel.Unique && (
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
            {Array(pagesLen)
              .fill(null)
              .map((_, index) => (
                /* eslint-disable-next-line react/no-array-index-key */
                <span key={index} className={index === page ? "" : "text-light"}>
                  {index === page ? <>&#x2716;</> : <>&bull;</>}
                </span>
              ))}
          </div>
        </div>
      )}

      {createNoteData !== null && (
        <SlideOver.Dialog show onClose={onCloseCreateNote}>
          <SlideOver.Header title="Create Note" onClose={onCloseCreateNote} />
          <NoteForm data={createNoteData} onClose={onCloseCreateNote} />
        </SlideOver.Dialog>
      )}
    </div>
  )
}

export default SourceShow
