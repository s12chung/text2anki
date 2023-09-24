/* eslint-disable max-lines */
import NotificationsContext from "../../contexts/NotificationsContext.ts"
import { CommonLevel } from "../../services/Lang.ts"
import {
  CreateNoteData,
  createNoteDataFromSourceTerm,
  NoteUsage,
} from "../../services/NotesService.ts"
import {
  PosPunctuation,
  Source,
  Token,
  TokenizedText,
  tokenPreviousPunct,
  tokenPreviousSpace,
} from "../../services/SourcesService.ts"
import { Term } from "../../services/TermsService.ts"
import { unique } from "../../utils/ArrayUntil.ts"
import { joinClasses, menuClass, pageSize, paginate, totalPages } from "../../utils/HtmlUtil.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import { queryString } from "../../utils/RequestUtil.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import DetailMenu from "../DetailMenu.tsx"
import SlideOver from "../SlideOver.tsx"
import NoteCreate from "../notes/NoteCreate.tsx"
import { Menu } from "@headlessui/react"
import React, {
  MouseEventHandler,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react"
import { Form, useFetcher } from "react-router-dom"

export interface ISourceShowData {
  source: Promise<Source>
}
interface ISourceShowProps {
  data: ISourceShowData
}

const SourceShow: React.FC<ISourceShowProps> = ({ data }) => {
  return (
    <AwaitWithFallback resolve={data.source}>
      {(source: Source) => <SourceComponent source={source} />}
    </AwaitWithFallback>
  )
}

let stopKeyboardEvents = false

const SourceComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [nav, setNav] = useState<boolean>(true)
  const [show, setShow] = useState<boolean>(true)

  const onEdit: () => void = () => setShow(false)
  const onCancel: () => void = () => setShow(true)

  const onReadKorean: MouseEventHandler<HTMLAnchorElement> = (e) => {
    e.preventDefault()
    setNav(!nav)
  }

  return (
    <>
      <div className="grid-std">
        {show ? (
          <SourceShowHeader source={source} onEdit={onEdit} />
        ) : (
          <SourceEditHeader source={source} onCancel={onCancel} />
        )}
        <div className="flex justify-center mt-std mb-10">
          <a href="#" className="btn" onClick={onReadKorean}>
            Read Korean
          </a>
        </div>
      </div>

      {nav ? <SourceNavComponent source={source} /> : <SourceShowComponent source={source} />}
    </>
  )
}

const SourceShowHeader: React.FC<{
  source: Source
  onEdit: () => void
}> = ({ source, onEdit }) => {
  const onEditClick: MouseEventHandler<HTMLAnchorElement> = (e) => {
    e.preventDefault()
    onEdit()
  }
  return (
    <Form
      action={`/sources/${source.id}`}
      method="delete"
      className="flex-std"
      onSubmit={(event) => {
        // eslint-disable-next-line no-alert
        if (!window.confirm("Delete Source?")) event.preventDefault()
      }}
    >
      <div className="flex-grow">
        <h2>{source.name}</h2>
        <div>{source.reference}</div>
      </div>

      <DetailMenu>
        <Menu.Item>
          {({ active }) => (
            <button type="submit" className={menuClass(active)}>
              Delete
            </button>
          )}
        </Menu.Item>
        <Menu.Item>
          {({ active }) => (
            <a href="#" className={menuClass(active)} onClick={onEditClick}>
              Edit
            </a>
          )}
        </Menu.Item>
      </DetailMenu>
    </Form>
  )
}

