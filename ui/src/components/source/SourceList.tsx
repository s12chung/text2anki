import { Source } from "../../services/SourceService"
import AwaitError from "../AwaitError.tsx"
import React from "react"
import { Await } from "react-router-dom"

export interface ISourceListData {
  sources: Promise<Source[]>
}
interface ISourceListProps {
  data: ISourceListData
}

const SourceList: React.FC<ISourceListProps> = ({ data }) => {
  return (
    <React.Suspense fallback={<div>Loading....</div>}>
      <Await resolve={data.sources} errorElement={<AwaitError />}>
        {(sources: Source[]) =>
          sources.length === 0 ? (
            <div>No sources created</div>
          ) : (
            <ul>
              {sources.map((source) => (
                <li key={`source-${source.id}`}>{source.name}</li>
              ))}
            </ul>
          )
        }
      </Await>
    </React.Suspense>
  )
}

export default SourceList
