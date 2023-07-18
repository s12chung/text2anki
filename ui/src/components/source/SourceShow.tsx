import { Source, Token } from "../../services/SourceService.ts"
import AwaitError from "../AwaitError.tsx"
import React, { useEffect, useRef, useState } from "react"
import { Await, Link } from "react-router-dom"

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
  const textRefs = useRef<(HTMLDivElement | null)[]>([])
  const tokenRefs = useRef<(HTMLDivElement | null)[][]>([])

  // eslint-disable-next-line prefer-destructuring
  const tokenizedTexts = source.parts[0].tokenizedTexts

  useEffect(() => {
    textRefs.current[textFocusIndex]?.focus()
    if (tokenFocusIndex === -1) return
    tokenRefs.current[textFocusIndex][tokenFocusIndex]?.focus()
  }, [textFocusIndex, tokenFocusIndex])

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const textLen = tokenizedTexts.length
      const tokenLen = tokenizedTexts[textFocusIndex]?.tokens.length

      let changedTextFocusIndex: number | null = null
      let changedTokenFocusIndex: number | null = null
      switch (event.code) {
        case "ArrowUp":
        case "KeyW":
          changedTextFocusIndex = decrement(textFocusIndex, textLen)
          break
        case "ArrowDown":
        case "KeyS":
          changedTextFocusIndex = increment(textFocusIndex, textLen)
          break
        case "ArrowLeft":
        case "KeyA":
          changedTokenFocusIndex = decrement(tokenFocusIndex, tokenLen)
          break
        case "ArrowRight":
        case "KeyD":
          changedTokenFocusIndex = increment(tokenFocusIndex, tokenLen)
          break
        default:
          return
      }
      event.preventDefault()

      if (changedTextFocusIndex !== null) {
        setTokenFocusIndex(-1)
        setTextFocusIndex(changedTextFocusIndex)
      }
      if (changedTokenFocusIndex !== null) {
        if (textFocusIndex === -1) setTextFocusIndex(0)
        setTokenFocusIndex(changedTokenFocusIndex)
      }
    }

    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [tokenizedTexts, textFocusIndex, tokenFocusIndex])

  const handleTextClick = (index: number) => setTextFocusIndex(index)
  const handleTokenClick = (index: number) => setTokenFocusIndex(index)

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
        {tokenizedTexts.map((tokenizedText, textIndex) => (
          /* eslint-disable-next-line react/no-array-index-key */
          <div key={`${tokenizedText.text}-${textIndex}`}>
            {Boolean(tokenizedText.previousBreak) && <div className="text-4xl">&nbsp;</div>}
            <div
              ref={(ref) => (textRefs.current[textIndex] = ref)}
              tabIndex={-1}
              className="group py-2 focin:py-4 focin:bg-gray-std"
              onClick={() => handleTextClick(textIndex)}
            >
              <div className="ko-sans text-2xl focgrin:text-light">{tokenizedText.text}</div>
              <div className="ko-sans hidden text-4xl mb-4 justify-center focgrin:flex">
                {tokenizedText.tokens.map((token, index) => {
                  const previousSpace = tokenPreviousSpace(tokenizedText.tokens, index)
                  return (
                    /* eslint-disable-next-line react/no-array-index-key */
                    <div key={`${token.text}-${token.partOfSpeech}-${index}`} className="flex">
                      {!previousSpace && index !== 0 && <span>&middot;</span>}
                      {Boolean(previousSpace) && index !== 0 && <span>&nbsp;&nbsp;</span>}
                      <div
                        ref={(ref) => {
                          if (!tokenRefs.current[textIndex]) tokenRefs.current[textIndex] = []
                          tokenRefs.current[textIndex][index] = ref
                        }}
                        tabIndex={-1}
                        onClick={() => handleTokenClick(index)}
                      >
                        <div>{token.text}</div>
                      </div>
                    </div>
                  )
                })}
              </div>
              <div className="text-lg focgrin:text-2xl">{tokenizedText.translation}</div>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}

export default SourceShow
