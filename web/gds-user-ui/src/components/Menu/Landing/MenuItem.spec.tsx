import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen, waitFor } from 'utils/test-utils';
import MenuItem from './MenuItem';

describe('<MenuItem />', () => {
  beforeAll(async () => {
    await waitFor(() => {
      dynamicActivate('en');
    });
  });
  it('should', () => {
    const { debug } = render(<MenuItem to={'/about'}>Test</MenuItem>);

    const linkEl = screen.getByRole('link', { name: /test/i });
    expect(linkEl).toHaveTextContent('Test');
    expect(linkEl).toHaveAttribute('href', '/about');
  });
});