const SourceEditHeader: React.FC<{
  source: Source
  onCancel: () => void
}> = ({ source, onCancel }) => {
  useEffect(() => {
    stopKeyboardEvents = true
  }, [])
  const onCancelWrap = useCallback(() => {
    stopKeyboardEvents = false
    onCancel()
  }, [onCancel])

  const fetcher = useFetcher<Source>()
  const { error, success } = useContext(NotificationsContext)

  useEffect(() => {
    if (!fetcher.data) return
    success(`Updated Source`)
    onCancelWrap()
  }, [fetcher, success, error, onCancelWrap])

  const onCancelClick: MouseEventHandler<HTMLAnchorElement> = (e) => {
    e.preventDefault()
    onCancelWrap()
  }

  return (
    <fetcher.Form action={`/sources/${source.id}`} method="patch" className="space-y-std">
      <label>
        Name:
        <input name="name" type="text" defaultValue={source.name} />
      </label>
      <label>
        Reference:
        <input name="reference" type="text" defaultValue={source.reference} />
      </label>

      <div className="flex justify-end space-x-basic">
        <a href="#" className="btn" onClick={onCancelClick}>
          Cancel
        </a>
        <button type="submit" className="btn-primary">
          Save
        </button>
      </div>
    </fetcher.Form>
  )
}

const SourceWrapper: React.FC<{
  source: Source
  children: (tokenizedText: TokenizedText, partIndex: number, textIndex: number) => React.ReactNode
}> = ({ source, children }) => {
  return (
    <div className="text-center">
      {source.parts.map((part, partIndex) => (
        // eslint-disable-next-line react/no-array-index-key
        <div key={`part-${partIndex}`}>
          {part.tokenizedTexts.map((tokenizedText, textIndex) => (
            /* eslint-disable-next-line react/no-array-index-key */
            <div key={`${tokenizedText.text}-${textIndex}`}>
              {Boolean(tokenizedText.previousBreak) && <div className="text-4xl">&nbsp;</div>}
              {children(tokenizedText, partIndex, textIndex)}
            </div>
          ))}
          {part.media.imageUrl ? (
            <div className="grid-std">
              <img src={part.media.imageUrl} alt="Part Image" />
            </div>
          ) : null}
        </div>
      ))}
    </div>
  )
}

const textClassBase = "ko-sans text-2xl focgrin:text-light"
const translationClassBase = "text-lg focgrin:text-2xl"

const SourceShowComponent: React.FC<{ source: Source }> = ({ source }) => {
  return (
    <SourceWrapper source={source}>
      {(tokenizedText) => (
        <>
          <div className={textClassBase}>{tokenizedText.text}</div>
          <div className={translationClassBase}>{tokenizedText.translation}</div>
        </>
      )}
    </SourceWrapper>
  )
}

