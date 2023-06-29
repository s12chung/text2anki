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
    <>
      <div className="flex-std my-std">
        <div className="flex-grow">
          <h2>{source.name}</h2>
        </div>
        <div className="flex">
          <Link to={`/sources/${source.id}/edit`} className="btn">
            Edit
          </Link>
        </div>
      </div>

      {source.tokenizedTexts.map((tokenizedText) => (
        <div key={`text-${tokenizedText.text}`}>
          <div>{tokenizedText.text}</div>
          <div>{tokenizedText.translation}</div>
        </div>
      ))}
    </>
  )
}

export default SourceShow
