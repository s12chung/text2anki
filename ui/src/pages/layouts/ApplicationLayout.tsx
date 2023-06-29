import { Link, Outlet } from "react-router-dom"

const ApplicationLayout: React.FC = () => {
  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 md:px-8">
      <div className="my-std">
        <Link to="/">text2anki</Link>
      </div>
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default ApplicationLayout