function getTermsComponentProps(
  source: Source,
  tokenizedText: TokenizedText,
  tokenFocusIndex: number
): ITermsComponentProps {
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

// eslint-disable-next-line max-lines-per-function
const SourceNavComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [partFocusIndex, setPartFocusIndex] = useState<number>(0)
  const [textFocusIndex, setTextFocusIndex] = useState<number>(0)
  const [termsComponentProps, setTermsComponentProps] = useState<ITermsComponentProps | null>(null)

  const textRefs = useRef<(HTMLDivElement | null)[][]>([])
  const [lastFocusedElement, setLastFocusedElement] = useState<HTMLDivElement | null>(null)
  const focusElement = (element: HTMLDivElement) => {
    setLastFocusedElement(element)
    element.focus()
  }

  const currentTokenizedTexts = useMemo<TokenizedText[]>(
    () => source.parts[partFocusIndex].tokenizedTexts,
    [partFocusIndex, source.parts]
  )
  const termsFocused = termsComponentProps !== null
  const onTokenChange = (tokenElement: HTMLDivElement) => focusElement(tokenElement)
  const onTokenSelect = (tokenFocusIndex: number) => {
    setTermsComponentProps(
      getTermsComponentProps(source, currentTokenizedTexts[textFocusIndex], tokenFocusIndex)
    )
  }

  const partsLength = source.parts.length
  const decrementText = useCallback(() => {
    const result = decrement(textFocusIndex, currentTokenizedTexts.length)
    if (result !== currentTokenizedTexts.length - 1) {
      setTextFocusIndex(result)
      return
    }
    const partIndex = decrement(partFocusIndex, partsLength)
    setPartFocusIndex(partIndex)
    setTextFocusIndex(source.parts[partIndex].tokenizedTexts.length - 1)
  }, [textFocusIndex, currentTokenizedTexts.length, partFocusIndex, partsLength, source.parts])
  const incrementText = useCallback(() => {
    const result = increment(textFocusIndex, currentTokenizedTexts.length)
    if (result !== 0) {
      setTextFocusIndex(result)
      return
    }
    const partIndex = increment(partFocusIndex, partsLength)
    setPartFocusIndex(partIndex)
    setTextFocusIndex(0)
  }, [textFocusIndex, currentTokenizedTexts.length, partFocusIndex, partsLength])

  useEffect(() => {
    setTermsComponentProps(null)
    const textElement = textRefs.current[partFocusIndex][textFocusIndex]
    if (!textElement) return
    focusElement(textElement)
    window.scrollTo({
      top: textElement.getBoundingClientRect().top + window.scrollY - 150,
      behavior: "smooth",
    })
  }, [partFocusIndex, textFocusIndex])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return

      switch (e.code) {
        case "Escape":
          if (!termsFocused) return
          lastFocusedElement?.focus()
          setTermsComponentProps(null)
          break
        default:
      }

      if (termsFocused) return

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
    [termsFocused, lastFocusedElement, decrementText, incrementText]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const textOnClick = (index: number) => setTextFocusIndex(index)

  return (
    <SourceWrapper source={source}>
      {(tokenizedText, partIndex, textIndex) => {
        const textFocused = partIndex === partFocusIndex && textIndex === textFocusIndex
        return (
          <div
            ref={(ref) => {
              if (!textRefs.current[partIndex]) textRefs.current[partIndex] = []
              textRefs.current[partIndex][textIndex] = ref
            }}
            tabIndex={-1}
            className={joinClasses(
              textFocused ? "py-4 bg-gray-std" : "",
              "group py-2 focin:py-4 focin:bg-gray-std"
            )}
            onClick={() => textOnClick(textIndex)}
          >
            <div className={joinClasses(textClassBase, textFocused ? "text-light" : "")}>
              {tokenizedText.text}
            </div>
            {textFocused ? (
              <TokensComponent
                tokens={tokenizedText.tokens}
                termsFocused={termsFocused}
                onTokenChange={onTokenChange}
                onTokenSelect={onTokenSelect}
              />
            ) : null}
            <div className={joinClasses(translationClassBase, textFocused ? "text-2xl" : "")}>
              {tokenizedText.translation}
            </div>
            {textFocused && termsFocused ? (
              <TermsComponent token={termsComponentProps.token} usage={termsComponentProps.usage} />
            ) : null}
          </div>
        )
      }}
    </SourceWrapper>
  )
}

function skipPunct(
  tokens: Token[],
  index: number,
  change: (index: number, length: number) => number
): number {
  index = change(index, tokens.length)
  if (tokens[index].partOfSpeech !== PosPunctuation) return index
  return skipPunct(tokens, index, change)
}

