import { Modal } from '@/components/Modal';
import BusinessInfosForm from '@/pages/app/details/BasicDetails/components/BusinessInfosForm';
import { act, fireEvent, render, screen } from '@/utils/test-utils';

describe('BusinessInfosForm', () => {
  const data = {
    name: 'Opalcliff, Inc.',
    vasp: {
      business_category: 'BUSINESS_ENTITY',
      established_on: '275760-08-08',
      vasp_categories: ['P2P'],
      website: 'http://opalcliff.com',
    },
  };

  describe('website', () => {
    it('should be a valid website', async () => {
      render(
        <Modal>
          <BusinessInfosForm data={data} />
        </Modal>
      );

      const websiteEl = screen.getByRole('textbox', { name: /website/i });
      await act(() => {
        fireEvent.change(websiteEl, { target: { value: 'http://opalcliff.co' } });
      });

      expect(websiteEl).toHaveClass('is-valid');
    });

    it('should not be a valid website', async () => {
      render(
        <Modal>
          <BusinessInfosForm data={data} />
        </Modal>
      );

      const websiteEl = screen.getByRole('textbox', { name: /website/i });
      await act(() => {
        fireEvent.change(websiteEl, { target: { value: 'http:/opalcliff' } });
      });
      const websiteErrorMessageEl = screen.getByText(/website should be a valid url/i);

      expect(websiteErrorMessageEl).toBeInTheDocument();
      expect(websiteEl).toHaveClass('is-invalid');
    });
  });

  describe('Date of Incorporation/Establishment', () => {
    it('should be a valid date', async () => {
      render(
        <Modal>
          <BusinessInfosForm data={data} />
        </Modal>
      );

      const establishedOnEl = screen.getByLabelText(/date of incorporation\/establishment/i);
      await act(() => {
        fireEvent.change(establishedOnEl, { target: { value: '2020-10-12' } });
      });

      expect(establishedOnEl).toHaveClass('is-valid');
    });

    it('should not be a valid date', async () => {
      render(
        <Modal>
          <BusinessInfosForm data={data} />
        </Modal>
      );

      const establishedOnEl = screen.getByLabelText(/date of incorporation\/establishment/i);
      await act(async () => {
        fireEvent.change(establishedOnEl, { target: { value: '10-12-2020' } });
      });

      const establishedOnErrorMessageEl = screen.getByText(
        /date of incorporation\/establishment should be a valid date/i
      );

      expect(establishedOnErrorMessageEl).toBeInTheDocument();
      expect(establishedOnEl).toHaveClass('is-invalid');
    });
  });
});
