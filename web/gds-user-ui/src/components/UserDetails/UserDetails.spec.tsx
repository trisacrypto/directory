import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import UserDetails from '.';

describe('<UserDetails />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render props', () => {
    const userDetails = {
      userId: 'C0000213',
      createdDate: '01/01/2020',
      status: 'Active',
      permissions: 'Admin',
      lastLogin: '01/01/2020'
    };
    render(
      <UserDetails
        userId={userDetails.userId}
        createdDate={userDetails.createdDate}
        status={userDetails.status}
        permissions={userDetails.permissions}
        lastLogin={userDetails.lastLogin}
      />
    );

    expect(screen.getByTestId('user_id').textContent).toBe('User ID: C0000213');
    expect(screen.getByTestId('profile_created').textContent).toBe('Profile Created: 01/01/2020');
    expect(screen.getByTestId('status').textContent).toBe('Status: Active');
    expect(screen.getByTestId('last_login').textContent).toBe('Last Login: 01/01/2020');
    expect(screen.getByTestId('permissions').textContent).toBe('Permission: Admin');
  });

  it('should match inline snapshot', () => {
    const userDetails = {
      userId: 'C0000213',
      createdDate: '01/01/2020',
      status: 'Active',
      permissions: 'Admin',
      lastLogin: '01/01/2020'
    };
    const { container } = render(
      <UserDetails
        userId={userDetails.userId}
        createdDate={userDetails.createdDate}
        status={userDetails.status}
        permissions={userDetails.permissions}
        lastLogin={userDetails.lastLogin}
      />
    );

    expect(container).toMatchInlineSnapshot(`
      .emotion-0 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        margin-top: var(--chakra-space-10);
      }

      .emotion-1 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-align-items: center;
        -webkit-box-align: center;
        -ms-flex-align: center;
        align-items: center;
        -webkit-flex-direction: column;
        -ms-flex-direction: column;
        flex-direction: column;
      }

      .emotion-1>*:not(style)~*:not(style) {
        margin-top: var(--chakra-space-4);
        -webkit-margin-end: 0px;
        margin-inline-end: 0px;
        margin-bottom: 0px;
        -webkit-margin-start: 0px;
        margin-inline-start: 0px;
      }

      .emotion-2 {
        margin-top: var(--chakra-space-2);
      }

      .emotion-3 {
        font-family: var(--chakra-fonts-heading);
        font-weight: var(--chakra-fontWeights-bold);
        font-size: var(--chakra-fontSizes-xl);
        line-height: 1.2;
        padding-bottom: var(--chakra-space-3);
      }

      <div>
        <div
          class="emotion-0"
        >
          <div
            class="chakra-stack emotion-1"
          >
            <div
              class="emotion-2"
            >
              <h2
                class="chakra-heading emotion-3"
              >
                User Details
              </h2>
              <p
                class="chakra-text emotion-4"
                data-testid="user_id"
              >
                User ID:
                 
                C0000213
              </p>
              <p
                class="chakra-text emotion-4"
                data-testid="profile_created"
              >
                Profile Created:
                 
                01/01/2020
              </p>
              <p
                class="chakra-text emotion-4"
                data-testid="status"
              >
                Status:
                 
                Active
              </p>
              <p
                class="chakra-text emotion-4"
                data-testid="permissions"
              >
                Permission:
                 
                Admin
              </p>
              <p
                class="chakra-text emotion-4"
                data-testid="last_login"
              >
                Last Login:
                 
                01/01/2020
              </p>
            </div>
          </div>
        </div>
      </div>
    `);
  });
});
