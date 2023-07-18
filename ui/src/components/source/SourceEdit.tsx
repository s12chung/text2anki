import { Source } from "../../services/SourceService.ts"
import AwaitError from "../AwaitError.tsx"
import React from "react"
import { Await, Form, Link } from "react-router-dom"

export interface ISourceEditData {
  source: Promise<Source>
}
interface ISourceEditProps {
  data: ISourceEditData
}

const SourceEdit: React.FC<ISourceEditProps> = ({ data }) => {
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
    <Form action={`/sources/${source.id}`} method="patch">
      <div className="flex-std">
        <input name="name" type="text" defaultValue={source.name} className="flex-grow" />

        <Link to={`/sources/${source.id}`} className="btn">
          Cancel
        </Link>
        <button type="submit">Submit</button>
      </div>
    </Form>
  )
}

export default SourceEdit
