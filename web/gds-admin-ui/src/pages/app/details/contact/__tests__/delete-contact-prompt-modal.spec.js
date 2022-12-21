import { render, screen, fireEvent } from "utils/test-utils"
import DeleteContactPromptModal from "pages/app/details/contact/DeleteContactPromptModal"
import { Modal } from "components/Modal"

describe('<DeleteContactPromptModal />', () => {

    it("should delete a contact", async () => {
        const handleDelete = jest.fn()

        render(
            <Modal>
                <DeleteContactPromptModal onDelete={handleDelete} type="legal" />
            </Modal>
        )

        const deleteEl = screen.getByRole('button', { name: /delete/i });
        fireEvent.click(deleteEl)
        expect(handleDelete).toHaveBeenCalledTimes(1)
    })
})