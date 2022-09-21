import { TRISA_BASE_URL } from 'constants/trisa-base-url';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { waitFor, render, screen } from 'utils/test-utils';
import LandingHeader from './LandingHeader';

describe('<LandingHeader />', () => {
  describe('Documentation menu', () => {
    it('should target the english website', async () => {
      await waitFor(() => {
        dynamicActivate('en');
      });
      const { debug } = render(<LandingHeader />, { locale: 'en' });
      const documentation = screen.getByText(/documentation/i);
      expect(documentation).toHaveAttribute('href', `${TRISA_BASE_URL}/en`);
    });

    it('should target the german website', async () => {
      await waitFor(() => {
        dynamicActivate('de');
      });
      render(<LandingHeader />, { locale: 'de' });

      const documentation = screen.getByText(/Dokumentation/i);
      expect(documentation).toHaveAttribute('href', `${TRISA_BASE_URL}/de`);
    });

    it('should target the french website', async () => {
      await waitFor(() => {
        dynamicActivate('fr');
      });
      render(<LandingHeader />, { locale: 'fr' });

      const documentation = screen.getByText(/documentation/i);
      expect(documentation).toHaveAttribute('href', `${TRISA_BASE_URL}/fr`);
    });
  });
});
