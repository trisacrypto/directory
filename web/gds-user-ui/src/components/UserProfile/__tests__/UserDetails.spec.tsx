import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import UserDetails from '../UserDetails';
const userDetailsMock = {
  createAt: '2021-05-05T12:00:00.000Z',
  role: 'admin',
  permissions: ['admin', 'user'],
  lastLogin: '2021-05-05T12:00:00.000Z'
};
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useSelector: jest.fn().mockReturnValueOnce({
    user: userDetailsMock
  })
}));

describe('<UserDetails />  ', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render user details', () => {
    render(<UserDetails />);

    expect(screen.getByTestId('user_created_At').textContent).toBe(userDetailsMock.createAt);
    expect(screen.getByTestId('user_role').textContent).toBe(userDetailsMock.role);
    expect(screen.getByTestId('user_last_login').textContent).toBe(userDetailsMock.lastLogin);
    // get all permissions by index and check if they are in the array
    expect(screen.getAllByTestId('user_permissions')).toHaveLength(
      userDetailsMock.permissions.length
    );
    // afterEach(() => {
    //   useSelector.mockClear();
    // });
  });
});
