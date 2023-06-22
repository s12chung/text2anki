import SourceCreate from "./components/source/SourceCreate.tsx"
import SourceList, { ISourceListData } from "./components/source/SourceList"
import { useLoaderData } from "react-router-dom"

const HomePage: React.FC = () => {
  const data = useLoaderData() as ISourceListData
  return (
    <div>
      <SourceCreate />
      <h1>Sources</h1>
      <SourceList data={data} />
    </div>
  )
}

export default HomePage
