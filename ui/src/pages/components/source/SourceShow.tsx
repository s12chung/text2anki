import { Source } from "../../../services/SourceService.ts"
import AwaitError from "../AwaitError.tsx"
import React from "react"
import { Await } from "react-router-dom"

export interface ISourceShowData {
  source: Promise<Source[]>
}
interface ISourceShowProps {
  data: ISourceShowData
}

const SourceShow: React.FC<ISourceShowProps> = ({ data }) => {
  return (
    <React.Suspense fallback={<div>Loading....</div>}>
      <Await resolve={data.source} errorElement={<AwaitError />}>
        {(source: Source) => (
          <div>
            <div>{source.name}</div>
            {source.tokenizedTexts.map((tokenizedText) => (
              <div key={`text-${tokenizedText.text}`}>
                <div>{tokenizedText.text}</div>
                <div>{tokenizedText.translation}</div>
              </div>
            ))}
          </div>
        )}
      </Await>
    </React.Suspense>
  )
}

export default SourceShow
