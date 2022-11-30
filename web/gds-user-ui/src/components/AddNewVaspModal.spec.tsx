import { Modal, ModalContent } from '@chakra-ui/react';
import userEvent from '@testing-library/user-event';
import { ReactNode } from 'react';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import AddNewVaspForm from './AddNewVaspForm';

const ModalWrapper = ({ children }: { children: ReactNode }) => (
  <Modal isOpen onClose={jest.fn()}>
    <ModalContent>{children}</ModalContent>
  </Modal>
);

describe('<AddNewVaspModal />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  // it('should be disabled when checkbox is unchecked', () => {
  //   render(
  //     <ModalWrapper>
  //       <AddNewVaspForm isCreatingVasp={false} onSubmit={jest.fn()} closeModal={jest.fn()} />
  //     </ModalWrapper>
  //   );

  //   expect(screen.getByTestId('accept').querySelector('input[type="checkbox"]')).not.toBeChecked();

  //   expect(screen.getByTestId('name')).toBeDisabled();
  //   expect(screen.getByTestId('domain')).toBeDisabled();
  //   expect(screen.getByRole('button', { name: /next/i })).toBeDisabled();
  // });

  it('should be enabled when checkbox is checked', () => {
    render(
      <ModalWrapper>
        <AddNewVaspForm isCreatingVasp={false} onSubmit={jest.fn()} closeModal={jest.fn()} />
      </ModalWrapper>
    );

    const acceptElement = screen.getByTestId('accept').querySelector('input[type="checkbox"]');
    userEvent.click(acceptElement!);

    expect(screen.getByTestId('accept').querySelector('input[type="checkbox"]')).toBeChecked();

    expect(screen.getByTestId('name')).toBeEnabled();
    expect(screen.getByTestId('domain')).toBeEnabled();
    expect(screen.getByRole('button', { name: /next/i })).toBeEnabled();
  });
});
