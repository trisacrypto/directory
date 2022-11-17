/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { fireEvent, screen } from '@testing-library/react';
import { useMutation } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import nock from 'nock';
import EditCollaboratorModal from '../EditCollaboratorModal';
import { act, render } from 'utils/test-utils';
import { collaboratorMockValue } from 'components/Collaborators/AddCollaborator/__mocks__';
import * as useCollaborators from 'components/Collaborators/useFetchCollaborator';
import * as useUpdateCollaborator from '../useUpdateCollaborator';
const mockUseMutation = useMutation as jest.Mock;

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
    divWithChildrenMock(children, 'update-collaborator-modal')
  ),
  ModalFooter: jest.fn(({ children }) => divWithChildrenMock(children, 'footer')),
  ModalBody: jest.fn(({ children }) => divWithChildrenMock(children, 'body')),
  ModalCloseButton: jest.fn(() => divWithoutChildrenMock('close'))
}));

// render delete collaborator component
function renderComponent() {
  const Props = {
    collaboratorId: '1',
    roles: ['ADMIN']
  };
  return render(<EditCollaboratorModal {...Props} />);
}

const mockCollaborators = jest.fn();
const mockGetAllCollaborators = jest.fn();
const mockUpdateCollaborator = jest.fn();
const mockCollaboratorsData = collaboratorMockValue.data;

const useFetchCollaboratorsMock = jest.spyOn(useCollaborators, 'useFetchCollaborators');
const useUpdateCollaboratorMock = jest.spyOn(useUpdateCollaborator, 'useUpdateCollaborator');
describe('UpdateCollaboratorModal', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
    useFetchCollaboratorsMock.mockReturnValue({
      collaborators: mockCollaboratorsData,
      getAllCollaborators: jest.fn(),
      hasCollaboratorsFailed: false,
      wasCollaboratorsFetched: false,
      isFetchingCollaborators: false
    });
    useUpdateCollaboratorMock.mockReturnValue({
      isUpdating: false,
      wasCollaboratorUpdated: false,
      updateCollaborator: mockUpdateCollaborator,
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
    expect(screen.getByTestId('update-collaborator-modal')).toBeInTheDocument();
  });

  it('useCollaborators should be called', () => {
    renderComponent();
    expect(useFetchCollaboratorsMock).toHaveBeenCalled();
  });

  it('useUpdateCollaborator should be called', () => {
    renderComponent();
    expect(useUpdateCollaboratorMock).toHaveBeenCalled();
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

  // collaborator without update:collaborator permission should not be able to update collaborator
  it('should disable edit button if user does not have update:collaborator permission', () => {
    // mock shouldDisableDeleteButton to return true
    const mockShouldDisableDeleteButton = jest.fn().mockReturnValue(true);
    jest.mock('components/Collaborators/useSafeDisableButton', () => ({
      ...jest.requireActual('hooks/useSafeDisableButton'),
      useSafeDisableButton: () => ({
        isDisabled: mockShouldDisableDeleteButton()
      })
    }));
    renderComponent();
    const editButton = screen.getByTestId('collaborator-button');
    fireEvent.click(editButton);
    expect(editButton).toBeDisabled();
  });

  // it('should call deleteHandler function when delete button is clicked', () => {
  //   renderComponent();
  //   userEvent.click(screen.getByTestId('delete-collaborator-button'));

  //   expect(mockDeleteCollaborator).toHaveBeenCalled();
  // });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
