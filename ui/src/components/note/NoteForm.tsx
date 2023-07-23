import React from "react"

const NoteForm: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  return (
    <form className="m-std space-y-std">
      <div>
        <label>Name</label>
        <input type="text" />
      </div>

      <div className="my-std justify-end flex-std">
        <button type="button" onClick={onClose}>
          Cancel
        </button>
        <button type="submit">Create</button>
      </div>
    </form>
  )
}

export default NoteForm
