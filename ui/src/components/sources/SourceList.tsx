import { Source } from "../../services/models/Source.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import { DocumentTextIcon } from "@heroicons/react/24/outline"
import React from "react"
import { Link } from "react-router-dom"

export interface ISourceListData {
  sources: Promise<Source[]>
}
interface ISourceListProps {
  readonly data: ISourceListData
}

const SourceList: React.FC<ISourceListProps> = ({ data }) => {
  return (
    <AwaitWithFallback resolve={data.sources}>
      {(sources: Source[]) =>
        sources.length === 0 ? (
          <div>No sources created</div>
        ) : (
          <ul>
            {sources.map((source) => (
              <li key={`source-${source.id}`} className="my-half">
                <Link to={`sources/${source.id}`} className="flex space-x-basic">
                  <DocumentTextIcon className="h-6 w-6" />
                  <div>{source.name}</div>
                </Link>
              </li>
            ))}
          </ul>
        )
      }
    </AwaitWithFallback>
  )
}

export default SourceList
