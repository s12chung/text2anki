/* eslint-disable max-lines */
import { CommonLevel } from "../../services/models/Lang.ts"
import {
  CreateNoteData,
  createNoteDataFromSourceTerm,
  createNoteDataFromUsage,
} from "../../services/models/Note.ts"
import {
  PosPunctuation,
  Source,
  SourcePart,
  Token,
  TokenizedText,
  tokenPreviousPunct,
  tokenPreviousSpace,
} from "../../services/models/Source.ts"
import { Term } from "../../services/models/Term.ts"
import { joinClasses, paginate, scrollTo } from "../../utils/HtmlUtil.ts"
import { preventDefault, SafeSet, useSafeSet } from "../../utils/JSXUtil.ts"
import { queryString } from "../../utils/RequestUtil.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import SlideOver from "../SlideOver.tsx"
import NoteCreate from "../notes/NoteCreate.tsx"
import { StopKeyboardContext, useStopKeyboard } from "./SourceShow_SourceComponent.ts"
import { getUsage, useFocusTextWithKeyboard } from "./SourceShow_SourceNavComponent.ts"
import { otherTranslationTexts, useChangeTermWithKeyboard } from "./SourceShow_TermsComponent.ts"
import { SelectedToken, useFocusTokenWithKeyboard } from "./SourceShow_TokensComponent.ts"
import {
  PartCreateForm,
  PartUpdateForm,
  SourceDetailMenu,
  SourceEditHeader,
  SourcePartDetailMenu,
} from "./SourceShow_Update.tsx"
import React, { SyntheticEvent, useContext, useEffect, useMemo, useRef, useState } from "react"
import { useFetcher } from "react-router-dom"

export interface ISourceShowData {
  source: Promise<Source>
}
interface ISourceShowProps {
  readonly data: ISourceShowData
}

const SourceShow: React.FC<ISourceShowProps> = ({ data }) => {
  return (
    <AwaitWithFallback resolve={data.source}>
      {(source: Source) => <SourceComponent source={source} />}
    </AwaitWithFallback>
  )
}

const SourceComponent: React.FC<{ readonly source: Source }> = ({ source }) => {
  const stopKeyboard = useStopKeyboard()
  const [nav, setNav] = useState<boolean>(true)
  const [edit, setEdit] = useState<boolean>(false)
  const [expandPartsCreate, setExpandPartsCreate] = useState<boolean>(false)

  const safeSet = useSafeSet((val) => stopKeyboard.setStopKeyboardEvents(val), [stopKeyboard])
  useEffect(() => {
    safeSet.addReset(() => {
      setEdit(false)
      setExpandPartsCreate(false)
    })
  }, [safeSet])

  return (
    <StopKeyboardContext.Provider value={stopKeyboard}>
      <div className="grid-std">
        {edit ? (
          <SourceEditHeader source={source} onCancel={() => safeSet.safeSet(setEdit, false)} />
        ) : (
          <SourceShowHeader
            source={source}
            onAddParts={() => safeSet.safeSet(setExpandPartsCreate, true)}
            onEdit={() => safeSet.safeSet(setEdit, true)}
          />
        )}
        <div className="flex justify-center mt-std mb-10">
          <a href="#" className="btn" onClick={preventDefault(() => setNav(!nav))}>
            Read Korean
          </a>
        </div>
      </div>

      {nav ? (
        <SourceNavComponent source={source} safeSet={safeSet} />
      ) : (
        <SourceShowComponent source={source} />
      )}
      <div className="grid-std pt-std pb-std2">
        {expandPartsCreate ? (
          <PartCreateForm
            sourceId={source.id}
            onCancel={() => safeSet.safeSet(setExpandPartsCreate, false)}
          />
        ) : (
          <div className="flex justify-center">
            <button
              type="button"
              className="btn"
              onClick={preventDefault(() => safeSet.safeSet(setExpandPartsCreate, true))}
            >
              Create Part
            </button>
          </div>
        )}
      </div>
    </StopKeyboardContext.Provider>
  )
}

