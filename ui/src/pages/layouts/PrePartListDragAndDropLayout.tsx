import PrePartListDragAndDrop from "../../components/pre_part_lists/PrePartListDragAndDrop.tsx"
import ApplicationLayout from "./ApplicationLayout.tsx"
import React from "react"

const SourceDragAndDropLayout: React.FC = () => {
  return (
    <PrePartListDragAndDrop minHeight="min-h-screen">
      <ApplicationLayout />
    </PrePartListDragAndDrop>
  )
}

export default SourceDragAndDropLayout
