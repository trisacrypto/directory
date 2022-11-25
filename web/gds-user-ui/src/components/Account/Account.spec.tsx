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
    const { name, domain, onClose } = {
      name: faker.internet.domainName(),
      domain: faker.internet.userName(),
      onClose: jest.fn()
    };
    render(<Account name={name} domain={domain} onClose={onClose} />);

    expect(screen.getByTestId('vaspName')).toBeInTheDocument();
    expect(screen.getByTestId('vaspDomain')).toBeInTheDocument();
  });
});
