/* eslint-disable max-lines */
import NotificationsContext from "../../contexts/NotificationsContext.ts"
import { CommonLevel } from "../../services/Lang.ts"
import { CreateNoteData, createNoteDataFromSourceTerm } from "../../services/NotesService.ts"
import {
  PosPunctuation,
  Source,
  SourcePart,
  Token,
  TokenizedText,
  tokenPreviousPunct,
  tokenPreviousSpace,
} from "../../services/SourcesService.ts"
import { Term } from "../../services/TermsService.ts"
import { joinClasses, menuClass, paginate, scrollTo } from "../../utils/HtmlUtil.ts"
import { preventDefault } from "../../utils/JSXUtil.ts"
import { queryString } from "../../utils/RequestUtil.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import DetailMenu from "../DetailMenu.tsx"
import SlideOver from "../SlideOver.tsx"
import NoteCreate from "../notes/NoteCreate.tsx"
import PrePartListDragAndDrop from "../pre_part_lists/PrePartListDragAndDrop.tsx"
import { StopKeyboardContext, useStopKeyboard } from "./SourceShow_SourceComponent.ts"
import {
  getTermProps,
  ITermsComponentProps,
  useFocusTextWithKeyboard,
} from "./SourceShow_SourceNavComponent.ts"
import { otherTranslationTexts, useChangeTermWithKeyboard } from "./SourceShow_TermsComponent.ts"
import { useFocusTokenWithKeyboard } from "./SourceShow_TokensComponent.ts"
import { Menu } from "@headlessui/react"
import React, { useContext, useEffect, useMemo, useRef, useState } from "react"
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

const SourceComponent: React.FC<{ source: Source }> = ({ source }) => {
  const stopKeyboardValue = useStopKeyboard()
  const [nav, setNav] = useState<boolean>(true)
  const [show, setShow] = useState<boolean>(true)
  const [expandPartsCreate, setExpandPartsCreate] = useState<boolean>(false)

  const resetAndSet = (set: (val: boolean) => void, val: boolean) => {
    setShow(true)
    setExpandPartsCreate(false)
    stopKeyboardValue.setStopKeyboardEvents(val)
    set(val)
  }

  return (
    <StopKeyboardContext.Provider value={stopKeyboardValue}>
      <div className="grid-std">
        {show ? (
          <SourceShowHeader
            source={source}
            onAddParts={() => resetAndSet(setExpandPartsCreate, true)}
            onEdit={() => resetAndSet(setShow, false)}
          />
        ) : (
          <SourceEditHeader source={source} onCancel={() => resetAndSet(setShow, false)} />
        )}
        <div className="flex justify-center mt-std mb-10">
          <a href="#" className="btn" onClick={preventDefault(() => setNav(!nav))}>
            Read Korean
          </a>
        </div>
      </div>

      {nav ? <SourceNavComponent source={source} /> : <SourceShowComponent source={source} />}
      <PartsCreate
        sourceId={source.id}
        expand={expandPartsCreate}
        setExpand={(val: boolean) => resetAndSet(setExpandPartsCreate, val)}
      />
    </StopKeyboardContext.Provider>
  )
}

const SourceShowHeader: React.FC<{
  source: Source
  onAddParts: () => void
  onEdit: () => void
}> = ({ source, onAddParts, onEdit }) => {
  return (
    <Form
      action={`/sources/${source.id}`}
      method="delete"
      className="flex"
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
            <button type="submit" className={joinClasses("w-full", menuClass(active))}>
              Delete
            </button>
          )}
        </Menu.Item>
        <Menu.Item>
          {({ active }) => (
            <a href="#" className={menuClass(active)} onClick={preventDefault(onAddParts)}>
              Add Parts
            </a>
          )}
        </Menu.Item>
        <Menu.Item>
          {({ active }) => (
            <a href="#" className={menuClass(active)} onClick={preventDefault(onEdit)}>
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
  const fetcher = useFetcher<Source>()
  const { error, success } = useContext(NotificationsContext)

  useEffect(() => {
    if (!fetcher.data) return
    success(`Updated Source`)
    onCancel()
  }, [fetcher, success, error, onCancel])

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
        <button type="button" className="btn" onClick={preventDefault(onCancel)}>
          Cancel
        </button>
        <button type="submit" className="btn-primary">
          Save
        </button>
      </div>
    </fetcher.Form>
  )
}

