import Header from "../../components/Header.tsx"
import { Outlet } from "react-router-dom"

const FullLayout: React.FC = () => {
  return (
    <div className="m-std">
      <Header />
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default FullLayout
