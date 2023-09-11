import { Source } from "../../services/SourcesService.ts"
import AwaitWithFallback from "../AwaitWithFallback.tsx"
import React from "react"
import { Form, Link } from "react-router-dom"

export interface ISourceEditData {
  source: Promise<Source>
}
interface ISourceEditProps {
  data: ISourceEditData
}

const SourceEdit: React.FC<ISourceEditProps> = ({ data }) => {
  return (
    <AwaitWithFallback resolve={data.source}>
      {(source: Source) => <SourceComponent source={source} />}
    </AwaitWithFallback>
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

        <div className="flex justify-end space-x-basic">
          <Link to={`/sources/${source.id}`} className="btn">
            Cancel
          </Link>
          <button type="submit" className="btn-primary">
            Submit
          </button>
        </div>
      </div>
    </Form>
  )
}

export default SourceEdit