const TokensComponent: React.FC<{
  tokens: Token[]
  termsFocused: boolean
  onTokenChange: (tokenElement: HTMLDivElement) => void
  onTokenSelect: (tokenFocusIndex: number) => void
}> = ({ tokens, termsFocused, onTokenChange, onTokenSelect }) => {
  const [tokenFocusIndex, setTokenFocusIndex] = useState<number>(0)
  const tokenRefs = useRef<(HTMLDivElement | null)[]>([])

  const isAllPunct = useMemo<boolean>(
    () => tokens.every((token) => token.partOfSpeech === PosPunctuation),
    [tokens]
  )

  useEffect(() => {
    const element = tokenRefs.current[tokenFocusIndex]
    if (element) onTokenChange(element)
  }, [onTokenChange, tokenFocusIndex])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents || termsFocused || isAllPunct) return

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
          onTokenSelect(tokenFocusIndex)
          break
        default:
          return
      }

      e.preventDefault()
    },
    [termsFocused, isAllPunct, tokens, tokenFocusIndex, onTokenSelect]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const tokenOnClick = (index: number) => setTokenFocusIndex(index)

  return (
    <div className="ko-sans text-4xl justify-center mb-2 child:py-2 flex">
      {tokens.map((token, index) => {
        const previousSpace = tokenPreviousSpace(tokens, index)
        const isPunct = token.partOfSpeech === PosPunctuation
        return (
          /* eslint-disable-next-line react/no-array-index-key */
          <React.Fragment key={`${token.text}-${token.partOfSpeech}-${index}`}>
            {!previousSpace && index !== 0 && (
              <div className={isPunct || tokenPreviousPunct(tokens, index) ? "text-faded" : ""}>
                &middot;
              </div>
            )}
            {Boolean(previousSpace) && index !== 0 && <div>&nbsp;&nbsp;</div>}
            <div
              ref={(ref) => (tokenRefs.current[index] = ref)}
              className={joinClasses(
                "focus:text-white focus:bg-ink",
                index === tokenFocusIndex ? " text-white bg-ink" : "",
                isPunct ? " text-faded" : ""
              )}
              /* eslint-disable-next-line no-undefined */
              tabIndex={isPunct ? undefined : -1}
              /* eslint-disable-next-line no-undefined */
              onClick={isPunct ? undefined : () => tokenOnClick(index)}
            >
              <div>{token.text}</div>
            </div>
          </React.Fragment>
        )
      })}
    </div>
  )
}

interface ITermsComponentProps {
  token: Token
  usage: NoteUsage
}

interface ITermsShowData {
  terms: Term[]
}

const maxPageSize = 5

// eslint-disable-next-line max-lines-per-function
const TermsComponent: React.FC<ITermsComponentProps> = ({ token, usage }) => {
  const fetcher = useFetcher<ITermsShowData>()
  const terms = useMemo<Term[]>(() => (fetcher.data ? fetcher.data.terms : []), [fetcher.data])

  const [termFocusIndex, setTermFocusIndex] = useState<number>(0)
  const termRefs = useRef<(HTMLDivElement | null)[]>([])

  const [pageIndex, setPageIndex] = useState<number>(0)
  const pagesLen = useMemo<number>(() => totalPages(terms, maxPageSize), [terms])

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
  }, [terms, pageIndex, termFocusIndex]) // trigger from terms/page to do initial focus

  useEffect(() => {
    stopKeyboardEvents = createNoteData !== null
  }, [createNoteData])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (stopKeyboardEvents) return

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
          setCreateNoteData(createNoteDataFromSourceTerm(terms[termFocusIndex], usage))
          break
        default:
          return
      }
      e.preventDefault()
    },
    [termFocusIndex, pageIndex, pagesLen, terms, usage]
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
          {paginate(terms, maxPageSize, pageIndex).map((term, index) => (
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
              <div className="ml-std2">
                {unique(term.translations.map((translation) => translation.text))
                  .slice(1, 6)
                  .join("; ")}
              </div>
            </div>
          ))}
          <div className="text-center space-x-half">
            {Array(pagesLen)
              .fill(null)
              .map((_, index) => (
                /* eslint-disable-next-line react/no-array-index-key */
                <span key={index} className={index === pageIndex ? "" : "text-light"}>
                  {index === pageIndex ? <>&#x2716;</> : <>&bull;</>}
                </span>
              ))}
          </div>
        </div>
      )}

      {createNoteData !== null && (
        <SlideOver.Dialog show onClose={onCloseCreateNote}>
          <SlideOver.Header title="Create Note" onClose={onCloseCreateNote} />
          <NoteCreate data={createNoteData} onClose={onCloseCreateNote} />
        </SlideOver.Dialog>
      )}
    </div>
  )
}

export default SourceShow
