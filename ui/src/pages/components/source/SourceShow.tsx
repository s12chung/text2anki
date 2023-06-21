import { Source, sourceService } from "../../../services/SourceService.ts"
import { printAndAlertError } from "../../../utils/ErrorUtil.ts"
import AwaitError from "../AwaitError.tsx"
import { Ring } from "@uiball/loaders"
import React, { ChangeEventHandler, FormEventHandler, MouseEventHandler, useState } from "react"
import { Await } from "react-router-dom"

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

const SourceComponent: React.FC<{ source: Source }> = ({ source }) => {
  return (
    <div>
      <SourceNameComponent source={source} />

      {source.tokenizedTexts.map((tokenizedText) => (
        <div key={`text-${tokenizedText.text}`}>
          <div>{tokenizedText.text}</div>
          <div>{tokenizedText.translation}</div>
        </div>
      ))}
    </div>
  )
}

const SourceNameComponent: React.FC<{ source: Source }> = ({ source }) => {
  const [isEditing, setIsEditing] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [editName, setEditName] = useState(source.name)

  const handleEdit: MouseEventHandler = () => setIsEditing(true)
  const handleInput: ChangeEventHandler<HTMLInputElement> = (e) => setEditName(e.target.value)

  const handleSubmit: FormEventHandler = (e) => {
    e.preventDefault()
    setIsLoading(true)

    sourceService
      .update(source.id, { name: editName })
      .then((s) => {
        source.name = s.name
      })
      .catch((error) => printAndAlertError(error))
      .finally(() => {
        setIsEditing(false)
        setIsLoading(false)
      })
  }

  const handleCancel: MouseEventHandler = () => {
    setIsEditing(false)
    setEditName(source.name)
  }

  return isEditing ? (
    <form className="flex" onSubmit={handleSubmit}>
      <input className="flex-grow" type="text" value={editName} onChange={handleInput} />

      {Boolean(isLoading) && <Ring />}
      <button type="button" disabled={isLoading} onClick={handleCancel}>
        Cancel
      </button>
      <button type="submit" disabled={isLoading}>
        Submit
      </button>
    </form>
  ) : (
    <div className="flex">
      <div className="flex-grow">{source.name}</div>
      <button type="button" onClick={handleEdit}>
        Edit
      </button>
    </div>
  )
}

export default SourceShow
