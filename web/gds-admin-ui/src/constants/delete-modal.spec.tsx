import DeleteModal from '@/components/DeleteModal';
import { Modal, ModalContent } from '@/components/Modal';
import { fireEvent, render, screen } from '@/utils/test-utils';
import { ReactNode } from 'react';

const ModalWrapper = ({ children }: { children: ReactNode }) => (
    <Modal
        value={{
            isOpen: true,
        }}>
        <ModalContent>{children}</ModalContent>
    </Modal>
);

describe('<DeleteModal />', () => {
    it('should delete when clicking on delete button', async () => {
        const isLoading = false;
        const handleDeleteClick = jest.fn();

        render(
            <ModalWrapper>
                <DeleteModal onDelete={handleDeleteClick} isLoading={isLoading} vaspId="" vasp={{}} />
            </ModalWrapper>
        );

        const deleteBtn = screen.getByTestId('deleteBtn');

        fireEvent.click(deleteBtn);

        expect(handleDeleteClick).toHaveBeenCalledTimes(1);
    });
});
