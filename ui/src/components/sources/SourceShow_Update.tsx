import NotificationsContext from "../../contexts/NotificationsContext.ts"
import { partString, Source, SourcePart } from "../../services/models/Source.ts"
import { joinClasses, menuClass } from "../../utils/HtmlUtil.ts"
import { preventDefault } from "../../utils/JSXUtil.ts"
import DetailMenu from "../DetailMenu.tsx"
import PrePartListDragAndDrop from "../pre_part_lists/PrePartListDragAndDrop.tsx"
import { Menu } from "@headlessui/react"
import React, { useContext, useEffect, useRef, useState } from "react"
import { Form, useFetcher } from "react-router-dom"

interface ISourceResponse {
  source: Source
}

export const SourceDetailMenu: React.FC<{
  readonly source: Source
  readonly onAddParts: () => void
  readonly onEdit: () => void
}> = ({ source, onAddParts, onEdit }) => {
  return (
    <Form
      action={`/sources/${source.id}`}
      method="delete"
      onSubmit={(event) => {
        // eslint-disable-next-line no-alert
        if (!window.confirm("Delete Source?")) event.preventDefault()
      }}
    >
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

export const PartsCreate: React.FC<{
  readonly sourceId: number
  readonly expand: boolean
  readonly setExpand: (expand: boolean) => void
}> = ({ sourceId, expand, setExpand }) => {
  return (
    <div className="grid-std pt-std pb-std2">
      {expand ? (
        <PartCreateForm sourceId={sourceId} onCancel={() => setExpand(false)} />
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

export const PartCreateForm: React.FC<{
  readonly sourceId: number
  readonly onCancel: () => void
}> = ({ sourceId, onCancel }) => {
  const [text, setText] = useState<string>("")
  const textAreaRef = useRef<HTMLTextAreaElement>(null)
  useEffect(() => textAreaRef.current?.focus(), [textAreaRef])

  const fetcher = useFetcher<ISourceResponse>()
  const { success } = useContext(NotificationsContext)
  useEffect(() => {
    if (!fetcher.data) return
    success(`Created Part`)
    onCancel()
  }, [fetcher, success, onCancel])

  return (
    <PrePartListDragAndDrop sourceId={sourceId} minHeight="h-third">
      <fetcher.Form action={`/sources/${sourceId}/parts`} method="post">
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
      </fetcher.Form>
    </PrePartListDragAndDrop>
  )
}

export const SourceEditHeader: React.FC<{
  readonly source: Source
  readonly onCancel: () => void
}> = ({ source, onCancel }) => {
  const fetcher = useFetcher<ISourceResponse>()
  const { success } = useContext(NotificationsContext)
  useEffect(() => {
    if (!fetcher.data) return
    success(`Updated Source`)
    onCancel()
  }, [fetcher, success, onCancel])

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

export const SourcePartDetailMenu: React.FC<{
  readonly sourceId: number
  readonly partIndex: number
  readonly onEdit: () => void
}> = ({ sourceId, partIndex, onEdit }) => {
  const fetcher = useFetcher<ISourceResponse>()
  const { success } = useContext(NotificationsContext)
  const didRun = useRef(false)

  useEffect(() => {
    if (!fetcher.data || didRun.current) return
    didRun.current = true
    success(`Deleted Part`)
  }, [fetcher, success])

  return (
    <fetcher.Form
      action={`/sources/${sourceId}/parts/${partIndex}`}
      method="delete"
      className="group-hover:block hidden absolute top-0 right-0"
    >
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
            <a href="#" className={menuClass(active)} onClick={preventDefault(onEdit)}>
              Edit
            </a>
          )}
        </Menu.Item>
      </DetailMenu>
    </fetcher.Form>
  )
}

export const PartUpdateForm: React.FC<{
  readonly sourceId: number
  readonly partIndex: number
  readonly part: SourcePart
  readonly onCancel: () => void
}> = ({ sourceId, partIndex, part, onCancel }) => {
  const fetcher = useFetcher<ISourceResponse>()
  const { success } = useContext(NotificationsContext)
  useEffect(() => {
    if (!fetcher.data) return
    success(`Updated Part`)
    onCancel()
  }, [fetcher, success, onCancel])

  const textAreaRef = useRef<HTMLTextAreaElement | null>(null)
  useEffect(() => textAreaRef.current?.focus(), [])

  return (
    <fetcher.Form
      action={`/sources/${sourceId}/parts/${partIndex}`}
      method="patch"
      className="grid-std"
    >
      <textarea ref={textAreaRef} name="text" className="h-third" defaultValue={partString(part)} />
      <div className="mt-half flex justify-end space-x-basic">
        <button type="button" className="btn" onClick={preventDefault(onCancel)}>
          Cancel
        </button>
        <button type="submit" className="btn-primary">
          Update Part
        </button>
      </div>
    </fetcher.Form>
  )
}
