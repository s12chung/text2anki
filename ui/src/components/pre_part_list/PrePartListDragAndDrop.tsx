import {
  prePartListService,
  PrePartSignData,
  PreSignedHTTPRequest,
} from "../../services/PrePartListsService.ts"
import { printError } from "../../utils/ErrorUtil.ts"
import { headers } from "../../utils/RequestUtil.ts"
import { XMarkIcon } from "@heroicons/react/24/outline"
import React, {
  DragEventHandler,
  MouseEventHandler,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react"
import { useNavigate } from "react-router-dom"
import { DotLoader } from "react-spinners"

enum DragState {
  None,
  Dragging,
  Dropped,
}

const PrePartListDragAndDrop: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [files, setFiles] = useState<File[]>([])
  const [dragState, setDragState] = useState<DragState>(DragState.None)

  const onClose = useCallback(() => setDragState(DragState.None), [])
  const onCloseMouse: MouseEventHandler<HTMLAnchorElement> = (e) => {
    e.preventDefault()
    onClose()
  }
  const handleKeyDown = useCallback(
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

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [handleKeyDown])

  const handleDragOver: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault()
    setDragState(DragState.Dragging)
  }

  const handleDragLeave: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault()
    setDragState(DragState.None)
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    // eslint-disable-next-line prefer-destructuring
    const files = e.dataTransfer.files
    if (files.length === 0) {
      return
    }

    setFiles(Array.from(files).sort((a: File, b: File) => a.name.localeCompare(b.name)))
    setDragState(DragState.Dropped)
  }

  return (
    <div
      className="min-h-screen"
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      {dragState === DragState.None ? (
        children
      ) : (
        <div className="min-h-screen flex relative">
          <a href="#" className="absolute top-5 right-5 a-btn" onClick={onCloseMouse}>
            <XMarkIcon className="h-10 w-10" />
          </a>
          {dragState === DragState.Dragging ? (
            <div className="m-auto text-4xl">Dragging files to create Source</div>
          ) : (
            <PrePartListDrop files={files} />
          )}
        </div>
      )}
    </div>
  )
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
      .filter((imageRequest): imageRequest is PreSignedHTTPRequest => Boolean(imageRequest))
      .map((imageRequest, index) =>
        fetch(imageRequest.url, {
          method: imageRequest.method,
          headers: headers(imageRequest.signedHeader),
          body: files[index],
        })
      )
  ).then(() => signedResponse.id)
}

const PrePartListDrop: React.FC<{ files: File[] }> = ({ files }) => {
  const navigate = useNavigate()

  const didRun = useRef(false)
  const [errorMessage, setErrorMessage] = useState<string>("")

  useEffect(() => {
    if (didRun.current || files.length === 0) return
    didRun.current = true
    uploadFiles(files)
      .then((id) => navigate(`/sources/pre_part_lists/${id}`))
      .catch((error) => setErrorMessage(printError(error).message))
  }, [files, navigate])

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
      <div className="mt-8">
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

export default PrePartListDragAndDrop
