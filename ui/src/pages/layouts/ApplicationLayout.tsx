import Header from "../../components/Header.tsx"
import { Outlet } from "react-router-dom"

const ApplicationLayout: React.FC = () => {
  return (
    <div className="grid-std">
      <Header />
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default ApplicationLayout
