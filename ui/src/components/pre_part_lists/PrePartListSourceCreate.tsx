import { PrePartList } from "../../services/PrePartListsService.ts"
import { imageToClipboard } from "../../utils/ClipboardUtils.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import AwaitError from "../AwaitError.tsx"
import Header from "../Header.tsx"
import SlideOver from "../SlideOver.tsx"
import React, { useCallback, useEffect, useRef, useState } from "react"
import { Await, Form } from "react-router-dom"

export interface IPrePartListSourceCreateData {
  prePartList: Promise<PrePartList>
}

interface IPrePartListSourceCreateProps {
  data: IPrePartListSourceCreateData
}

const PrePartListSourceCreate: React.FC<IPrePartListSourceCreateProps> = ({ data }) => {
  return (
    <React.Suspense fallback={<div>Loading....</div>}>
      <Await resolve={data.prePartList} errorElement={<AwaitError />}>
        {(prePartList: PrePartList) => <PrePartsForm prePartList={prePartList} />}
      </Await>
    </React.Suspense>
  )
}

// eslint-disable-next-line max-lines-per-function
const PrePartsForm: React.FC<{ prePartList: PrePartList }> = ({ prePartList }) => {
  const { preParts } = prePartList
  const [currentIndex, setCurrentIndex] = useState<number>(0)
  const [partTextsMap, setPartTextsMap] = useState<Record<number, string>>({})

  const setCurrentIndexWithClipboard = useCallback(
    (changeFunction: (index: number, length: number) => number) => {
      const index = changeFunction(currentIndex, preParts.length)
      setCurrentIndex(index)
      const { imageUrl } = preParts[index]
      if (!imageUrl) return
      imageToClipboard(imageUrl)
    },
    [currentIndex, preParts]
  )
  const next = useCallback(
    () => setCurrentIndexWithClipboard(increment),
    [setCurrentIndexWithClipboard]
  )
  const prev = useCallback(
    () => setCurrentIndexWithClipboard(decrement),
    [setCurrentIndexWithClipboard]
  )
  const setPartTextsAt = (index: number, value: string) => {
    const c = { ...partTextsMap }
    c[index] = value
    setPartTextsMap(c)
  }

  const textAreaRefs = useRef<(HTMLTextAreaElement | null)[]>([])
  useEffect(() => {
    const textArea = textAreaRefs.current[currentIndex]
    if (!textArea) return
    textArea.focus()
  }, [currentIndex])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      switch (e.code) {
        case "F1":
          prev()
          break
        case "F2":
          next()
          break
        default:
          return
      }
      e.preventDefault()
    },
    [next, prev]
  )

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  return (
    <SlideOver.Dialog show leftNode={<PrePartLeft image={preParts[currentIndex].imageUrl ?? ""} />}>
      <SlideOver.Header title="Create Source from Parts" />
      <Form action="/sources" method="post" className="m-std space-y-std">
        <div className="text-center">
          Part {currentIndex + 1}/{preParts.length}
        </div>
        <div className="flex">
          <button type="button" className="btn flex-grow" onClick={prev}>
            ←
          </button>
          <button type="button" className="btn flex-grow" onClick={next}>
            →
          </button>
        </div>

        <input type="hidden" name="prePartListId" value={prePartList.id} />

        {preParts.map((prePart, index) => (
          <textarea
            key={prePart.imageUrl}
            ref={(ref) => (textAreaRefs.current[index] = ref)}
            autoFocus={index === 0}
            name={`parts[${index}].text`}
            value={partTextsMap[index]}
            placeholder="Each text line followed by an optional translation"
            rows={20}
            className={`text-xl${index === currentIndex ? "" : " hidden"}`}
            onChange={(e) => setPartTextsAt(index, e.target.value)}
          />
        ))}

        <div className="flex justify-end">
          <button type="submit" className="btn-primary">
            Create Source from Parts
          </button>
        </div>
      </Form>
    </SlideOver.Dialog>
  )
}

const PrePartLeft: React.FC<{ image: string }> = ({ image }) => (
  <div className="h-screen flex flex-col">
    <div className="m-std">
      <Header />
    </div>
    <img className="flex-grow" src={image} alt="Drag and Dropped image" />
  </div>
)

export default PrePartListSourceCreate