const SourceShowHeader: React.FC<{
  readonly source: Source
  readonly onAddParts: () => void
  readonly onEdit: () => void
}> = ({ source, onAddParts, onEdit }) => {
  return (
    <div className="flex">
      <div className="flex-grow">
        <h2>{source.name}</h2>
        <div>{source.reference}</div>
      </div>
      <SourceDetailMenu source={source} onAddParts={onAddParts} onEdit={onEdit} />
    </div>
  )
}

const SourcePartsWrapper: React.FC<{
  readonly sourceId: number
  readonly parts: SourcePart[]
  readonly safeSet?: SafeSet
  readonly children: (
    tokenizedText: TokenizedText,
    partIndex: number,
    textIndex: number,
  ) => React.ReactNode
}> = ({ sourceId, parts, safeSet, children }) => {
  return (
    <div className="text-center space-y-std2">
      {parts.map((part, partIndex) => (
        <SourcePartWrapper
          // eslint-disable-next-line react/no-array-index-key
          key={`part-${partIndex}`}
          // eslint-disable-next-line react/no-children-prop
          children={children}
          sourceId={sourceId}
          partIndex={partIndex}
          part={part}
          safeSet={safeSet}
        />
      ))}
    </div>
  )
}

SourcePartsWrapper.defaultProps = {
  // eslint-disable-next-line no-undefined
  safeSet: undefined,
}

const SourcePartWrapper: React.FC<{
  readonly sourceId: number
  readonly partIndex: number
  readonly part: SourcePart
  readonly safeSet?: SafeSet
  readonly children: (
    tokenizedText: TokenizedText,
    partIndex: number,
    textIndex: number,
  ) => React.ReactNode
}> = ({ sourceId, partIndex, part, safeSet, children }) => {
  const [edit, setEdit] = useState<boolean>(false)
  useEffect(() => {
    // eslint-disable-next-line no-empty-function
    if (!safeSet) return () => {}
    const key = safeSet.addReset(() => setEdit(false))
    return () => safeSet.removeReset(key)
  }, [safeSet])

  return (
    <div className="group relative">
      {edit ? (
        <PartUpdateForm
          sourceId={sourceId}
          partIndex={partIndex}
          part={part}
          onCancel={() => safeSet?.safeSet(setEdit, false)}
        />
      ) : (
        <>
          {Boolean(safeSet) && (
            <SourcePartDetailMenu
              sourceId={sourceId}
              partIndex={partIndex}
              onEdit={() => safeSet?.safeSet(setEdit, true)}
            />
          )}
          {part.tokenizedTexts.map((tokenizedText, textIndex) => (
            /* eslint-disable-next-line react/no-array-index-key */
            <div key={`${tokenizedText.text}-${textIndex}`}>
              {Boolean(tokenizedText.previousBreak) && <div className="text-4xl">&nbsp;</div>}
              {children(tokenizedText, partIndex, textIndex)}
            </div>
          ))}
        </>
      )}
      {part.media.imageUrl ? (
        <div className="grid-std">
          <img src={part.media.imageUrl} alt="Part Image" />
        </div>
      ) : (
        <hr className="mt-std2" />
      )}
    </div>
  )
}

SourcePartWrapper.defaultProps = {
  // eslint-disable-next-line no-undefined
  safeSet: undefined,
}

const textClassBase = "ko-sans text-2xl"
const translationClassBase = "text-lg"

const SourceShowComponent: React.FC<{ readonly source: Source }> = ({ source }) => {
  return (
    <SourcePartsWrapper sourceId={source.id} parts={source.parts}>
      {(tokenizedText) => (
        <>
          <div className={textClassBase}>{tokenizedText.text}</div>
          <div className={translationClassBase}>{tokenizedText.translation}</div>
        </>
      )}
    </SourcePartsWrapper>
  )
}

