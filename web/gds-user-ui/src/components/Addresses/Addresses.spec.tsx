import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import Addresses from '.';

describe('<Addresses />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should add address row', () => {
    render(<Addresses />);

    const addAddress = screen.getByRole('button', { name: /add address/i });

    userEvent.click(addAddress);
    userEvent.click(addAddress);

    expect(screen.getAllByTestId('address-row')).toHaveLength(2);
  });
});
