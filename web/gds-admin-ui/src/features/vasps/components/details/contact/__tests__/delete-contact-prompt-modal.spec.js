import { Modal } from '@/components/Modal';
import { fireEvent, render, screen } from '@/utils/test-utils';
import DeleteContactPromptModal from '../DeleteContactPromptModal';

describe('<DeleteContactPromptModal />', () => {
  it('should delete a contact', async () => {
    const handleDelete = jest.fn();

    render(
      <Modal>
        <DeleteContactPromptModal onDelete={handleDelete} type="legal" />
      </Modal>
    );

    const deleteEl = screen.getByRole('button', { name: /delete/i });
    fireEvent.click(deleteEl);
    expect(handleDelete).toHaveBeenCalledTimes(1);
  });
});