// eslint-disable-next-line max-lines-per-function
const SourceNavComponent: React.FC<{ readonly source: Source; readonly safeSet: SafeSet }> = ({
  source,
  safeSet,
}) => {
  const [createNoteData, setCreateNoteData] = useState<CreateNoteData | null>(null)
  const { setStopKeyboardEvents } = useContext(StopKeyboardContext)
  useEffect(
    () => setStopKeyboardEvents(createNoteData !== null),
    [createNoteData, setStopKeyboardEvents],
  )

  const [selectedToken, setSelectedToken] = useState<SelectedToken | null>(null)
  const isTokenSelected = selectedToken !== null
  const [partFocusIndex, textFocusIndex, setText] = useFocusTextWithKeyboard(
    source.parts,
    isTokenSelected,
    () =>
      setCreateNoteData(createNoteDataFromUsage(getUsage(source, partFocusIndex, textFocusIndex))),
    () => setSelectedToken(null),
  )
  const setTextOnClick = (
    e: SyntheticEvent<HTMLDivElement>,
    textFocused: boolean,
    partIndex: number,
    textIndex: number,
    // eslint-disable-next-line max-params
  ) => {
    if (textFocused) return
    e.preventDefault()
    setText(partIndex, textIndex)
    setCustomToken(null)
    setSelectedToken(null)
  }

  const [customToken, setCustomToken] = useState<SelectedToken | null>(null)
  const onCustomToken = (token: SelectedToken | null) => {
    setCustomToken(token)
    setSelectedToken(null)
  }
  const onTokenSearch = (token: SelectedToken) => {
    setCustomToken(null)
    setSelectedToken(token)
  }

  return (
    <>
      {createNoteData !== null && (
        <CreateNoteDialog
          data={createNoteData}
          onCreate={() => {
            setCreateNoteData(null)
            setSelectedToken(null)
          }}
          onClose={() => setCreateNoteData(null)}
        />
      )}
      <SourcePartsWrapper sourceId={source.id} parts={source.parts} safeSet={safeSet}>
        {(tokenizedText, partIndex, textIndex) => {
          const textFocused = partIndex === partFocusIndex && textIndex === textFocusIndex
          return (
            <div
              className={joinClasses(textFocused ? "py-4 bg-gray-std" : "", "group py-2")}
              onClick={(e) => setTextOnClick(e, textFocused, partIndex, textIndex)}
            >
              <div className={joinClasses(textClassBase, textFocused ? "text-light" : "")}>
                {tokenizedText.text}
              </div>
              {textFocused ? (
                <TokensComponent
                  tokens={tokenizedText.tokens}
                  isTokenSelected={isTokenSelected}
                  onTokenSelect={setSelectedToken}
                  onCustomToken={onCustomToken}
                />
              ) : null}
              <div className={textFocused ? "text-2xl" : translationClassBase}>
                {tokenizedText.translation}
              </div>
              {textFocused && selectedToken ? (
                <TermsComponent
                  token={selectedToken}
                  onTermSelect={(term: Term) => {
                    setCreateNoteData(
                      createNoteDataFromSourceTerm(
                        term,
                        getUsage(source, partFocusIndex, textFocusIndex),
                      ),
                    )
                  }}
                />
              ) : null}
              {textFocused && customToken ? (
                <SearchTokensComponent customToken={customToken} onTokenSearch={onTokenSearch} />
              ) : null}
            </div>
          )
        }}
      </SourcePartsWrapper>
    </>
  )
}

const TokensComponent: React.FC<{
  readonly tokens: Token[]
  readonly isTokenSelected: boolean
  readonly onTokenSelect: (token: SelectedToken) => void
  readonly onCustomToken: (token: SelectedToken | null) => void
}> = ({ tokens, isTokenSelected, onTokenSelect, onCustomToken }) => {
  const [tokenFocusIndex] = useFocusTokenWithKeyboard(
    tokens,
    isTokenSelected,
    onTokenSelect,
    onCustomToken,
  )

  const tokenRefs = useRef<(HTMLDivElement | null)[]>([])
  useEffect(() => {
    const element = tokenRefs.current[tokenFocusIndex]
    if (element) scrollTo(element)
  }, [tokenFocusIndex])

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
                index === tokenFocusIndex ? "text-white bg-ink" : "",
                isPunct ? "text-faded" : "",
              )}
              /* eslint-disable-next-line no-undefined */
              tabIndex={isPunct ? undefined : -1}
            >
              <div>{token.text}</div>
            </div>
          </React.Fragment>
        )
      })}
    </div>
  )
}

