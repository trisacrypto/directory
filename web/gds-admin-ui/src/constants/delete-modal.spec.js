import DeleteModal from '@/components/DeleteModal';
import { Modal, ModalContent } from '@/components/Modal';
import { fireEvent, render, screen } from '@/utils/test-utils';

const ModalWrapper = ({ children }) => (
  <Modal value={[true, () => {}]}>
    <ModalContent>{children}</ModalContent>
  </Modal>
);

describe('<DeleteModal />', () => {
  it('should delete when clicking on delete button', async () => {
    const isLoading = false;
    const handleDeleteClick = jest.fn();

    render(
      <ModalWrapper>
        <DeleteModal onDelete={handleDeleteClick} isLoading={isLoading} />
      </ModalWrapper>
    );

    const deleteBtn = screen.getByTestId('deleteBtn');

    fireEvent.click(deleteBtn);

    expect(handleDeleteClick).toHaveBeenCalledTimes(1);
  });
});
