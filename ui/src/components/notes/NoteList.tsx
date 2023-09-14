import { Note, notesService } from "../../services/NotesService.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import React, { MouseEventHandler } from "react"

export interface INoteListData {
  sources: Promise<Note[]>
}
interface INoteListProps {
  data: INoteListData
}

const SourceList: React.FC<INoteListProps> = ({ data }) => {
  const onDownload: MouseEventHandler<HTMLAnchorElement> = () => {
    setTimeout(() => {
      window.location.reload()
    }, 1000)
  }

  return (
    <>
      <div className="flex justify-center mb-10">
        <a download href={notesService.downloadUrl()} className="btn" onClick={onDownload}>
          Download Notes
        </a>
      </div>

      <AwaitWithFallback resolve={data.sources}>
        {(notes: Note[]) =>
          notes.length === 0 ? (
            <div>No notes created</div>
          ) : (
            <ul>
              {notes.map((note) => (
                <li key={`source-${note.id}`} className={note.downloaded ? "text-faded" : ""}>
                  <div>
                    {note.text} - {note.usage}
                  </div>
                  <div className="ml-std2">{note.usageTranslation}</div>
                </li>
              ))}
            </ul>
          )
        }
      </AwaitWithFallback>
    </>
  )
}

export default SourceList
