/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { fireEvent, screen } from '@testing-library/react';
import { useMutation } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import nock from 'nock';
import DeleteCollaboratorModal from '../DeleteCollaboratorModal';
import { act, render } from 'utils/test-utils';
import { collaboratorMockValue } from 'components/AddCollaboratorModal/__mocks__';
import * as useCollaborators from 'components/Collaborators/useFetchCollaborator';
import * as useDeleteCollaborator from '../useDeleteCollaborator';
const mockUseMutation = useMutation as jest.Mock;
// mock chakra ui modal component to be able to test it

const divWithChildrenMock = (children: any, identifier: any) => (
  <div data-testId={identifier}>{children}</div>
);
const divWithoutChildrenMock = (identifier: any) => <div data-testId={identifier} />;

jest.mock('@chakra-ui/react', () => ({
  ...jest.requireActual('@chakra-ui/react'),
  Modal: jest.fn(({ children }) => divWithChildrenMock(children, 'modal')),
  ModalOverlay: jest.fn(({ children }) => divWithChildrenMock(children, 'overlay')),
  ModalContent: jest.fn(({ children }) => divWithChildrenMock(children, 'content')),
  ModalHeader: jest.fn(({ children }) =>
    divWithChildrenMock(children, 'delete-collaborator-modal')
  ),
  ModalFooter: jest.fn(({ children }) => divWithChildrenMock(children, 'footer')),
  ModalBody: jest.fn(({ children }) => divWithChildrenMock(children, 'body')),
  ModalCloseButton: jest.fn(() => divWithoutChildrenMock('close'))
}));

// mock deletecollaborator component render function
function renderComponent() {
  const Props = {
    collaboratorId: '1'
  };
  return render(<DeleteCollaboratorModal {...Props} />);
}

const mockCollaborators = jest.fn();
const mockGetAllCollaborators = jest.fn();
const mockDeleteCollaborator = jest.fn();

const useFetchCollaboratorsMock = jest.spyOn(useCollaborators, 'useFetchCollaborators');
const useDeleteCollaboratorMock = jest.spyOn(useDeleteCollaborator, 'useDeleteCollaborator');
describe('DeleteCollaboratorModal', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
    useFetchCollaboratorsMock.mockReturnValue({
      collaborators: {
        data: {
          collaborators: collaboratorMockValue.data
        },
        getAllCollaborators: mockGetAllCollaborators
      },
      getAllCollaborators: jest.fn(),
      hasCollaboratorsFailed: false,
      wasCollaboratorsFetched: false,
      isFetchingCollaborators: false
    });
    useDeleteCollaboratorMock.mockReturnValue({
      isDeleting: false,
      wasCollaboratorDeleted: false,
      deleteCollaborator: mockDeleteCollaborator,
      hasCollaboratorFailed: false,
      errorMessage: '',
      reset(): void {
        throw new Error('Function not implemented.');
      }
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the modal', () => {
    renderComponent();
    expect(screen.getByTestId('delete-collaborator-modal')).toBeInTheDocument();
  });

  it('useCollaborators should be called', () => {
    renderComponent();
    expect(useFetchCollaboratorsMock).toHaveBeenCalled();
  });

  it('useDeleteCollaborator should be called', () => {
    renderComponent();
    expect(useDeleteCollaboratorMock).toHaveBeenCalled();
  });

  it('should show collaborator email in the modal', () => {
    renderComponent();
    expect(screen.getByTestId('collaborator-email').textContent).toBe(
      collaboratorMockValue.data[0].email
    );
  });

  it('should show collaborator name in the modal', () => {
    renderComponent();
    expect(screen.getByTestId('collaborator-name').textContent).toBe(
      collaboratorMockValue.data[0].name
    );
  });

  // delete collaborator click event test
  // it('should call deleteHandler function when delete button is clicked', () => {
  //   renderComponent();
  //   fireEvent.click(screen.getByTestId('delete-collaborator-button'));
  //
  //   expect(mockDeleteCollaborator).toHaveBeenCalled();
  // });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
