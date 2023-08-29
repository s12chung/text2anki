import { Source } from "../../services/SourcesService.ts"
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
      <div className="space-y-std">
        <label>
          Name:
          <input name="name" type="text" defaultValue={source.name} />
        </label>
        <label>
          Reference:
          <input name="reference" type="text" defaultValue={source.reference} />
        </label>

        <Link to={`/sources/${source.id}`} className="btn">
          Cancel
        </Link>
        <button type="submit">Submit</button>
      </div>
    </Form>
  )
}

export default SourceEdit