const SearchTokensComponent: React.FC<{
  readonly customToken: SelectedToken
  readonly onTokenSearch: (token: SelectedToken) => void
}> = ({ customToken, onTokenSearch }) => {
  useEffect(() => {
    // keyboard trigger to show component adds a character, so timeout
    const id = setTimeout(() => textRef.current?.focus(), 50)
    return () => clearTimeout(id)
  }, [])

  const { setStopKeyboardEvents } = useContext(StopKeyboardContext)
  useEffect(() => {
    setStopKeyboardEvents(true)
    return () => setStopKeyboardEvents(false)
  }, [setStopKeyboardEvents])

  const textRef = useRef<HTMLInputElement>(null)
  return (
    <form
      className="mt-std space-x-basic"
      onSubmit={preventDefault(() =>
        onTokenSearch({ text: textRef.current?.value || "", partOfSpeech: "" }),
      )}
    >
      <input ref={textRef} defaultValue={customToken.text} />
      <button type="submit" className="btn-primary">
        Submit
      </button>
    </form>
  )
}

interface ITermsShowData {
  terms: Term[]
}

const termsComponentClass = "grid-std text-left text-lg py-2 space-y-2"

const TermsComponent: React.FC<{
  readonly token: SelectedToken
  readonly onTermSelect: (term: Term) => void
}> = ({ token, onTermSelect }) => {
  const fetcher = useFetcher<ITermsShowData>()
  const terms = useMemo<Term[]>(() => (fetcher.data ? fetcher.data.terms : []), [fetcher.data])
  useEffect(() => {
    if (fetcher.state !== "idle" || fetcher.data) return
    fetcher.load(`/terms/search?${queryString({ query: token.text, pos: token.partOfSpeech })}`)
  }, [fetcher, token])

  const [termFocusIndex, pageIndex, pagesLen, maxPageSize, shake] = useChangeTermWithKeyboard(
    terms,
    onTermSelect,
  )

  if (!fetcher.data) return <div className={termsComponentClass}>Loading...</div>
  return (
    <div className={termsComponentClass}>
      {terms.length === 0 ? (
        <div>No terms found</div>
      ) : (
        <div className={shake ? "shake" : ""}>
          {paginate(terms, maxPageSize, pageIndex).map((term, index) => (
            <div
              key={term.id}
              className={joinClasses(
                index === termFocusIndex ? "underline focus-ring" : "",
                "py-1",
              )}
            >
              <div className="text-xl">
                <span className="font-bold">{term.text}</span>&nbsp;
                <span className="text-light text-base">{term.partOfSpeech}</span>
                {term.commonLevel !== CommonLevel.Unique && (
                  <span className="relative top-2">&nbsp;{"*".repeat(term.commonLevel)}</span>
                )}
                : {term.translations[0].text}&nbsp;&mdash;&nbsp;
                {term.translations[0].explanation}
              </div>
              <div className="ml-std2">{otherTranslationTexts(term.translations)}</div>
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
    </div>
  )
}

const CreateNoteDialog: React.FC<{
  readonly data: CreateNoteData
  readonly onCreate: () => void
  readonly onClose: () => void
}> = ({ data, onCreate, onClose }) => {
  return (
    <SlideOver.Dialog show onClose={onClose}>
      <SlideOver.Header title="Create Note" onClose={onClose} />
      <NoteCreate data={data} onCreate={onCreate} onClose={onClose} />
    </SlideOver.Dialog>
  )
}

export default SourceShow
