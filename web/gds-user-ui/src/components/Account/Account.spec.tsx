import { act, render, screen } from 'utils/test-utils';
import { Account } from '.';
import { faker } from '@faker-js/faker';
import { dynamicActivate } from 'utils/i18nLoaderHelper';

describe('<Account />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should display props correctly', () => {
    const { username, vaspName } = {
      vaspName: faker.internet.domainName(),
      username: faker.internet.userName()
    };
    render(<Account username={username} vaspName={vaspName} />);

    expect(screen.getByTestId('vaspName')).toBeInTheDocument();
    expect(screen.getByTestId('username')).toBeInTheDocument();
  });
});
