import { sourceIdQueryParam } from "../../controllers/PrePartListsController.ts"
import { printError } from "../../services/Format.ts"
import { prePartListService, PrePartSignData } from "../../services/PrePartListsService.ts"
import { Source } from "../../services/models/Source.ts"
import { joinClasses } from "../../utils/HtmlUtil.ts"
import { preventDefault, useKeyDownEffect } from "../../utils/JSXUtil.ts"
import { headers, queryString } from "../../utils/RequestUtil.ts"
import { removeExtension } from "../../utils/StringUtil.ts"
import { XMarkIcon } from "@heroicons/react/24/outline"
import React, { useCallback, useEffect, useRef, useState } from "react"
import { useFetcher, useNavigate } from "react-router-dom"
import { DotLoader } from "react-spinners"

enum DragState {
  None,
  Dragging,
  Dropped,
}

const textFileExts: Record<string, boolean> = {
  "text/plain": true,
  "text/markdown": true,
}

function useDropFiles(onDrop: () => void): readonly [File[], (e: React.DragEvent) => void] {
  const [files, setFiles] = useState<File[]>([])
  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      // eslint-disable-next-line prefer-destructuring
      const files = e.dataTransfer.files
      if (files.length === 0) {
        return
      }
      setFiles(Array.from(files).sort((a: File, b: File) => a.name.localeCompare(b.name)))
      onDrop()
    },
    [onDrop]
  )
  return [files, handleDrop] as const
}

const PrePartListDragAndDrop: React.FC<{
  sourceId?: number
  minHeight: string
  children: React.ReactNode
}> = ({ sourceId, minHeight, children }) => {
  const [dragState, setDragState] = useState<DragState>(DragState.None)

  const [files, handleDrop] = useDropFiles(() => setDragState(DragState.Dropped))
  const onClose = useCallback(() => setDragState(DragState.None), [])

  useKeyDownEffect(
    (e: KeyboardEvent) => {
      switch (e.code) {
        case "Escape":
          onClose()
          break
        default:
          return
      }
      e.preventDefault()
    },
    [onClose]
  )

  return (
    <div
      className={minHeight}
      onDragOver={preventDefault(() => setDragState(DragState.Dragging))}
      onDragLeave={preventDefault(() => setDragState(DragState.None))}
      onDrop={handleDrop}
    >
      {dragState === DragState.None ? (
        children
      ) : (
        <div className={joinClasses(minHeight, "flex relative")}>
          <a href="#" className="absolute top-5 right-5 a-btn" onClick={preventDefault(onClose)}>
            <XMarkIcon className="h-10 w-10" />
          </a>
          {dragState === DragState.Dragging ? (
            <div className="m-auto text-4xl">Dragging files to create Source</div>
          ) : (
            <PrePartListDrop sourceId={sourceId ? sourceId : 0} files={files} />
          )}
        </div>
      )}
    </div>
  )
}

PrePartListDragAndDrop.defaultProps = {
  sourceId: 0,
}

async function uploadFiles(files: File[]): Promise<string> {
  const exts = files.map((file) => {
    const splitName = file.name.split(".")
    return splitName.length > 1 ? `.${splitName[splitName.length - 1]}` : ""
  })
  const signedResponse = await prePartListService.sign({
    preParts: exts.map<PrePartSignData>((ext) => ({ imageExt: ext })),
  })

  return Promise.all<Response>(
    signedResponse.preParts
      .map((part) => part.imageRequest)
      .filter((imageRequest) => Boolean(imageRequest))
      .map((imageRequest, index) =>
        fetch(imageRequest.url, {
          method: imageRequest.method,
          headers: headers(imageRequest.signedHeader),
          body: files[index],
        })
      )
  ).then(() => signedResponse.id)
}

interface ISourceCreateResponse {
  source: Source
}

const PrePartListDrop: React.FC<{ sourceId: number; files: File[] }> = ({ sourceId, files }) => {
  const navigate = useNavigate()
  const fetcher = useFetcher<ISourceCreateResponse>()

  const didRun = useRef(false)
  const [errorMessage, setErrorMessage] = useState<string>("")

  useEffect(() => {
    if (didRun.current || files.length === 0 || onlyTextFile(files)) return
    didRun.current = true

    const query = sourceId ? `?${queryString({ [sourceIdQueryParam]: String(sourceId) })}` : ""
    uploadFiles(files)
      .then((id) => navigate(`/sources/pre_part_lists/${id}${query}`))
      .catch((error) => setErrorMessage(printError(error).message))
  }, [files, navigate, sourceId])

  useEffect(() => {
    const file = onlyTextFile(files)
    if (didRun.current || !file || fetcher.data || fetcher.state !== "idle") return

    didRun.current = true
    file
      .text()
      .then((text) => {
        fetcher.submit(
          { name: removeExtension(file.name), reference: file.name, "parts[0].text": text },
          { method: "post", action: "/sources" }
        )
      })
      .catch((error) => setErrorMessage(printError(error).message))
  }, [fetcher, files, navigate])

  return (
    <div className="m-auto text-2xl flex-col items-center text-center">
      Uploading files:
      <ul>
        {files.map((file, index) => (
          <li key={file.name}>
            {index + 1}. {file.name}
          </li>
        ))}
      </ul>
      <div className="mt-std2">
        {errorMessage ? (
          <>
            <p>Error: {errorMessage}</p>
            <p className="text-lg">Try dragging and dropping again</p>
          </>
        ) : (
          <DotLoader className="m-auto" />
        )}
      </div>
    </div>
  )
}

function onlyTextFile(files: File[]): File | null {
  if (files.length === 1) {
    // eslint-disable-next-line prefer-destructuring
    const file = files[0]
    if (textFileExts[file.type]) {
      return file
    }
  }
  return null
}

export default PrePartListDragAndDrop
