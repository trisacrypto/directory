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

      render(<LandingHeader />, { locale: 'en' });

      expect(screen.getByText(/documentation/i)).toHaveAttribute('href', `${TRISA_BASE_URL}/en`);
      expect(screen.getByText(/about trisa/i)).toHaveAttribute('href', 'https://trisa.io');
    });

    it('should target the german website', async () => {
      await waitFor(() => {
        dynamicActivate('de');
      });
      render(<LandingHeader />, { locale: 'de' });

      expect(screen.getByText(/Dokumentation/i)).toHaveAttribute('href', `${TRISA_BASE_URL}/de`);
    });

    it('should target the french website', async () => {
      await waitFor(() => {
        dynamicActivate('fr');
      });
      render(<LandingHeader />, { locale: 'fr' });

      expect(screen.getByText(/documentation/i)).toHaveAttribute('href', `${TRISA_BASE_URL}/fr`);
    });
  });
});
