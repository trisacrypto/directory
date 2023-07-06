import { act, render, screen, fireEvent, waitFor } from 'utils/test-utils';
import ShowMemberModal from '../';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { memberDetailMock } from '../../../__mocks__';

// mock chakra modal
jest.mock('@chakra-ui/react', () => {
  const originalModule = jest.requireActual('@chakra-ui/react');

  return {
    __esModule: true,
    ...originalModule,
    useDisclosure: jest.fn(() => ({
      isOpen: true,
      onOpen: jest.fn(),
      onClose: jest.fn()
    }))
  };
});

function renderComponent(memberId?: string) {
  const props = memberId || memberDetailMock.summary.id;

  return render(<ShowMemberModal memberId={props} />);
}

describe('Member Modal', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the modal', async () => {
    const { getByTestId } = renderComponent();
    const modalBtn = getByTestId('member-modal-button');

    fireEvent.click(modalBtn);

    waitFor(() => {
      expect(screen.getByTestId('member-modal')).toBeInTheDocument();
    });
  });

  it('should render the member details in the modal', () => {
    const { getByTestId } = renderComponent();
    const modalBtn = getByTestId('member-modal-button');

    fireEvent.click(modalBtn);

    waitFor(() => {
      expect(screen.getByText(`${memberDetailMock.summary.name}`)).toBeInTheDocument();
      expect(screen.getByText(`${memberDetailMock.summary.website}`)).toBeInTheDocument();
      expect(screen.getByText(`${memberDetailMock.summary.country}`)).toBeInTheDocument();
      expect(screen.getByText(`${memberDetailMock.summary.endpoint}`)).toBeInTheDocument();
      expect(screen.getByText(`${memberDetailMock.summary.common_name}`)).toBeInTheDocument();
    });
  });
});
