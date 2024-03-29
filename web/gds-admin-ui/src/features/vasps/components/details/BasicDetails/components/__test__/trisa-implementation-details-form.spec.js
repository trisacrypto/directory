/* eslint-disable testing-library/no-unnecessary-act */
import { Modal } from '@/components/Modal';
import TrisaImplementationDetailsForm from '../TrisaImplementationDetailsForm';
import { act, fireEvent, render, screen } from '@/utils/test-utils';

const TrisaForm = ({ data }) => (
  <Modal>
    <TrisaImplementationDetailsForm data={data} />
  </Modal>
);

describe('TrisaImplementationDetailsForm', () => {
  const data = {
    vasp: {
      common_name: 'traveler.ciphertrace.com',
      trisa_endpoint: 'traveler.ciphertrace.com:443',
    },
  };

  it('submit button should be disable when the form is not dirty', () => {
    render(<TrisaForm data={data} />);

    const submitEl = screen.getByText(/save/i);

    expect(submitEl).toBeDisabled();
  });

  it('trisa endpoint should not start by zero', async () => {
    render(<TrisaForm data={data} />);

    const trisaEndpointEl = screen.getByRole('textbox', { name: /trisa endpoint/i });

    await act(async () => {
      fireEvent.change(trisaEndpointEl, { target: { value: 'traveler.ciphertrace.com:043' } });
    });

    const errorMessageEl = screen.getByText(/trisa endpoint is not valid/i);

    expect(trisaEndpointEl).toHaveClass('is-invalid');
    expect(errorMessageEl).toBeInTheDocument();
    expect(errorMessageEl).toHaveClass('invalid-feedback');
  });

  it('trisa endpoint should not start by a http', async () => {
    render(<TrisaForm data={data} />);

    const trisaEndpointEl = screen.getByRole('textbox', { name: /trisa endpoint/i });

    await act(async () => {
      fireEvent.change(trisaEndpointEl, {
        target: { value: 'https://traveler.ciphertrace.com:443' },
      });
    });
    const errorMessageEl = screen.getByText(/trisa endpoint is not valid/i);

    expect(trisaEndpointEl).toHaveClass('is-invalid');
    expect(errorMessageEl).toBeInTheDocument();
    expect(errorMessageEl).toHaveClass('invalid-feedback');
  });

  it("should show a warning when common name doesn't match trisa endpoint without the port", async () => {
    render(<TrisaForm data={data} />);

    const commonNameEl = screen.getByRole('textbox', { name: /certificate common name/i });

    await act(async () => {
      fireEvent.change(commonNameEl, {
        target: {
          value: 'traveler.ciphertrace.co',
        },
      });
    });

    const warningMessageEl = screen.getByText(
      /common name should match the trisa endpoint without the port/i
    );
    expect(warningMessageEl).toBeInTheDocument();
    expect(warningMessageEl).toHaveClass('text-warning');
  });
});
