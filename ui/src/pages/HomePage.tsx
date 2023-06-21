import SourceForm from "./components/source/SourceForm.tsx"
import SourceList, { ISourceListData } from "./components/source/SourceList"
import { useLoaderData } from "react-router-dom"

const HomePage: React.FC = () => {
  const data = useLoaderData() as ISourceListData
  return (
    <div>
      <SourceForm />
      <h1>Sources</h1>
      <SourceList data={data} />
    </div>
  )
}

export default HomePage
