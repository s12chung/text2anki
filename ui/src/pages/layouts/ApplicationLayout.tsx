import { Outlet } from "react-router-dom"

const ApplicationLayout: React.FC = () => {
  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 md:px-8">
      <div className="my-4">text2anki</div>
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default ApplicationLayout
