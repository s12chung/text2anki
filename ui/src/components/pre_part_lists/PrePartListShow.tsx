import { printAndAlertError } from "../../services/Format.ts"
import { PrePartList } from "../../services/PrePartListsService.ts"
import { imageToClipboard } from "../../utils/ClipboardUtils.ts"
import { decrement, increment } from "../../utils/NumberUtil.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import Header from "../Header.tsx"
import SlideOver from "../SlideOver.tsx"
import React, { useCallback, useEffect, useRef, useState } from "react"
import { Form } from "react-router-dom"

export interface IPrePartListShowData {
  prePartList: Promise<PrePartList>
  sourceId: number
}

interface IPrePartListShowProps {
  data: IPrePartListShowData
}

const PrePartListShow: React.FC<IPrePartListShowProps> = ({ data }) => {
  return (
    <AwaitWithFallback resolve={data.prePartList}>
      {(prePartList: PrePartList) => (
        <PrePartsForm sourceId={data.sourceId} prePartList={prePartList} />
      )}
    </AwaitWithFallback>
  )
}

// eslint-disable-next-line max-lines-per-function
const PrePartsForm: React.FC<{ sourceId: number; prePartList: PrePartList }> = ({
  sourceId,
  prePartList,
}) => {
  const formAction = sourceId ? `/sources/${sourceId}/parts/multi` : "/sources"
  const action = sourceId ? "Append Parts" : "Create Source from Parts"
  const { preParts } = prePartList
  const [currentIndex, setCurrentIndex] = useState<number>(0)
  const [partTextsMap, setPartTextsMap] = useState<Record<number, string>>({})

  const setImageToClipboard = useCallback(
    (index: number) => {
      const { imageUrl } = preParts[index]
      if (!imageUrl) return
      imageToClipboard(imageUrl).catch((error) => {
        printAndAlertError(error)
      })
    },
    [preParts]
  )
  const setCurrentIndexWithClipboard = useCallback(
    (changeFunction: (index: number, length: number) => number) => {
      const index = changeFunction(currentIndex, preParts.length)
      setCurrentIndex(index)
      setImageToClipboard(index)
    },
    [currentIndex, preParts.length, setImageToClipboard]
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
    setImageToClipboard(currentIndex)

    const textArea = textAreaRefs.current[currentIndex]
    if (!textArea) return
    textArea.focus()
  }, [currentIndex, setImageToClipboard])

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
    <SlideOver.Dialog
      show
      leftNode={<PrePartLeft image={preParts[currentIndex].imageUrl} prev={prev} next={next} />}
    >
      <SlideOver.Header title={action} />
      <Form action={formAction} method="post" className="m-std space-y-std">
        <div className="text-center">
          Part {currentIndex + 1}/{preParts.length}
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
            {action}
          </button>
        </div>
      </Form>
    </SlideOver.Dialog>
  )
}

const PrePartLeft: React.FC<{ image: string; prev: () => void; next: () => void }> = ({
  image,
  prev,
  next,
}) => (
  <div className="h-screen flex flex-col">
    <div className="m-std">
      <Header />
    </div>
    <div className="flex flex-grow relative">
      <ImageNav char="<" changeF={prev} />
      <div className="flex-1" />
      <ImageNav char=">" changeF={next} />
      <img
        className="absolute flex-grow h-full w-full object-contain -z-10"
        src={image}
        alt="Drag and Dropped image"
      />
    </div>
  </div>
)

const ImageNav: React.FC<{ char: string; changeF: () => void }> = ({ char, changeF }) => {
  return (
    <a
      href="#"
      className="flex flex-1 bg-black justify-center items-center opacity-0 hover:opacity-50 transition ease-out duration-300"
      onClick={(e) => {
        e.preventDefault()
        changeF()
      }}
    >
      <span className="text-white text-8xl font-bold">{char}</span>
    </a>
  )
}

export default PrePartListShow
