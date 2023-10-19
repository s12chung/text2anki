import NotificationsContext from "../../contexts/NotificationsContext.ts"
import { CommonLevel } from "../../services/models/Lang.ts"
import { CreateNoteData, CreateNoteDataEmpty, Note } from "../../services/models/Note.ts"
import { filterKeys } from "../../utils/ArrayUntil.ts"
import { camelToTitle } from "../../utils/StringUtil.ts"
import React, { useContext, useEffect, useRef, useState } from "react"
import { useFetcher } from "react-router-dom"

interface INoteFormData {
  note: Promise<Note>
}

const termKeys: (keyof CreateNoteData)[] = ["text", "partOfSpeech", "translation", "explanation"]
const usageKeys: (keyof CreateNoteData)[] = ["usage", "usageTranslation"]
const commonLevelKey: keyof CreateNoteData = "commonLevel"
const otherKeys = filterKeys(
  Object.keys(CreateNoteDataEmpty) as (keyof CreateNoteData)[],
  usageKeys,
  termKeys,
  [commonLevelKey],
)

const NoteCreate: React.FC<{ readonly data: CreateNoteData; readonly onClose: () => void }> = ({
  data,
  onClose,
}) => {
  const fetcher = useFetcher<INoteFormData>()
  const { error, success } = useContext(NotificationsContext)
  useEffect(() => {
    if (!fetcher.data) return
    fetcher.data.note
      .then((note) => success(`Created new Note: ${note.text}`))
      .catch(() => error("Failed to create Note"))
      .finally(() => onClose())
  }, [fetcher, onClose, success, error])

  const [submitted, setSubmitted] = useState<boolean>(false)
  const submitButtonRef = useRef<HTMLButtonElement>(null)
  useEffect(() => submitButtonRef.current?.focus(), [])

  return (
    <fetcher.Form
      action="/notes"
      method="post"
      className="m-std space-y-std2"
      onSubmit={() => setSubmitted(true)}
    >
      <div className="space-y-half">
        <TextFormGroup dataKeys={termKeys} data={data} />
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
      </div>

      <TextFormGroup dataKeys={usageKeys} data={data} />
      <TextFormGroup dataKeys={otherKeys} data={data} />

      <div className="my-std justify-end flex space-x-basic">
        <button type="button" className="btn" onClick={onClose}>
          Cancel
        </button>
        <button ref={submitButtonRef} type="submit" className="btn-primary" disabled={submitted}>
          Create
        </button>
      </div>
    </fetcher.Form>
  )
}

const TextFormGroup: React.FC<{
  readonly dataKeys: (keyof CreateNoteData)[]
  readonly data: CreateNoteData
}> = ({ dataKeys, data }) => (
  <div className="space-y-half">
    {dataKeys.map((key) => (
      <TextFormField key={key} data={data} dataKey={key} />
    ))}
  </div>
)

const TextFormField: React.FC<{
  readonly data: CreateNoteData
  readonly dataKey: keyof CreateNoteData
}> = ({ data, dataKey }) => {
  return (
    <div>
      <label>{camelToTitle(dataKey)}</label>
      <input name={dataKey} type="text" defaultValue={data[dataKey]} />
    </div>
  )
}

export default NoteCreate
