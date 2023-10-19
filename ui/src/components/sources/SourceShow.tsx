/* eslint-disable max-lines */
import { CommonLevel } from "../../services/models/Lang.ts"
import { CreateNoteData, createNoteDataFromSourceTerm } from "../../services/models/Note.ts"
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
import {
  getTermProps,
  ITermsComponentProps,
  useFocusTextWithKeyboard,
} from "./SourceShow_SourceNavComponent.ts"
import { otherTranslationTexts, useChangeTermWithKeyboard } from "./SourceShow_TermsComponent.ts"
import { useFocusTokenWithKeyboard } from "./SourceShow_TokensComponent.ts"
import {
  PartCreateForm,
  PartUpdateForm,
  SourceDetailMenu,
  SourceEditHeader,
  SourcePartDetailMenu,
} from "./SourceShow_Update.tsx"
import React, { useEffect, useMemo, useRef, useState } from "react"
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

const SourceNavComponent: React.FC<{ readonly source: Source; readonly safeSet: SafeSet }> = ({
  source,
  safeSet,
}) => {
  const [termProps, setTermProps] = useState<ITermsComponentProps | null>(null)

  const termsFocused = termProps !== null
  const [partFocusIndex, textFocusIndex, focusElement, setText] = useFocusTextWithKeyboard(
    source.parts,
    termsFocused,
    () => setTermProps(null),
  )

  const textRefs = useRef<(HTMLDivElement | null)[][]>([])
  useEffect(() => {
    setTermProps(null)
    const textElement = textRefs.current[partFocusIndex][textFocusIndex]
    if (!textElement) return
    focusElement(textElement)
    scrollTo(textElement)
  }, [focusElement, partFocusIndex, textFocusIndex])

  return (
    <SourcePartsWrapper sourceId={source.id} parts={source.parts} safeSet={safeSet}>
      {(tokenizedText, partIndex, textIndex) => {
        const textFocused = partIndex === partFocusIndex && textIndex === textFocusIndex
        return (
          <div
            ref={(ref) => {
              if (!textRefs.current[partIndex]) textRefs.current[partIndex] = []
              textRefs.current[partIndex][textIndex] = ref
            }}
            tabIndex={-1}
            className={joinClasses(textFocused ? "py-4 bg-gray-std" : "", "group py-2")}
            onClick={preventDefault(() => setText(partIndex, textIndex))}
          >
            <div className={joinClasses(textClassBase, textFocused ? "text-light" : "")}>
              {tokenizedText.text}
            </div>
            {textFocused ? (
              <TokensComponent
                tokens={tokenizedText.tokens}
                termsFocused={termsFocused}
                onTokenSelect={(tokenIndex) =>
                  setTermProps(getTermProps(source, partFocusIndex, textFocusIndex, tokenIndex))
                }
                onTokenChange={(tokenElement) => focusElement(tokenElement)}
              />
            ) : null}
            <div className={textFocused ? "text-2xl" : translationClassBase}>
              {tokenizedText.translation}
            </div>
            {textFocused && termsFocused ? (
              <TermsComponent token={termProps.token} usage={termProps.usage} />
            ) : null}
          </div>
        )
      }}
    </SourcePartsWrapper>
  )
}

const TokensComponent: React.FC<{
  readonly tokens: Token[]
  readonly termsFocused: boolean
  readonly onTokenSelect: (tokenFocusIndex: number) => void
  readonly onTokenChange: (tokenElement: HTMLDivElement) => void
}> = ({ tokens, termsFocused, onTokenSelect, onTokenChange }) => {
  const [tokenFocusIndex] = useFocusTokenWithKeyboard(tokens, termsFocused, onTokenSelect)

  const tokenRefs = useRef<(HTMLDivElement | null)[]>([])
  useEffect(() => {
    const element = tokenRefs.current[tokenFocusIndex]
    if (element) onTokenChange(element)
  }, [onTokenChange, tokenFocusIndex])

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

interface ITermsShowData {
  terms: Term[]
}

const termsComponentClass = "grid-std text-left text-lg py-2 space-y-2"

const TermsComponent: React.FC<ITermsComponentProps> = ({ token, usage }) => {
  const fetcher = useFetcher<ITermsShowData>()
  const terms = useMemo<Term[]>(() => (fetcher.data ? fetcher.data.terms : []), [fetcher.data])
  useEffect(() => {
    if (fetcher.state !== "idle" || fetcher.data) return
    fetcher.load(`/terms/search?${queryString({ query: token.text, pos: token.partOfSpeech })}`)
  }, [fetcher, token])

  const [createNoteData, setCreateNoteData] = useState<CreateNoteData | null>(null)
  const [termFocusIndex, pageIndex, pagesLen, maxPageSize, shake] = useChangeTermWithKeyboard(
    terms,
    (term: Term) => setCreateNoteData(createNoteDataFromSourceTerm(term, usage)),
    () => createNoteData !== null,
  )

  const termRefs = useRef<(HTMLDivElement | null)[]>([])
  useEffect(() => termRefs.current[termFocusIndex]?.focus(), [terms, pageIndex, termFocusIndex])

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
              ref={(ref) => (termRefs.current[index] = ref)}
              tabIndex={-1}
              className={joinClasses(index === termFocusIndex ? "underline" : "", "py-1")}
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

      {createNoteData !== null && (
        <NoteDialog data={createNoteData} onClose={() => setCreateNoteData(null)} />
      )}
    </div>
  )
}

const NoteDialog: React.FC<{ readonly data: CreateNoteData; readonly onClose: () => void }> = ({
  data,
  onClose,
}) => {
  return (
    <SlideOver.Dialog show onClose={onClose}>
      <SlideOver.Header title="Create Note" onClose={onClose} />
      <NoteCreate data={data} onClose={onClose} />
    </SlideOver.Dialog>
  )
}

export default SourceShow
