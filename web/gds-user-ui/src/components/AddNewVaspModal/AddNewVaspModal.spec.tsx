import { dynamicActivate } from 'utils/i18nLoaderHelper';
import * as permission from 'utils/permission';
import { render, screen } from 'utils/test-utils';
import AddNewVaspModal from './AddNewVaspModal';

// jest.mock('utils/permission');

describe('<AddNewVaspModal />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should be disable when user has not the right permission', () => {
    render(<AddNewVaspModal />);

    jest.mock('utils/permission', () => ({
      canCreateOrganization: () => false
    }));

    const addNewVaspButton = screen.getByTestId('add-new-vasp');
    expect(addNewVaspButton).toBeDisabled();
  });

  it('should be enable when user has the right permission', () => {
    render(<AddNewVaspModal />);

    jest.mock('utils/permission', () => ({
      canCreateOrganization: () => true
    }));

    const addNewVaspButton = screen.getByTestId('add-new-vasp');
    expect(addNewVaspButton).toBeDisabled();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });
});
