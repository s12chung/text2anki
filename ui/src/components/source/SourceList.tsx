import { Source } from "../../services/SourceService.ts"
import AwaitError from "../AwaitError.tsx"
import { DocumentTextIcon } from "@heroicons/react/24/outline"
import React from "react"
import { Await, Link } from "react-router-dom"

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
                <li key={`source-${source.id}`} className="my-half">
                  <Link to={`sources/${source.id}`} className="flex-std">
                    <DocumentTextIcon className="h-6 w-6" />
                    <div>{source.name}</div>
                  </Link>
                </li>
              ))}
            </ul>
          )
        }
      </Await>
    </React.Suspense>
  )
}

export default SourceList
