import userEvent from "@testing-library/user-event"
import DeleteModal from "components/DeleteModal"
import { Modal, ModalContent } from "components/Modal"
import { render, screen, waitFor } from "utils/test-utils"

const ModalWrapper = ({ children }) => {
    return (
        // this mount the modal component
        <Modal
            value={[true, () => { }]}
        >
            <ModalContent>
                {children}
            </ModalContent>
        </Modal>
    )
}

describe('<DeleteModal />', () => {

    it("should delete when clicking on delete button", async () => {
        const isLoading = false
        const handleDeleteClick = jest.fn()

        render(
            <ModalWrapper>
                <DeleteModal onDelete={handleDeleteClick} isLoading={isLoading} />
            </ModalWrapper>
        )

        const deleteBtn = screen.getByTestId('deleteBtn')
        await waitFor(() => {
            userEvent.click(deleteBtn)
        })

        expect(handleDeleteClick).toHaveBeenCalledTimes(1)
    })
})