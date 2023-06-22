import SourceEdit, { ISourceEditData } from "./components/source/SourceEdit.tsx"
import { useLoaderData } from "react-router-dom"

const SourceShowPage: React.FC = () => {
  const data = useLoaderData() as ISourceEditData
  return <SourceEdit data={data} />
}

export default SourceShowPage
