import { Source } from "../../../services/SourceService.ts"
import AwaitError from "../AwaitError.tsx"
import React from "react"
import { Await, Link } from "react-router-dom"

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
      <div className="flex">
        <div className="flex-grow">{source.name}</div>
        <Link to={`/sources/${source.id}/edit`}>Edit</Link>
      </div>

      {source.tokenizedTexts.map((tokenizedText) => (
        <div key={`text-${tokenizedText.text}`}>
          <div>{tokenizedText.text}</div>
          <div>{tokenizedText.translation}</div>
        </div>
      ))}
    </div>
  )
}

export default SourceShow