const SourcePartsWrapper: React.FC<{
  parts: SourcePart[]
  children: (tokenizedText: TokenizedText, partIndex: number, textIndex: number) => React.ReactNode
}> = ({ parts, children }) => {
  return (
    <div className="text-center space-y-std2">
      {parts.map((part, partIndex) => (
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
          ) : (
            partIndex !== parts.length - 1 && <hr className="mt-std2" />
          )}
        </div>
      ))}
    </div>
  )
}

const textClassBase = "ko-sans text-2xl"
const translationClassBase = "text-lg"

const SourceShowComponent: React.FC<{ source: Source }> = ({ source }) => {
  return (
    <SourcePartsWrapper parts={source.parts}>
      {(tokenizedText) => (
        <>
          <div className={textClassBase}>{tokenizedText.text}</div>
          <div className={translationClassBase}>{tokenizedText.translation}</div>
        </>
      )}
    </SourcePartsWrapper>
  )
}

const SourceNavComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [termProps, setTermProps] = useState<ITermsComponentProps | null>(null)

  const termsFocused = termProps !== null
  const [partFocusIndex, textFocusIndex, focusElement, setText] = useFocusTextWithKeyboard(
    source.parts,
    termsFocused,
    () => setTermProps(null)
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
    <SourcePartsWrapper parts={source.parts}>
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
  tokens: Token[]
  termsFocused: boolean
  onTokenSelect: (tokenFocusIndex: number) => void
  onTokenChange: (tokenElement: HTMLDivElement) => void
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
                isPunct ? "text-faded" : ""
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
  const [termFocusIndex, pageIndex, pagesLen, maxPageSize] = useChangeTermWithKeyboard(
    terms,
    (term: Term) => setCreateNoteData(createNoteDataFromSourceTerm(term, usage)),
    () => createNoteData !== null
  )

  const termRefs = useRef<(HTMLDivElement | null)[]>([])
  useEffect(() => termRefs.current[termFocusIndex]?.focus(), [terms, pageIndex, termFocusIndex])

  if (!fetcher.data) return <div className={termsComponentClass}>Loading...</div>
  return (
    <div className={termsComponentClass}>
      {terms.length === 0 ? (
        <div>No terms found</div>
      ) : (
        <div>
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

const NoteDialog: React.FC<{ data: CreateNoteData; onClose: () => void }> = ({ data, onClose }) => {
  return (
    <SlideOver.Dialog show onClose={onClose}>
      <SlideOver.Header title="Create Note" onClose={onClose} />
      <NoteCreate data={data} onClose={onClose} />
    </SlideOver.Dialog>
  )
}

const PartsCreate: React.FC<{
  sourceId: number
  expand: boolean
  setExpand: (expand: boolean) => void
}> = ({ sourceId, expand, setExpand }) => {
  return (
    <div className="grid-std pt-std pb-std2">
      {expand ? (
        <PartsCreateForm sourceId={sourceId} onCancel={() => setExpand(false)} />
      ) : (
        <div className="flex justify-center">
          <button type="button" className="btn" onClick={preventDefault(() => setExpand(true))}>
            Create Part
          </button>
        </div>
      )}
    </div>
  )
}

const PartsCreateForm: React.FC<{ sourceId: number; onCancel: () => void }> = ({
  sourceId,
  onCancel,
}) => {
  const [text, setText] = useState<string>("")
  const textAreaRef = useRef<HTMLTextAreaElement>(null)
  useEffect(() => textAreaRef.current?.focus(), [textAreaRef])

  return (
    <PrePartListDragAndDrop sourceId={sourceId} minHeight="h-third">
      <Form action={`/sources/${sourceId}/parts`} method="post">
        <textarea
          ref={textAreaRef}
          name="text"
          value={text}
          placeholder="You may also drag and drop here."
          className="h-third"
          onChange={(e) => setText(e.target.value)}
        />
        <div className="mt-half flex justify-end space-x-basic">
          <button type="button" className="btn" onClick={preventDefault(onCancel)}>
            Cancel
          </button>
          <button type="submit" className="btn-primary" disabled={!text}>
            Add Part
          </button>
        </div>
      </Form>
    </PrePartListDragAndDrop>
  )
}

export default SourceShow
