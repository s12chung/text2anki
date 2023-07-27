import NotificationsContext from "../../contexts/NotificationsContext.ts"
import { CommonLevel } from "../../services/LangService.ts"
import { CreateNoteData, CreateNoteDataEmpty, Note } from "../../services/NotesService.ts"
import { camelToTitle } from "../../utils/StringUtil.ts"
import React, { useContext, useEffect, useRef, useState } from "react"
import { useFetcher } from "react-router-dom"

interface INoteFormData {
  note: Promise<Note>
}

const NoteForm: React.FC<{ data: CreateNoteData; onClose: () => void }> = ({ data, onClose }) => {
  const fetcher = useFetcher<INoteFormData>()
  const commonLevelIndex = 3
  const createNoteDataKeys = Object.keys(CreateNoteDataEmpty) as (keyof CreateNoteData)[]
  const commonLevelKey = createNoteDataKeys[commonLevelIndex]
  const { error, success } = useContext(NotificationsContext)

  const [submitted, setSubmitted] = useState<boolean>(false)

  const submitButtonRef = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    submitButtonRef.current?.focus()
  }, [])

  useEffect(() => {
    if (!fetcher.data) return
    fetcher.data.note
      .then((note) => success(`Created new Note: ${note.text}`))
      .catch(() => error("Failed to create Note"))
      .finally(() => onClose())
  }, [fetcher, onClose, success, error])

  return (
    <fetcher.Form
      action="/notes"
      method="post"
      className="m-std space-y-std"
      onSubmit={() => setSubmitted(true)}
    >
      {createNoteDataKeys.slice(0, commonLevelIndex).map((key) => (
        <TextFormField key={key} data={data} dataKey={key} />
      ))}

      <div className="space-x-std">
        <label>{camelToTitle(commonLevelKey)}</label>
        <select name={commonLevelKey} defaultValue={data[commonLevelKey]}>
          {Array(CommonLevel.Common + 1)
            .fill(null)
            .map((_, index) => (
              // eslint-disable-next-line react/no-array-index-key
              <option key={index} value={index}>
                {index}
              </option>
            ))}
        </select>
      </div>

      {createNoteDataKeys.slice(commonLevelIndex + 1).map((key) => (
        <TextFormField key={key} data={data} dataKey={key} />
      ))}

      <div className="my-std justify-end flex-std">
        <button type="button" onClick={onClose}>
          Cancel
        </button>
        <button ref={submitButtonRef} type="submit" disabled={submitted}>
          Create
        </button>
      </div>
    </fetcher.Form>
  )
}

const TextFormField: React.FC<{ data: CreateNoteData; dataKey: keyof CreateNoteData }> = ({
  data,
  dataKey,
}) => {
  return (
    <div>
      <label>{camelToTitle(dataKey)}</label>
      <input name={dataKey} type="text" defaultValue={data[dataKey]} />
    </div>
  )
}

export default NoteForm
